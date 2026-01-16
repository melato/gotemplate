#!/bin/sh

# use gotemplate to concatenate the core usage and the file functions usage
gotemplate -t usage.tpl -o usage.yaml -F core=../usage.yaml -F file_funcs=usage_files.yaml
