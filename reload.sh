#!/bin/bash
echo "Username: $1";

/usr/bin/flock -w 0 /var/run/192.168.11.100.lock /home/mac/goServer/reconect.sh -r 4G  -i 192.168.11.1 /etc/init.d/3proxy start192.168.11.1 >/dev/null 2>&1
