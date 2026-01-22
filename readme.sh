#!/bin/sh

# generates README.md

gotemplate exec -o README.md -t readme.tpl doc/*.tpl
