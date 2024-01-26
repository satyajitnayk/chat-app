#!/bin/bash

echo "creating server.key"
openssl genrsa -out server.key 2048
openssl ecparam -genkey -name secp384r1 -out server.key
echo "creating server.cert"
openssl req -new -x509 -sha256 -key server.key -out server.cert -batch -days 365