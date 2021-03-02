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
.import {tmp_csv_path} imported_data
.mode list
SELECT "Imported " || COUNT(*) || " records" FROM imported_data;
"""
    # Reference: datasette's handling of full-text search in sqlite https://docs.datasette.io/en/latest/full_text_search.html
    sqlite3_schema = f"""{sqlite3_perf_header}
-- Create main PUP_data table
CREATE TABLE PUP_data AS
    SELECT
        -- rowid # assigned contiguously ascending values, starting with 1, in the order that they are returned by the SELECT statement
        CAST(npi AS TEXT) as npi,
        CAST(nppes_provider_last_org_name AS TEXT) as nppes_provider_last_org_name,
        CAST(nppes_provider_first_name AS TEXT) as nppes_provider_first_name,
        CAST(nppes_provider_mi AS TEXT) as nppes_provider_mi,
        CAST(nppes_credentials AS TEXT) as nppes_credentials,
        CAST(nppes_provider_gender AS TEXT) as nppes_provider_gender,
        CAST(nppes_entity_code AS TEXT) as nppes_entity_code,
        CAST(nppes_provider_street1 AS TEXT) as nppes_provider_street1,
        CAST(nppes_provider_street2 AS TEXT) as nppes_provider_street2,
        CAST(nppes_provider_city AS TEXT) as nppes_provider_city,
        CAST(nppes_provider_zip AS TEXT) as nppes_provider_zip,
        CAST(nppes_provider_state AS TEXT) as nppes_provider_state,
        CAST(nppes_provider_country AS TEXT) as nppes_provider_country,
        CAST(provider_type AS TEXT) as provider_type,
        CAST(medicare_participation_indicator AS TEXT) as medicare_participation_indicator,
        CAST(place_of_service AS TEXT) as place_of_service,
        CAST(hcpcs_code AS TEXT) as hcpcs_code,
        -- CAST(hcpcs_description AS TEXT) as hcpcs_description, # NOTE: removing this field to reduce DB size; use hcpcs table instead.
        CAST(hcpcs_drug_indicator AS TEXT) as hcpcs_drug_indicator,
        CAST(line_srvc_cnt AS NUMERIC) as line_srvc_cnt,
        CAST(bene_unique_cnt AS NUMERIC) as bene_unique_cnt,
        CAST(bene_day_srvc_cnt AS NUMERIC) as bene_day_srvc_cnt,
        CAST(average_Medicare_allowed_amt AS REAL) as average_Medicare_allowed_amt,
        CAST(average_submitted_chrg_amt AS REAL) as average_submitted_chrg_amt,
        CAST(average_Medicare_payment_amt AS REAL) as average_Medicare_payment_amt,
        CAST(average_Medicare_standard_amt AS REAL) as average_Medicare_standard_amt
    FROM imported_data
    WHERE npi != "0000000001"; -- copyright details are saved as a malformed first row
;

-- Create auxiliary hcpcs table
CREATE TABLE hcpcs AS
    SELECT
        CAST(hcpcs_code AS TEXT) as code,
        CAST(hcpcs_description AS TEXT) as description
    FROM imported_data
    GROUP BY hcpcs_code, hcpcs_description
;
CREATE INDEX index_hcpcs_code ON hcpcs(code);

-- Drop transient import table
DROP TABLE imported_data;
"""
    sqlite3_fts = f"""{sqlite3_perf_header}
-- Create full-text search table
CREATE VIRTUAL TABLE PUP_data_fts USING fts5(
    npi,
    hcpcs_code,
    hcpcs_description,
    name,
    address
);
INSERT INTO PUP_data_fts (rowid, npi, hcpcs_code, hcpcs_description, name, address)
    SELECT
        p.rowid,
        p.npi,
        p.hcpcs_code,
        h.description,
        p.nppes_provider_first_name || ' ' || p.nppes_provider_mi || ' ' || p.nppes_provider_last_org_name as name,
        p.nppes_provider_street1 || ' ' || p.nppes_provider_street2 || ' ' || p.nppes_provider_city || ' ' || p.nppes_provider_zip || ' ' || p.nppes_provider_state || ' ' || p.nppes_provider_country as address
    FROM PUP_data AS p
    INNER JOIN hcpcs AS h ON p.hcpcs_code = h.code
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
    print("= Transforming into structured schema")
    subprocess.run(["sqlite3", path], input=sqlite3_schema.encode(), check=True)
    print_db_size()
    print("= Creating full-text search table")
    subprocess.run(["sqlite3", path], input=sqlite3_fts.encode(), check=True)
    print_db_size()
    if OPTIMIZE:
        print("= Optimizing DB using VACUUM and fts optimize")
        subprocess.run(["sqlite3", path], input=sqlite3_optimize.encode(), check=True)
        print_db_size()
    print("========")


file_cached(db_path, data_gen=prepare_sqlite3_db, step_name="sqlite3 DB")
