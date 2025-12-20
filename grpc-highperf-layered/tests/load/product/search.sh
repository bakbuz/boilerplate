#!/bin/bash

# Ensure ghz is installed: https://ghz.sh/docs/install

# Run load test for ProductService/Search
ghz --insecure \
  --proto ../../../api/protos/catalog/v1/product.proto \
  --call catalog.v1.ProductService/Search \
  -d '{"name":"Product", "limit":10}' \
  -c 10 \
  -n 1000 \
  127.0.0.1:50051
