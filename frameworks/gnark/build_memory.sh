#!/bin/sh
for CMD in `ls cmd`
do
  if [ -d ./cmd/$CMD ]; then
    go build -o $CMD ./cmd/$CMD
  fi
done
