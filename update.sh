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
    device=$(echo $line | cut -d ',' -f 3)
    model=$(echo $line | cut -d ',' -f 4)
    if [ "$brand" != "" ] && [ "$marketing_name" != "" ] && [ "$device" != "" ] && [ "$model" != "" ] ; then
        if [ ! -d "$brand" ]; then
            mkdir "$brand"
        fi

        if [[ "$device" == *"/"* ]]; then
            device=$(echo "$device" | sed 's/\//-/g')
        fi
        if [ ! -f "$brand/$device" ]; then
            echo "$marketing_name" > "$brand/$device"
        fi

        if [[ "$model" == *"/"* ]]; then
            model=$(echo "$model" | sed 's/\//-/g')
        fi
        if [ ! -f "$brand/$model" ]; then
            echo "$marketing_name" > "$brand/$model"
        fi
    fi
done

rm -f supported_devices.csv
