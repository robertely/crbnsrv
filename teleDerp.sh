#!/bin/bash

# Default 1
MAXTHREADS=${1:-1}

function kall {
  echo "-> Killing all telegraf instances."
  pkill -9 telegraf
}
trap kall INT TERM EXIT

COUNT=1
while [ $COUNT -le $MAXTHREADS ]; do
  export FAKEHOST=$(hostname)$(printf %04d $COUNT)
  echo "-> Starting Thread with host: $FAKEHOST"
  telegraf --config teleDerp.conf --input-filter mem --output-filter graphite &
  ((COUNT++))
  sleep .01 # You can brush up against pid exhaustion if you do this too fast
done
sleep .5 
echo "-----------------------------------------"
echo "-> Started $MAXTHREADS threads"
echo "-----------------------------------------"
wait
