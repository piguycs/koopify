THE DATASET IS ORIGINALLY SOURCED FROM: https://www.kaggle.com/datasets/adampq/pokemon-tcg-all-cards-1999-2023
LICENSES FROM THAT DATASET APPLY

## Why is this not in the migrations?

These scripts will add almost 600 products to the database, from a dataset which I do not own. I find it safer to do it this way. You have no real need of using this dataset locally, as I have already deployed it to <https://koopify.piguy.nl>

## How to use this dataset

1. Setup the venv
```sh
python3 -m venv .venv
source .venv/bin/activate # or activate.fish or activate.csh
pip install psycopg2
```
2. Set the `PGDB` enviornment variable to your postgres url
3. Insert the products in the database
```sh
# while still in the venv
python3 insert_products.py
```
4. Tag the products
```sh
# while still in the venv
python3 insert_categories.py
```
