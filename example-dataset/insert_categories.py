#!/usr/bin/env python3

import argparse
import csv
import getpass
import os
import re
import sys
import unicodedata

import psycopg2


def make_slug(name):
    if not name:
        return ""
    slug = name.lower()
    slug = unicodedata.normalize("NFKD", slug)
    slug = slug.encode("ascii", "ignore").decode("ascii")
    slug = re.sub(r"[^a-z0-9]+", "-", slug)
    slug = slug.strip("-")
    slug = re.sub(r"-+", "-", slug)
    return slug


def make_category_slug(category):
    if not category:
        return ""
    return make_slug(category)


def main():
    parser = argparse.ArgumentParser(description="Add categories to products from CSV")
    parser.add_argument(
        "--csv",
        required=True,
        help="Path to the CSV file produced by parse_pokemon_cards.py",
    )
    parser.add_argument(
        "--db", help="PostgreSQL connection string (will prompt if not provided)"
    )
    args = parser.parse_args()

    conn_str = (
        args.db
        or os.environ.get("PGDB")
        or input("PostgreSQL connection string: ").strip()
    )

    try:
        conn = psycopg2.connect(conn_str)
    except Exception as e:
        print(f"Failed to connect to database: {e}", file=sys.stderr)
        sys.exit(1)

    cursor = conn.cursor()

    with open(args.csv, newline="", encoding="utf-8") as f:
        reader = csv.DictReader(f)
        rows = list(reader)

    total_rows = len(rows)
    categories_seen = {}
    product_links = 0
    product_not_found = 0

    print(f"Processing {total_rows} rows...", file=sys.stderr)

    for i, row in enumerate(rows, 1):
        product_slug = make_slug(row.get("name", ""))
        category = row.get("category", "").strip()

        if not category:
            continue

        cat_slug = make_category_slug(category)

        if cat_slug not in categories_seen:
            cursor.execute("SELECT id FROM categories WHERE slug = %s", (cat_slug,))
            existing = cursor.fetchone()
            if existing:
                cat_id = existing[0]
            else:
                cursor.execute(
                    "INSERT INTO categories (name, slug) VALUES (%s, %s) RETURNING id",
                    (category, cat_slug),
                )
                cat_id = cursor.fetchone()[0]
            categories_seen[cat_slug] = cat_id
        else:
            cat_id = categories_seen[cat_slug]

        cursor.execute("SELECT id FROM products WHERE slug = %s", (product_slug,))
        product = cursor.fetchone()

        if not product:
            product_not_found += 1
            if product_not_found <= 5:
                print(f"Product not found: {product_slug}", file=sys.stderr)
            continue

        product_id = product[0]

        cursor.execute(
            """
            INSERT INTO product_categories (product_id, category_id)
            VALUES (%s, %s)
            ON CONFLICT DO NOTHING
            """,
            (product_id, cat_id),
        )
        if cursor.rowcount > 0:
            product_links += 1

        if i % 500 == 0 or i == total_rows:
            print(f"Progress: {i}/{total_rows}", end="\r", file=sys.stderr)

    conn.commit()
    cursor.close()
    conn.close()

    print(f"\nDone.")
    print(f"  Categories added: {len(categories_seen)}")
    print(f"  Product links created: {product_links}")
    if product_not_found > 0:
        print(f"  Products not found (skipped): {product_not_found}")


if __name__ == "__main__":
    main()
