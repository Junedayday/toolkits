#!/bin/bash
wget http://ftp.gnu.org/gnu/glibc/glibc-2.14.tar.gz 
wget http://ftp.gnu.org/gnu/glibc/glibc-ports-2.14.tar.gz 
tar -xvf  glibc-2.14.tar.gz 
tar -xvf  glibc-ports-2.14.tar.gz
mv glibc-ports-2.14 glibc-2.14/ports
mkdir glibc-2.14/build
cd glibc-2.14/build 
../configure  --prefix=/usr --disable-profile --enable-add-ons --with-headers=/usr/include --with-binutils=/usr/bin
make
make install
