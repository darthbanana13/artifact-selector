#!/usr/bin/env bash

INPUT_FILE="assets/extension_list.txt"
OUTPUT_FILE="pkg/extensionlist/map.go"

if [[ ! -f "$INPUT_FILE" ]]; then
    echo "Error: Extension list file '${INPUT_FILE}' not found."
    exit 1
fi

cat > "${OUTPUT_FILE}" <<EOL
package extensionlist

// exts holds a set of known file extensions, generated from ${INPUT_FILE}.
var exts = map[string]bool{
EOL

while IFS= read -r line || [[ -n "$line" ]]; do
    if [[ -n "$line" ]]; then
        printf "\t\"%s\": true,\n" "$line" >> "$OUTPUT_FILE"
    fi
done < "$INPUT_FILE"

echo "}" >> "$OUTPUT_FILE"

echo "Successfully generated ${OUTPUT_FILE}"
