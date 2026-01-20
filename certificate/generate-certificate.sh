#!/bin/bash

#--- Generate CA key & certificate
openssl req -newkey rsa:2048 -nodes -days 9999 -x509 -keyout ca.key -out ca.crt -subj "/CN=*"

#--- Generate agent key & cert request
openssl req -newkey rsa:2048 -nodes -keyout agent.key -out agent.csr -subj "/C=EU/ST=Italy/L=Rome/O=Algotiqa/OU=AlgotiqaAgent/CN=*"

#--- Sign agent certificate with CA
openssl x509 -req -days 9999 -sha256 -in agent.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out agent.crt \
    -extfile <(echo subjectAltName=IP:158.69.26.216)
