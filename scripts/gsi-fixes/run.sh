#!/usr/bin/env bash

export HUMAN_LOG=true

go build -o fix

./fix -col="test123" -dir="/Users/dave/Desktop/zebedee-data/content/zebedee"
