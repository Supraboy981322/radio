#!/bin/bash

config="$(cat \
  ~/.config/Supraboy981322/radio/tui/config.toml)"
server=$(printf "${config}" \
  | grep '^radio = ' \
  | awk -F'=' '{print $2}' \
  | tr -d ' '\")

libraryJSON=$(gum spin \
  -s minidot \
  --title "fetching library..." \
  -- curl ${server}library.json -s)

libraryLength=$((\
  $(echo "${libraryJSON}" \
  | jq 'length')-1 \
))
for item in $(seq ${libraryLength}); do
  stationNames+="$(echo "${libraryJSON}" \
    | jq -r .[${item}].[0])\n"
done

selectedStation=$(printf "${stationNames}" \
  | gum filter)
echo "${selectedStation}"
selectedStation=$(echo "${libraryJSON}" \
  | jq -r --arg v "${selectedStation}" '
  map(select(any(.[]; . == $v)))[0][1]
')
desc=$(echo "${libraryJSON}" \
  | jq -r --arg v "${selectedStation}" '
  map(select(any(.[]; . == $v)))[0][2]
')
echo "${desc}"
echo "${selectedStation}" 

gum spin \
  -s points \
  --title "playing..." \
  -- ffplay ${selectedStation} \
  -nodisp \
  -autoexit
