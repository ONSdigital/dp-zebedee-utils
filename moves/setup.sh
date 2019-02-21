#!/usr/bin/env bash

ECHO "making dir"

apt-get update && apt-get install vim

cd src && mkdir -p github.com/ONSdigital
cd github.com/ONSdigital

ECHO "clone Go script"
git clone -b feature/content-mover https://github.com/ONSdigital/dp-zebedee-utils.git

ECHO "getting Go dependencies"
go get github.com/satori/go.uuid
go get github.com/ONSdigital/log.go/log

cd dp-zebedee-utils/moves
ECHO "ready to go"


