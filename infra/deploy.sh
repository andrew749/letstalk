#!/bin/bash

set -e

./build_prod.sh
./build_frontend.sh
./prod.sh
