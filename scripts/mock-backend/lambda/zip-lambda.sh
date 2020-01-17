#!/usr/bin/env bash

RUN_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ACCOUNT_ID=$(aws --output text sts get-caller-identity --query 'Account')
STACK_NAME=aws-iot-loadsimulator-mockbackend
S3_BUCKET="${STACK_NAME}-${ACCOUNT_ID}"

mapfile -t FUNCTION_DIRS < <(find "${RUN_DIR}" -mindepth 1 -maxdepth 1 -type d ! -name venv | sort)
PV="3.8"

for dir in "${FUNCTION_DIRS[@]}"; do
    pushd "${dir}" || exit 127
    rm -rf *.zip venv/
    function_name=$(basename "${dir}")
    zip_file_name="${function_name}.zip"
    if [[ -e requirements.txt ]]; then
        virtualenv -p "$(command -v python${PV})" venv
        source ./venv/bin/activate
        ./venv/bin/pip${PV} --no-cache-dir install -r requirements.txt
        func_dir=$(pwd)
        mapfile -t sp_dirs < <(find . -type d -name site-packages -o -name dist-packages)
        for sp_dir in "${sp_dirs[@]}"; do
            pushd "${sp_dir}" || exit 127
            zip -r9 "${func_dir}"/"${zip_file_name}" .
            popd || exit 128
        done
        deactivate
        cd "${func_dir}" || exit 129
    fi
    zip -ur "${zip_file_name}" *.py
    aws s3 cp "${zip_file_name}" s3://"${S3_BUCKET}"/"${zip_file_name}"
    # aws lambda update-function-code --function-name "${function_name}" --s3-bucket "${S3_BUCKET}" --s3-key "${zip_file_name}"
    popd || exit 128
done
