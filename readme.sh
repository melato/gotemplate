#!/bin/sh

# generates README.md

gotemplate exec -o README.md -t README.tpl doc/common/*.tpl doc/*.tpl 
