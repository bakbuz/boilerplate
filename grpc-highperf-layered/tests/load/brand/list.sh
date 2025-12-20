#!/bin/bash

# Ensure ghz is installed: https://ghz.sh/docs/install

# Run load test for BrandService/List
ghz --insecure \
  --proto ../../../api/protos/catalog/v1/brand.proto \
  --call catalog.v1.BrandService/List \
  -c 10 \
  -n 1000 \
  127.0.0.1:50051
