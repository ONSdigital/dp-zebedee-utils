#!/usr/bin/env bash

export HUMAN_LOG=true

go build -o fix

./fix -col="pdf_test_do_not_publish" -type="article, bulletin, compendium_landing_page, compendium_chapter, static_methodology" -dir="/zebe-test"
#./fix -col="pdf_2048" -type="bulletin" -limit=20 -dir="/Users/dave/Desktop/zebedee-data/content/zebedee"
