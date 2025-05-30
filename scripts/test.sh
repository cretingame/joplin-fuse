#!/bin/bash

ADDRESS="http://127.0.0.1:41184"

PAYLOAD=$(
  cat <<EOF
{"key1":"value1", "key2":"value2"}
EOF
)

post_example() {
  curl -X POST -H "Content-Type: application/json" -d "$PAYLOAD" "$ADDRESS"
}

# https://joplinapp.org/fr/help/dev/spec/clipper_auth
# return
# {"auth_token":"EvEiLZa0NRxKhaV6d1EoQA"}
get_auth_token() {
  curl -X POST "$ADDRESS/auth" | jq '.auth_token' | sed 's/\"//g'
}

# https://joplinapp.org/fr/help/dev/spec/clipper_auth
# save the token in then token file in the current directory
save_token() {
  RESP=$(curl -X POST "$ADDRESS/auth")
  AUTH_TOKEN=$(echo "$RESP" | jq '.auth_token' | sed 's/\"//g')
  read -p "Check Joplin and Press any key to continue ..."
  RESP=$(curl "$ADDRESS/auth/check?auth_token=$AUTH_TOKEN")
  STATUS=$(echo "$RESP" | jq '.status' | sed 's/\"//g')
  if [[ "$STATUS" = "accepted" ]]; then
    TOKEN=$(echo "$RESP" | jq '.token' | sed 's/\"//g')
    echo ok
    echo "$TOKEN" >"$(dirname "$0")/token"
  else
    echo ko
  fi
}

# get the token from the token file in the current directory
get_token() {
  cat "$(dirname "$0")/token"
}

list_folders() {
  local TOKEN RESP PAGE HAS_MORE ITEM I
  HAS_MORE="true"
  PAGE=0
  TOKEN=$(get_token)
  while [ "$HAS_MORE" = "true" ]; do
    echo "page $PAGE"
    RESP=$(curl "http://localhost:41184/folders?token=$TOKEN&page=$PAGE")
    ITEM="undefined"
    I=0
    while ! [ "$ITEM" = "null" ]; do
      ITEM=$(echo "$RESP" | jq ".items[$I]")
      echo "$ITEM"
      I=$((I + 1))
    done
    HAS_MORE=$(echo "$RESP" | jq '.has_more')
    PAGE=$((PAGE + 1))
  done
}

list_resources() {
  local TOKEN RESP PAGE HAS_MORE ITEM I
  HAS_MORE="true"
  PAGE=0
  TOKEN=$(get_token)
  while [ "$HAS_MORE" = "true" ]; do
    echo "page $PAGE"
    RESP=$(curl "http://localhost:41184/resources?token=$TOKEN&page=$PAGE")
    ITEM="undefined"
    I=0
    while ! [ "$ITEM" = "null" ]; do
      ITEM=$(echo "$RESP" | jq ".items[$I]")
      echo "$ITEM"
      I=$((I + 1))
    done
    HAS_MORE=$(echo "$RESP" | jq '.has_more')
    PAGE=$((PAGE + 1))
  done
}

# list_folders
# list_resources

ID=73ca02a59f2741bcb98a77516b85b9a5
curl "http://localhost:41184/resources/$ID/file?token=$(get_token)"
