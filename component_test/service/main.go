package main

import (
	"flag"
	"github.com/apcera/gssapi"
	"log"
	"os"
)

type Context struct {
	ServiceName    string
	ServiceAddress string
	Krb5Ktname     string
	Krb5Config     string
	LibPathMIT     bool
	LibPathHeimdal bool
	LibPath        string

	*gssapi.Lib
}

var c = &Context{}

func init() {
	flag.StringVar(&c.ServiceName, "service-name", "SampleService", "[C,S] service name")
	flag.StringVar(&c.ServiceAddress, "service-address", ":8080", "[C,S] service address hostname:port")
	flag.StringVar(&c.Krb5Ktname, "krb5-ktname", "", "[S] path to the keytab file")
	flag.StringVar(&c.Krb5Config, "krb5-config", "", "[C,S] path to krb5.config file")
	flag.BoolVar(&c.LibPathMIT, "gssapi-mit", false, "[C,S] use the default MIT library path (libgssapi_krb5.so)")
	flag.BoolVar(&c.LibPathHeimdal, "gssapi-heimdal", false, "[C,S] use the default Heimdal library path (libgssapi.so)")
	flag.StringVar(&c.LibPath, "gssapi-path", "", "[C,S] use the specified path to libgssapi.so")
}

func main() {
	flag.Parse()
	logger := log.New(os.Stderr, "gssapi-sample:\t", log.LstdFlags)

	if c.Krb5Config != "" {
		err := os.Setenv("KRB5_CONFIG", c.Krb5Config)
		if err != nil {
			logger.Fatal(err)
		}
	}

	if c.Krb5Ktname != "" {
		err := os.Setenv("KRB5_KTNAME", c.Krb5Ktname)
		if err != nil {
			logger.Fatal(err)
		}
	}

	lib, err := gssapi.LoadLib(gssapi.LibPath(
		c.LibPath, c.LibPathMIT, c.LibPathHeimdal))
	if err != nil {
		logger.Fatal(err)
	}
	c.Lib = lib
	c.Lib.Printer = logger

	err = Server(c)
	if err != nil {
		logger.Fatal(err)
	}
}
