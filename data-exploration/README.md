# What Did Medicare Pay?: Data Exploration

This subfolder contains Jupyter notebooks for exploring the Medicare Provider Utilization & Payments data.

## Setup

[Poetry](https://python-poetry.org) is used to manage dependencies for this exploration. Data is built on the `data/` subfolder.

### Environment

1. Run `poetry install` to download, build, and setup all dependencies (like Jupyter).
2. Run `poetry run jupyter notebook` from this subfolder to load the Jupyter environment.

### Data

1. Make sure that a DB has been produced in the `data/` subfolder. See the instructions their for how that is handled.
2. Use `notebook/prepare_data_exp.ipynb` to convert the outputted DB in `data/` to something more optimized for exploration.
