package main

import (
	"flag"
	"github.com/apcera/gssapi"
	"log"
	"os"
)

type Context struct {
	Mode           string
	ServiceName    string
	ServiceAddress string
	RequestPath    string
	Krb5Keytab     string
	Krb5Config     string
	LibPathMIT     bool
	LibPathHeimdal bool
	LibPath        string

	*gssapi.Lib
}

var c = &Context{}

func init() {
	flag.StringVar(&c.Mode, "mode", "", `"client" or "service"`)
	flag.StringVar(&c.ServiceName, "service-name", "SampleService", "[C,S] service name")
	flag.StringVar(&c.ServiceAddress, "service-address", ":8080", "[C,S] service address hostname:port")
	flag.StringVar(&c.RequestPath, "request", "/", "[C,S] test path to use")
	flag.StringVar(&c.Krb5Keytab, "krb5-keytab", "", "[S] path to the keytab file")
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
	if c.Krb5Keytab != "" {
		err := os.Setenv("KRB5_KEYTAB", c.Krb5Keytab)
		if err != nil {
			logger.Fatal(err)
		}
	}

	path, err := gssapi.LibPath(c.LibPath, c.LibPathMIT, c.LibPathHeimdal)
	if err != nil {
		logger.Fatal(err)
	}
	lib, err := gssapi.LoadLib(path)
	if err != nil {
		logger.Fatal(err)
	}
	c.Lib = lib
	c.Lib.Printer = logger

	switch c.Mode {
	case "client":
		err = Client(c)
	case "service":
		err = Server(c)
	default:
		flag.Usage()
		os.Exit(1)
	}

	if err != nil {
		logger.Fatalf("mode:%v error:%v\n", c.Mode, err)
	}
}
