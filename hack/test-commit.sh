#!/bin/bash

echo "testing" >> .testfile

git add -A . 

git commit -a -m "testing"

git push

sleep 5

github-replay http://k3s-master:30081 lilley2412 samples-go-api 10s

