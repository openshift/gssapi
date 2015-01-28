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
	flag.StringVar(&c.ServiceName, "service-name", "SampleService", "service name")
	flag.StringVar(&c.ServiceAddress, "service-address", ":8080", "service address hostname:port")
	flag.StringVar(&c.Krb5Ktname, "krb5-ktname", "", "path to the keytab file")
	flag.StringVar(&c.Krb5Config, "krb5-config", "", "path to krb5.config file")
	flag.StringVar(&c.LibPath, "gssapi-path", "", "use the specified path to libgssapi.so")
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

	o := &gssapi.Options{
		LibPath: c.LibPath,
		Printer: logger,
	}
	var err error
	c.Lib, err = gssapi.LoadLib(o)
	if err != nil {
		logger.Fatal(err)
	}

	err = Server(c)
	if err != nil {
		logger.Fatal(err)
	}
}
