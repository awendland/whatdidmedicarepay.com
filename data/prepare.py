import urllib.request
import zipfile
import shutil
import tempfile
from pathlib import Path
import hashlib
import base64
import subprocess
import os

# CMS has produced a PDF describing the methodology and schema of the dataset
# https://www.cms.gov/Research-Statistics-Data-and-Systems/Statistics-Trends-and-Reports/Medicare-Provider-Charge-Data/Downloads/Medicare-Physician-and-Other-Supplier-PUF-Methodology.pdf

# TODO it appears that a 2018 dataset is available
DATA_URL = os.environ.get(
    "DATA_URL",
    "http://download.cms.gov/Research-Statistics-Data-and-Systems/Statistics-Trends-and-Reports/Medicare-Provider-Charge-Data/Downloads/Medicare_Provider_Util_Payment_PUF_CY2017.zip",
)
DATA_FILE_ID = os.environ.get(
    "DATA_FILE_ID", "Medicare_Provider_Util_Payment_PUF_CY2017.txt"
)

TMP_DIR = Path(os.environ.get("TMP_DIR", tempfile.gettempdir()))
OPTIMIZE = os.environ.get("OPTIMIZE", "False") == "True"


def stable_filename(value: str):
    """Generate a stable filename by Base32 encoding the sha256 of a string"""
    return str(
        base64.b32encode(hashlib.sha256(value.encode("utf-8")).digest()), "utf-8"
    ).replace("====", "")


def file_cached(file_name: str, data_readable=None, data_gen=None, step_name=None):
    """
    Wrap a time-intensive operation with a file-system cache. Values from the cache will always be
    returned, if available.

    * Provide a paramaterless function to data_readable if the value is accessible as a _Reader[AnyStr].
    * Provide a 1 parameter function to data_gen if the value will be generated by a function that writes
      directly to disk.

    If a value is written to the file system while it is being generated, the filename will be suffixed
    with '.part' to indicate that the generation has not completed. Once the value has been fully
    generated then the file will be renamed to remove the '.part' suffix. This is similar to how Safari
    downloads files.
    """
    if Path(file_name).exists():
        if step_name:
            print("Using cached {} at {}".format(step_name, file_name))
    else:
        Path(file_name).parent.mkdir(parents=True, exist_ok=True)
        if step_name:
            print("Caching {} at {}".format(step_name, file_name))
        file_name_part = Path(str(file_name) + ".part")
        if file_name_part.exists():
            file_name_part.unlink()
        if data_readable is not None:
            with open(file_name_part, mode="wb") as part_file:
                shutil.copyfileobj(data_readable(), part_file)
        elif data_gen is not None:
            data_gen(file_name_part)
        file_name_part.rename(file_name)
    return Path(file_name)


# Download the Internet-accessible zipped data file
tmp_zip_path = TMP_DIR / (stable_filename(DATA_URL) + ".zip")
file_cached(
    tmp_zip_path,
    data_readable=lambda: urllib.request.urlopen(DATA_URL),
    step_name="data zip",
)

# Extract the relevant .csv file from the zip
tmp_csv_path = TMP_DIR / (stable_filename(DATA_URL) + ".csv")


def unzip_data_csv():
    with zipfile.ZipFile(tmp_zip_path, mode="r") as unzipped_data:
        return unzipped_data.open(DATA_FILE_ID)


file_cached(tmp_csv_path, data_readable=unzip_data_csv, step_name="CSV data")

# Generate a sqlite database from the csv file
db_path = (Path(".") / "data.db").resolve()


def prepare_sqlite3_db(path):
    # Disable file locking, fsync, and journaling to improve performance since
    # if this import fails then the
    sqlite3_perf_header = """\
.output /dev/null
PRAGMA synchronous = OFF;
PRAGMA journal_mode = OFF;
PRAGMA locking_mode = EXCLUSIVE;
.output stdout
"""
    sqlite3_import = f"""{sqlite3_perf_header}
.mode csv
.separator '\t'
.import {tmp_csv_path} PUP_data
.mode list
SELECT "Imported " || COUNT(*) || " records" FROM PUP_data;
"""
    # Reference: datasette's handling of full-text search in sqlite https://docs.datasette.io/en/latest/full_text_search.html
    sqlite3_schema = f"""{sqlite3_perf_header}
-- CREATE UNIQUE INDEX index_npi ON PUP_data (npi, hcpcs_code, place_of_service);

CREATE VIRTUAL TABLE PUP_data_fts USING fts5(
    npi,
    hcpcs_code,
    hcpcs_description,
    name,
    address
);
INSERT INTO PUP_data_fts (rowid, npi, hcpcs_code, hcpcs_description, name, address)
    SELECT
        rowid,
        npi,
        hcpcs_code,
        hcpcs_description,
        nppes_provider_first_name || ' ' || nppes_provider_mi || ' ' || nppes_provider_last_org_name as name,
        nppes_provider_street1 || ' ' || nppes_provider_street2 || ' ' || nppes_provider_city || ' ' || nppes_provider_zip || ' ' || nppes_provider_state || ' ' || nppes_provider_country as address
    FROM PUP_data
;
"""
    sqlite3_optimize = f"""{sqlite3_perf_header}
INSERT INTO PUP_data_fts(PUP_data_fts) VALUES('optimize');
VACUUM;
"""

    def print_db_size():
        print("DB is now {} MiB".format(Path(path).stat().st_size / 1024 // 1024))

    print("========")
    print("= Importing CSV records")
    subprocess.run(["sqlite3", path], input=sqlite3_import.encode(), check=True)
    print_db_size()
    print("= Creating full-text search table")
    subprocess.run(["sqlite3", path], input=sqlite3_schema.encode(), check=True)
    print_db_size()
    if OPTIMIZE:
        print("= Optimizing DB using VACUUM and fts optimize")
        subprocess.run(["sqlite3", path], input=sqlite3_optimize.encode(), check=True)
        print_db_size()
    print("========")


file_cached(db_path, data_gen=prepare_sqlite3_db, step_name="sqlite3 DB")
