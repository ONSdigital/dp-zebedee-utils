#!/usr/bin/env bash

export HUMAN_LOG="true"

go build -o moves

./moves -zeb_root="/zeb-test" -collection="test123" -src="/aboutus/data.json" -dest="/aboutus/test/data.json"