#!/usr/bin/env sh 
#
# 
# This script is run before the tests are run. It is used to generate key files 


PEM="key.pem"
CERT="cert.pem"
TARGET_DIR="test/ssl/"
TEST_PASSWD="test"


if [ ! -f "${TARGET_DIR}${PEM}" ] || [ ! -f "${TARGET_DIR}${CERT}" ]; then 
    if [ ! -d "${TARGET_DIR}" ]; then 
        mkdir -p "${TARGET_DIR}"
    fi
    openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
    mv "${PEM}" "${TARGET_DIR}${PEM}"
    mv "${CERT}" "${TARGET_DIR}${CERT}"
fi






