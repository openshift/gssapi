package client

import (
	"flag"
	"github.com/apcera/gssapi"
	"log"
	"os"
	"testing"
)

type Context struct {
	ServiceName    string
	ServiceAddress string
	RequestPath    string
	Krb5Config     string
	LibPathMIT     bool
	LibPathHeimdal bool
	LibPath        string
}

var c = &Context{}

func init() {
	flag.StringVar(&c.ServiceName, "service-name", "SampleService", "[C,S] service name")
	flag.StringVar(&c.ServiceAddress, "service-address", ":8080", "[C,S] service address hostname:port")
	flag.StringVar(&c.RequestPath, "request", "/", "[C,S] test path to use")
	flag.StringVar(&c.Krb5Config, "krb5-config", "", "[C,S] path to krb5.config file")
	flag.BoolVar(&c.LibPathMIT, "gssapi-mit", false, "[C,S] use the default MIT library path (libgssapi_krb5.so)")
	flag.BoolVar(&c.LibPathHeimdal, "gssapi-heimdal", false, "[C,S] use the default Heimdal library path (libgssapi.so)")
	flag.StringVar(&c.LibPath, "gssapi-path", "", "[C,S] use the specified path to libgssapi.so")
}

func loadlib(tb testing.TB, verbose bool) *gssapi.Lib {
	lib, err := gssapi.LoadLib(gssapi.LibPath(
		c.LibPath, c.LibPathMIT, c.LibPathHeimdal))
	if err != nil {
		tb.Fatal(err)
	}
	if verbose {
		lib.Printer = log.New(os.Stderr, "gssapi-sample: ", log.LstdFlags)
	}
	return lib
}

func TestMain(m *testing.M) {
	//flag.Parse()

	if c.Krb5Config != "" {
		err := os.Setenv("KRB5_CONFIG", c.Krb5Config)
		if err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(m.Run())
}
