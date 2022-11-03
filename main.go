package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dgraph-io/badger"
	"github.com/tendermint/tendermint/libs/log"
	abciserver "github.com/tendermint/tendermint/abci/server"
)


var socketAddr string

func init(){
	flag.StringVar(&socketAddr, "socket-addr", "example address", "Unix domain socket address")
}


func main() {
	db, err := badger.Open(badger.DefaultOptions("/tmp/bager"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open badger db: $v", err)
		os.Exit(1)
	}
	defer db.Close()
	app:=NewKVStoreApplication(db)
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	server := abciserver.NewSocketServer(socketAddr, app)
	server.SetLogger(logger)
	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "err starting socket server: $v", err)
		os.Exit(1)
	}
	defer server.Stop()

	c:=make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(0)
}
