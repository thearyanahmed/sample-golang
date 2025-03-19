#!/bin/bash

# Cloudflare API endpoint
API_ENDPOINT="https://api.cloudflare.com/client/v4/zones/$CF_ZONE_ID/custom_hostnames/$CF_HOSTNAME_ID"

# Cloudflare authentication headers
AUTH_EMAIL="$CF_AUTH_EMAIL"
AUTH_KEY="$CF_AUTH_KEY"
#
# Check if required environment variables are set
if [[ -z "$CF_ZONE_ID" || -z "$CF_HOSTNAME_ID" || -z "$CF_AUTH_EMAIL" || -z "$CF_AUTH_KEY" ]]; then
  echo "Error: Missing required environment variables."
  echo "Please set CF_ZONE_ID, CF_HOSTNAME_ID, CF_AUTH_EMAIL, CF_AUTH_KEY"
  exit 1
fi


curl --request PATCH \
"https://api.cloudflare.com/client/v4/zones/$CF_ZONE_ID/custom_hostnames/$CF_HOSTNAME_ID" \
--header "X-Auth-Email: $AUTH_EMAIL" \
--header "X-Auth-Key: $AUTH_KEY" \
--header "Content-Type: application/json" \
--data '{
  "ssl": {
    "method": "http",
    "type": "dv"
  },
  "custom_metadata": {
    "cf_cache" : "disabled",
    "email_obfuscation" : "disabled",
    "some_other_data" : "helloworld1"
  }
}'

