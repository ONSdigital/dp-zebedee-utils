#!/usr/bin/env bash

export HUMAN_LOG=true

go build -o fix

./fix -col="gsi_datasets" -type="dataset" -dir="/zebe-test"
