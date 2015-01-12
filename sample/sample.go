package main

import (
	"flag"
	"fmt"
	"github.com/apcera/gssapi"
	"os"
)

type Context struct {
	Mode           string
	ServiceName    string
	ServiceAddress string
	RequestPath    string
	Keytab         string

	*gssapi.Lib
}

var c = &Context{}

func init() {
	flag.StringVar(&c.Mode, "mode", "", `"client" or "server"`)
	flag.StringVar(&c.ServiceName, "service-name", "SampleService", "[C,S] service name")
	flag.StringVar(&c.ServiceAddress, "service-address", ":8881", "[C,S] service address hostname:port")
	flag.StringVar(&c.RequestPath, "request-path", "/", "[C,S] test path to use")
	flag.StringVar(&c.Keytab, "keytab", "", "[S] path to keytab file")
}

func main() {
	flag.Parse()

	l, err := gssapi.LoadLib()
	if err != nil {
		panic(fmt.Errorf("failed to load gssapi library: %s\n", err))
	}
	c.Lib = l

	/*
		switch c.Mode {
		case "client":
			err = Client(c)
		case "server":
			err = Server(c)
		default:
			flag.Usage()
			os.Exit(1)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	*/

	u, err := Server(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error serving: %s\n", err)
		os.Exit(1)
	}

	c.ServiceAddress = u
	if err := Client(c); err != nil {
		fmt.Fprintf(os.Stderr, "Error on client: %s\n", err)
		os.Exit(1)
	}
}
