#!/bin/bash - 
#===============================================================================
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o  host_ages .
  5  docker build -f Dockerfile -t host_ages:latest .
  6  docker run -itd host_ages:latest /host_ages
