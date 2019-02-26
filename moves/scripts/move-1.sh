#!/usr/bin/env bash

export HUMAN_LOG="true"

go build -o moves

ECHO "executing move 1 - Experimental estimates...."

../moves -zeb_root="/zebe-test" \
    -create=true \
	-collection="move1_experimentalEstimates" \
	-src="/economy/economicoutputandproductivity/productivitymeasures/articles/experimentalestimatesofinvestmentinintangibleassetsintheuk2015" \
	-dest="/economy/economicoutputandproductivity/productivitymeasures/articles/experimentalestimatesofinvestmentinintangibleassetsintheuk"
