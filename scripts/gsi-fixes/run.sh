#!/usr/bin/env bash

export HUMAN_LOG=true

go build -o fix

./fix -col="gsi_timeseries" -dir="/Users/dave/Desktop/zebedee-data/content/zebedee"
