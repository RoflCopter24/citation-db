#!/bin/bash

HOST_VOLUME_PATH=/srv/mongodb
CONTAINER_BIND_IP=127.0.0.1
CONTAINER_PORT_MONGO=27017

docker run --name=citation-db-mongodb -h citation-db-mongodb -p ${CONTAINER_BIND_IP}:${CONTAINER_PORT_MONGO}:27017 -d mongo
