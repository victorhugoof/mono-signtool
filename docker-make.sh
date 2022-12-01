#!/bin/bash

MSYS_NO_PATHCONV=1 docker run --rm -ti \
 -v $PWD:/project \
 golang /bin/bash -c "cd /project && make $1"