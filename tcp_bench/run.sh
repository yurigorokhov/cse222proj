#!/bin/sh
/usr/sbin/apache2ctl -D FOREGROUND &
./bin/tcp_bench
