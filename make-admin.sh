#!/usr/bin/env bash

DB_URL="$1"
COMMAND="$2"
ARG="$3"

if [ -z "$DB_URL" ]; then
  echo "Error: missing database URL."
  echo "Usage: $0 <db_url> <command> [args]"
  exit 1
fi

function list_users() {
  echo "User list:"
  psql "$DB_URL" -c \
    "SELECT id, display_name, email, admin FROM users ORDER BY id;"
}

function make_admin() {
  USER_ID="$1"

  if [ -z "$USER_ID" ]; then
    read -p "Enter user ID to promote to admin: " USER_ID
  fi

  # basic validation
  if ! [[ "$USER_ID" =~ ^[0-9]+$ ]]; then
    echo "Invalid ID"
    return 1
  fi

  psql "$DB_URL" -c \
    "UPDATE users SET admin = true WHERE id = $USER_ID;"

  echo "User with ID $USER_ID is now an admin."
}

case "$COMMAND" in
  list)
    list_users
    ;;
  make-admin)
    make_admin "$ARG"
    ;;
  *)
    echo "Usage:"
    echo "  $0 <db_url> list"
    echo "  $0 <db_url> make-admin <user_id>"
    exit 1
    ;;
esac
