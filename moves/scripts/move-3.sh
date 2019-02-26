#!/usr/bin/env bash

export HUMAN_LOG="true"

go build -o moves

../moves -zeb_root="/zebe-test" \
    -create=true \
	-collection="move3-onsworkingpaperseries" \
	-src="/methodology/methodologicalpublications/generalmethodology/onsworkingpaperseries/onsmethodologyworkingpaperseriesnumber16syntheticdatapilot/onsworkingpaperseriesno17usingdatasciencefortheaddressmatchingservice" \
	-dest="/methodology/methodologicalpublications/generalmethodology/onsworkingpaperseries/onsworkingpaperseriesno17usingdatasciencefortheaddressmatchingservice"



