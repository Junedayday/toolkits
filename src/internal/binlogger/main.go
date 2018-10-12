package main

import (
	"flag"
	"pkg/gomysql"
	"proto/mysql/pbmysql"

	"github.com/golang/glog"
)

var (
	tUser     = "root"
	tPassword = "abc#123"
	tIP       = "192.168.33.202"
	tPort     = 3306
)

func main() {
	// glag parse is needed in glog package
	flag.Parse()
	// glog must be flushed before exit
	defer glog.Flush()
	glog.Info("service start!")

	// add the config for a mysql instance
	cfg := gomysql.NewConnCfger(tUser, tPassword, tIP, tPort)
	sc, err := cfg.NewSlaveConner()
	if err != nil {
		glog.Errorf("new slave connection error %v\n", err)
		return
	}
	defer sc.Close()

	// declare deal message function
	f := func(msg *pbmysql.Event) {
		glog.V(1).Info("read an event : %#v\n", msg)
	}
	go sc.StartBinlogDumpFromCurrentAsProto(f)
	for {

	}
}
