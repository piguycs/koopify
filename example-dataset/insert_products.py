#!/usr/bin/env python3

import argparse
import csv
import getpass
import os
import re
import sys
import unicodedata

import psycopg2
from psycopg2.extensions import AsIs


RARITY_LEVELS = {
    "common": 1,
    "uncommon": 2,
    "rare": 3,
    "rare holo": 4,
    "rare ultra": 5,
    "rare shiny": 4,
    "rare rainbow": 5,
    "rare secret": 6,
    "rare holo v": 5,
    "rare ultra holo": 6,
    "legend": 3,
    "mythic": 4,
}


def get_price_eur_cents(rarity, evolution_stage):
    rarity_key = rarity.lower().strip() if rarity else "common"
    rarity_level = RARITY_LEVELS.get(rarity_key, 1)
    return (rarity_level * 200) + (int(evolution_stage) * 50)


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


def main():
    parser = argparse.ArgumentParser(
        description="Insert Pokemon cards into products table"
    )
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

    total = len(rows)
    inserted = 0
    skipped = 0

    for i, row in enumerate(rows, 1):
        name = row.get("name", "").strip()
        category = row.get("category", "").strip()
        evolution_stage = int(row.get("evolution_stage") or 1)
        rarity = row.get("rarity", "").strip()
        description = row.get("description", "").strip()
        image_url = row.get("image_url", "").strip()

        slug = make_slug(name)
        price_eur_cents = get_price_eur_cents(rarity, evolution_stage)

        sql = """
            INSERT INTO products (
                name, slug, description, image_url,
                price_eur_cents, discount_percent,
                inventory_count, in_stock, is_active
            )
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)
            ON CONFLICT (slug) DO NOTHING
        """
        vals = (
            name,
            slug,
            description,
            image_url or None,
            price_eur_cents,
            0,
            10,
            True,
            True,
        )

        try:
            cursor.execute(sql, vals)
            if cursor.rowcount > 0:
                inserted += 1
            else:
                skipped += 1
        except Exception as e:
            print(f"Error inserting {name}: {e}", file=sys.stderr)
            skipped += 1

        if i % 100 == 0 or i == total:
            print(f"Progress: {i}/{total}", end="\r", file=sys.stderr)

    conn.commit()
    cursor.close()
    conn.close()

    print(
        f"\nDone. Inserted: {inserted}, Skipped (conflict): {skipped}, Total: {total}"
    )


if __name__ == "__main__":
    main()
