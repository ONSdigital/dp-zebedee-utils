#!/usr/bin/env bash

export HUMAN_LOG="true"

go build -o moves

ECHO "executing move 2 - augusy ...."

./moves -zeb_root="/zebe-test" \
    -create=true \
	-collection="move2-developingnewmeasuresofinfrastructureinvestment" \
	-src="/economy/economicoutputandproductivity/productivitymeasures/articles/developingnewmeasuresofinfrastructureinvestment/augusy2018" \
	-dest="/economy/economicoutputandproductivity/productivitymeasures/articles/developingnewmeasuresofinfrastructureinvestment/august2018"



