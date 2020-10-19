#!/usr/bin/env bash
set -e

CATALOG_PATH=docs/catalog/v1

CATALOG_INDEX_FILE=$CATALOG_PATH/index.json

rm $CATALOG_INDEX_FILE

echo "[" > $CATALOG_INDEX_FILE

for templateFile in $(find $CATALOG_PATH -type f -name 'index.json'); do
    if [ "docs/catalog/v1/index.json" != $templateFile ]; then
        cat $templateFile >> $CATALOG_INDEX_FILE
        echo "," >> $CATALOG_INDEX_FILE
    fi
done

sed -i '$ s/.$//' $CATALOG_INDEX_FILE
echo "]" >> $CATALOG_INDEX_FILE