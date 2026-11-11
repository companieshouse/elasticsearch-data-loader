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
USE_AWS_SIGV4="${USE_AWS_SIGV4:-${use_aws_sigv4:-false}}"

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

# When alphakey runs as an ECS sidecar it starts in parallel with this container,
# so wait for it to accept connections before the loader makes its first request.
ALPHAKEY_WAIT_SECONDS="${ALPHAKEY_WAIT_SECONDS:-120}"
if [[ "$ALPHAKEY_URL" =~ ^http://(localhost|127\.0\.0\.1):([0-9]+) ]]; then
  alphakey_host="${BASH_REMATCH[1]}"
  alphakey_port="${BASH_REMATCH[2]}"
  echo "Waiting up to ${ALPHAKEY_WAIT_SECONDS}s for alphakey sidecar on ${alphakey_host}:${alphakey_port}"
  waited=0
  until (exec 3<>"/dev/tcp/${alphakey_host}/${alphakey_port}") 2>/dev/null; do
    if [ "$waited" -ge "$ALPHAKEY_WAIT_SECONDS" ]; then
      echo "ERROR: alphakey sidecar not reachable at ${alphakey_host}:${alphakey_port} after ${ALPHAKEY_WAIT_SECONDS}s"
      exit 1
    fi
    sleep 1
    waited=$((waited + 1))
  done
  echo "alphakey sidecar reachable after ${waited}s"
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

# Run from the loader script directory so its relative file paths resolve.
script_dir="$(dirname "$RUN_SCRIPT")"
cd "$script_dir"

# Export AWS SigV4 env var for HTTP client
export USE_AWS_SIGV4

# Ensure errors directory exists for log file
mkdir -p /opt/errors

# Filter bulk JSON to file, keep status/error messages visible in CloudWatch
"${cmd[@]}" 2>&1 | tee >(grep -v '^{' >&2) > /opt/errors/bulk_operations.log
exit_code=${PIPESTATUS[0]}

if [ $exit_code -eq 0 ]; then
  echo "Load completed successfully. Sleeping 5 minutes before exit..."
  sleep 300
fi

  echo "Setting desired count of ECS service ${ECS_SERVICE} to 0"
  aws ecs update-service \
    --cluster "${ECS_CLUSTER}" \
    --service "${ECS_SERVICE}" \
    --desired-count 0 \
    --region "${AWS_REGION}"

exit $exit_code
