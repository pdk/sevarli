#!/bin/bash

# exit on any error
set -e

if [ -z "$1" ]
then
    echo usage: script.sh pattern
    exit 1
fi

setvars=`sevarli -data example.data -pattern $1`
eval $setvars

echo $CITY