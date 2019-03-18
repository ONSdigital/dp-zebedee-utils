#!/usr/bin/env bash

export HUMAN_LOG="true"

go build -o fixxxer

#./fixxxer -master="/Users/dave/Desktop/zebedee-data/content/zebedee/master" -collections="/Users/dave/Desktop/zebedee-data/content/zebedee/collections"
./fixxxer -master="/zebe-test/master" -collections="/zebe-test/collections"