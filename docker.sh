#!/bin/bash
docker stop socketiosample && docker container rm socketiosample
cd /opt/socketiosample
docker build --no-cache -t socketiosample .
docker run -d -p 12379:12379 --name=chitchat-push -ti -v /chitchat-data:/chitchat-data chitchat-push