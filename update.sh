#!/bin/bash

csv_url="https://storage.googleapis.com/play_public/supported_devices.csv"

echo "Downloading supported_devices.csv from $csv_url"
csvLines=$(curl -s $csv_url)

# split csvLines into array
IFS=$'\n' read -d '' -r -a csvLines <<< "$csvLines"

# remove first line
unset csvLines[0]

# for each line in csvLines
for line in "${csvLines[@]}" ; do
    brand=$(echo $line | cut -d ',' -f 1)
    marketing_name=$(echo $line | cut -d ',' -f 2)
    model=$(echo $line | cut -d ',' -f 4)
    if [ "$brand" != "" ] && [ "$marketing_name" != "" ] && [ "$model" != "" ]; then
        #if not exist create brand folder
        if [ ! -d "$brand" ]; then
            mkdir $brand
        fi
        #if not exist create device file with marketing_name content
        if [ ! -f "$brand/$device" ]; then
            echo $marketing_name > $brand/$device
        fi
    fi
done

rm -f supported_devices.csv

git add .
# if there are changes
if [ -n "$(git status --porcelain)" ]; then
    git commit -m "Update devices on $(date +%Y-%m-%d)"
    git push origin master
fi
