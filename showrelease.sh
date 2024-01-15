#!/bin/sh
TTT=/tmp/rel.dat
curl -s https://api.github.com/repos/ericdebeij/akamai-review/releases >"$TTT"
jq -r '.[0].tag_name' <"$TTT"
jq -r '.[0].name' <"$TTT"
