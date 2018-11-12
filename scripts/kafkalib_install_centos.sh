#!/bin/bash
yum install -y gcc gcc-c++ make 
git clone https://github.com/edenhill/librdkafka.git
cd librdkafka
./configure --prefix /usr
make
make install
export PKG_CONFIG_PATH=/usr/lib/pkgconfig
