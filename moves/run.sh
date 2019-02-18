#!/usr/bin/env bash

export HUMAN_LOG="true"

go build -o moves

./moves -zeb_root="/Users/dave/Desktop/zebedee-data/content/zebedee" -collection="test123" -action="mk"
#./moves -zeb_root="/Users/dave/Desktop/zebedee-data/content/zebedee" -collection="test123" -action="del"