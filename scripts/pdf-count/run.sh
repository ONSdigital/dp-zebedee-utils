#!/usr/bin/env bash

export HUMAN_LOG=true

go build -o fix

./fix -dir="/zebe-test"
#./fix -dir="/Users/dave/Desktop/zebedee-data/content/zebedee/master"
