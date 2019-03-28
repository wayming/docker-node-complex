#!/bin/sh

ApiServerIP=`docker ps | grep complex_api | awk '{print $1}' | xargs -I {} docker inspect --format '{{ .NetworkSettings.Networks.complex_default.IPAddress }}' {}
`
curl -d "{\"index\":\"7\"}" -H "Content-Type: application/json" --url  http://${ApiServerIP}:5000/values
#repeat 1000 "curl -d \"{\\\"index\\\":\\\"7\\\"}\" -H \"Content-Type: application/json\" --url  http://172.19.0.2:5000/values"
