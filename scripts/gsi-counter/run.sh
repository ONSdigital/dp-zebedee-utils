#!/usr/bin/env bash

export HUMAN_LOG=true

go build -o counter

./counter  -type="timeseries" -dir="/Users/dave/Desktop/zebedee-data/content/zebedee"
