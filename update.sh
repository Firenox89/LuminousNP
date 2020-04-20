#! /bin/bash
~/apps/nodemcu-firmware-3.0-master_20190907/luac.cross -f -o lfs.img lua/lfs/* && curl 192.168.2.111/OTA
