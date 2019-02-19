#!/usr/bin/env bash

export HUMAN_LOG="true"

go build -o moves

./moves -zeb_root="/Users/dave/Desktop/zebedee-data/content/zebedee" -collection="test123" -src="/aboutus/data.json" -dest="/aboutus/test/data.json"