#!/usr/bin/env sh
set -eu
base="${BASE_URL:-http://127.0.0.1:8080}"
curl -fsS "$base/healthz" >/dev/null
zone=$(curl -fsS -X POST "$base/api/v1/zones" -H "content-type: application/json" -d '{"name":"example.com."}')
id=$(printf "%s" "$zone" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
curl -fsS "$base/api/v1/zones/$id/validate" >/dev/null
printf "%s\n" "smoke ok zone=$id"

