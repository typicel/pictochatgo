#!/bin/bash

cd client
docker build -t client .
cd ..

cd server
docker build -t server .
cd ..

docker-compose up
