#!/bin/bash
rm -rf ~/.glide
mkdir -p ~/.glide
glide mirror set https://golang.org/x/mobile https://github.com/golang/mobile --vcs git
glide mirror set https://golang.org/x/crypto https://github.com/golang/crypto --vcs git
glide mirror set https://golang.org/x/net https://github.com/golang/net --vcs git
glide mirror set https://golang.org/x/tools https://github.com/golang/tools --vcs git
glide mirror set https://golang.org/x/text https://github.com/golang/text --vcs git
glide mirror set https://golang.org/x/image https://github.com/golang/image --vcs git
glide mirror set https://golang.org/x/sys https://github.com/golang/sys --vcs git
glide up