#!/usr/bin/env bash

export HUMAN_LOG="true"

go build -o moves

./moves -zeb_root="/Users/dave/Desktop/zebedee-data/content/zebedee" \
    -create=true \
    -collection="test123" \
    -src="/economy/environmentalaccounts/articles/environmentaltaxes/2015-06-01" \
    -dest="/economy/environmentalaccounts/articles/environmentaltaxes/2015-06-01/test"


#./moves -zeb_root="/zebe-test" \
 #   -create=true \
#	-collection="test123" \
#	-src="/economy/economicoutputandproductivity/productivitymeasures/articles/experimentalestimatesofinvestmentinintangibleassetsintheuk2015" \
#	-dest="/economy/economicoutputandproductivity/productivitymeasures/articles/experimentalestimatesofinvestmentinintangibleassetsintheuk"