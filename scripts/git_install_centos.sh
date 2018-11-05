#!/bin/bash
yum -y groupinstall "Development Tools"

yum -y install zlib-devel perl-ExtUtils-MakeMaker asciidoc xmlto openssl-devel

wget https://www.kernel.org/pub/software/scm/git/git-2.13.3.tar.gz

tar -zxvf git-2.13.3.tar.gz

cd git-2.13.3

./configure --prefix=/usr/local/git

make && make install

export PATH="/usr/local/git/bin:$PATH"

source /etc/profile

git --version