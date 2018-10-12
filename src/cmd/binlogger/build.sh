#!/bin/bash
go build internal/binlogger
cp ../../configs/binlogger/binlogger_sample.yaml ./
# tar -cvf binlogger.tar binlogger binlogger_sample.yaml
