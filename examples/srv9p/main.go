package main

import (
	"flag"
	"log"
	"net"
	"os"

	"aqwari.net/net/styx"
	"github.com/lemondevxyz/afero9p"
	"github.com/spf13/afero"
)

func main() {
	mode := flag.String("mode", "tcp", "which mode to listen to")
	addr := flag.String("addr", "127.0.0.1:12345", "the address to listen to")
	debug := flag.Bool("debug", false, "trace 9p messages")
	verbose := flag.Bool("verbose", false, "print extra info")

	flag.Parse()

	ln, err := net.Listen(*mode, *addr)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	var errorlog, tracelog styx.Logger
	if *debug {
		tracelog = log.New(os.Stderr, "", 0)
	}
	if *verbose {
		errorlog = log.New(os.Stderr, "", 0)
	}

	afero9p.NewServer(afero9p.ServerOptions{
		Listener: ln,
		ErrorLog: errorlog,
		TraceLog: tracelog,
	}, afero.NewOsFs())
}
