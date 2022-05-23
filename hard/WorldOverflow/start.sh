#!/bin/bash

docker build -t world_overflow .
docker run -i -dp 8888:8888 world_overflow
