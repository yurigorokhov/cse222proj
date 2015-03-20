#!/usr/bin/env bash

apt-get update && apt-get install -y docker.io

# Setup docker
ln -s /vagrant /opt/docker
cd /opt/docker
#docker build -t yurigorokhov/tcp_bench .
#docker run -i -t -p 8081:80 -p 11111:11111 yurigorokhov/tcp_bench
