#!/bin/sh
for f in *.go; do echo $f; head -3 $f | tail -1; done
