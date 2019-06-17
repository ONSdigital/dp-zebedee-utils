#!/usr/bin/env bash

export HUMAN_LOG=true

go build -o counter

#./counter -types="article, bulletin, compendium_landing_page, compendium_chapter, static_methodology" -dir="/zebe-test/master"
./counter -types="article, bulletin, compendium_landing_page, compendium_chapter, static_methodology" -dir="/Users/dave/Desktop/zebedee-data/content/zebedee/master"
