#!/usr/bin/env bash

export HUMAN_LOG="true"

go build -o fixxxer

#./fixxxer -master="/Users/dave/Desktop/zebedee-data/content/zebedee/master"
./fixxxer -master="/zebe-test/master"