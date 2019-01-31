#!/bin/bash

#### phosphor.localhost ####

cat > openssl.cnf <<-EOF
  [req]
  distinguished_name = req_distinguished_name
  x509_extensions = v3_req
  prompt = no
  [req_distinguished_name]
  CN = *.phosphor.localhost
  [v3_req]
  keyUsage = keyEncipherment, dataEncipherment
  extendedKeyUsage = serverAuth
  subjectAltName = @alt_names
  [alt_names]
  DNS.1 = *.phosphor.localhost
  DNS.2 = phosphor.localhost
EOF

openssl req \
  -new \
  -newkey rsa:2048 \
  -sha1 \
  -days 3650 \
  -nodes \
  -x509 \
  -keyout phosphor.localhost.key \
  -out phosphor.localhost.crt \
  -config openssl.cnf

rm openssl.cnf

#### eos.localhost ####

cat > openssl.cnf <<-EOF
  [req]
  distinguished_name = req_distinguished_name
  x509_extensions = v3_req
  prompt = no
  [req_distinguished_name]
  CN = *.eos.localhost
  [v3_req]
  keyUsage = keyEncipherment, dataEncipherment
  extendedKeyUsage = serverAuth
  subjectAltName = @alt_names
  [alt_names]
  DNS.1 = *.eos.localhost
  DNS.2 = eos.localhost
EOF

openssl req \
  -new \
  -newkey rsa:2048 \
  -sha1 \
  -days 3650 \
  -nodes \
  -x509 \
  -keyout eos.localhost.key \
  -out eos.localhost.crt \
  -config openssl.cnf

rm openssl.cnf

# On OSX run the following
# open /Applications/Utilities/Keychain\ Access.app phosphor.localhost.crt
# open /Applications/Utilities/Keychain\ Access.app eos.localhost.crt
