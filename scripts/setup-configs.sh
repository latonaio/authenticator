#!/bin/sh
HOST_NAME="mysql"
PORT=xxxx
DB_USER_NAME=xxxx
DB_USER_PASSWORD=xxxx

SCRIPT_DIR=$(cd $(dirname $0); pwd)
TARGET_FILE="configs.yaml"

cd ${SCRIPT_DIR}
cd ..
make generate-key-pair

cd configs
sed -e "s/host_name:.*$/host_name: ${HOST_NAME}/g" \
    -e "s/PORT/${PORT}/g" \
    -e "s/user_name:.*$/user_name: ${DB_USER_NAME}/g" \
    -e "s/user_password:.*$/user_password: ${DB_USER_PASSWORD}/g" \
    -e "s/private_key:.*$/private_key: |\n/g" \
    ${TARGET_FILE} > temp.yml

indentStdin() {
  indent='  '
  while IFS= read -r line; do
      echo "${indent}${line}"
  done
  echo
}

cat ../private.key | indentStdin >> temp.yml

rm ${TARGET_FILE}
mv temp.yml ${TARGET_FILE}