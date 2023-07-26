@echo off

start go run ./hertz-server/
cd ./kitex_service
start go run .
cd ../kitex_service_2
start go run .

