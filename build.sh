#!/bin/sh

# generates README.md

gotemplate build -c doc/build.yaml -o .
