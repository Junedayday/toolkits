#!/bin/bash
protoc --go_out=../../proto/mysql/pbmysql/ ../../proto/mysql/event.proto
protoc_mac --go_out=../../proto/mysql/pbmysql/ --proto_path=../../proto/mysql event.proto