#!/usr/bin/env bash

export HUMAN_LOG="true"

go build -o visualisations

./visualisations -zeb_root="/Users/carl/zebedee" \
    -reverse_changes=false \
    -collection="visualisationsGA"