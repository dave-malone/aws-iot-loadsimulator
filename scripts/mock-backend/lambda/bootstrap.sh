#!/usr/bin/env bash

PV="3.8"
if ! "pip${PV}" -V | grep -q "${PV}"; then
  echo "please install python/pip ${PV}"
  exit 1
fi
command -v virtualenv || pip${PV} install virtualenv

# crude bootstrap script to enable development in a virtualenv
virtualenv -p "$(command -v python${PV})" venv
source venv/bin/activate
pip${PV} install -r requirements.txt

echo "
to activate run:
    source venv/bin/activate
"
