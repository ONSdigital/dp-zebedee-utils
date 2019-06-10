#!/usr/bin/env bash

go build -o gsiFix

export HUMAN_LOG=1
./gsiFix -collections="/Users/dave/Desktop/zebedee-data/content/zebedee/collections" \
-master="/Users/dave/Desktop/zebedee-data/content/zebedee/master"

