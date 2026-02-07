#!/bin/sh

echo "GOGC=100で計測"
GOGC=100
go build
time ./prob-ch06-3

echo
echo "GOGC=80で計測"
GOGC=80
go build
time ./prob-ch06-3

echo
echo "GOGC=10で計測"
GOGC=10
go build
time ./prob-ch06-3
