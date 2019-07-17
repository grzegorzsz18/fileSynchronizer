# !/bin/bash

docker build  --no-cache --tag=test ../
docker container run -p 18080:18080 -p 22222:22222 test:latest
