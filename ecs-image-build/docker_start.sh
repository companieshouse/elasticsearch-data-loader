#!/usr/bin/env bash

set -euo pipefail

SEARCH="${SEARCH:-${search:-company}}"
INDEX="${INDEX:-${index:-}}"
ES_URL="${ES_URL:-${es_url:-}}"
MONGO_URL="${MONGO_URL:-${mongo_url:-}}"
ALPHAKEY_URL="${ALPHAKEY_URL:-${alphakey_url:-}}"
USERNAME="${USERNAME:-${username:-}}"
PASSWORD="${PASSWORD:-${password:-}}"
CREATE_MAPPING="${CREATE_MAPPING:-${create_mapping:-false}}"
COMPANY_LIMIT="${COMPANY_LIMIT:-${company_limit:-0}}"

RUN_SCRIPT="${RUN_SCRIPT:-/opt/run-elastic-search.sh}"


# Fallback for zipped build layouts where app files may live under /opt/build-*/.
# Avoid relying on external tools (for example, find) that may not exist in slim images.
if [ ! -x "$RUN_SCRIPT" ]; then
  for candidate in \
    /opt/run-elastic-search.sh \
    /opt/*/run-elastic-search.sh \
    /opt/*/*/run-elastic-search.sh \
    /opt/*/*/*/run-elastic-search.sh \
    /opt/*/*/*/*/run-elastic-search.sh
  do
    if [ -f "$candidate" ]; then
      chmod +x "$candidate" 2>/dev/null || true
      if [ -x "$candidate" ]; then
        RUN_SCRIPT="$candidate"
        break
      fi
    fi
  done
fi

if [ ! -x "$RUN_SCRIPT" ]; then
  echo "ERROR: Cannot execute loader script: $RUN_SCRIPT"
  exit 1
fi

if [ -z "$INDEX" ] || [ -z "$ES_URL" ] || [ -z "$MONGO_URL" ] || [ -z "$ALPHAKEY_URL" ]; then
  echo "ERROR: Missing required env vars. Required: INDEX, ES_URL, MONGO_URL, ALPHAKEY_URL"
  exit 1
fi

if [ "$CREATE_MAPPING" != "true" ] && [ "$CREATE_MAPPING" != "false" ]; then
  echo "ERROR: CREATE_MAPPING must be 'true' or 'false'"
  exit 1
fi

if ! [[ "$COMPANY_LIMIT" =~ ^[0-9]+$ ]]; then
  echo "ERROR: COMPANY_LIMIT must be a non-negative integer"
  exit 1
fi

if [ -n "$USERNAME" ] && [ -z "$PASSWORD" ]; then
  echo "ERROR: PASSWORD must be set when USERNAME is set"
  exit 1
fi

if [ -n "$PASSWORD" ] && [ -z "$USERNAME" ]; then
  echo "ERROR: USERNAME must be set when PASSWORD is set"
  exit 1
fi

cmd=(
  "$RUN_SCRIPT"
  -s "$SEARCH"
  -i "$INDEX"
  -e "$ES_URL"
  -m "$MONGO_URL"
  -a "$ALPHAKEY_URL"
  -c "$CREATE_MAPPING"
  -l "$COMPANY_LIMIT"
)

if [ -n "$USERNAME" ]; then
  cmd+=( -u "$USERNAME" -p "$PASSWORD" )
fi

echo "Starting loader for index '$INDEX' with search '$SEARCH'"
echo "Using create_mapping=$CREATE_MAPPING company_limit=$COMPANY_LIMIT"

exec "${cmd[@]}"
