@echo off
rem Кросскомпиляция из под Windows для Linux

set GOOS=linux
set GOARCH=amd64
go build -o app
