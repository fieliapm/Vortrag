#!/bin/sh

locust -f dummy.py --master --master-bind-host=`/sbin/ip route|awk '/docker0/ { print $9 }'` --master-bind-port=5557
