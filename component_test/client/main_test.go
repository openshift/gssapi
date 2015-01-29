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
	flag.StringVar(&c.ServiceName, "service-name", "SampleService", "service name")
	flag.StringVar(&c.ServiceAddress, "service-address", ":8080", "service address hostname:port")
	flag.StringVar(&c.RequestPath, "request", "/", "test path to use")
	flag.StringVar(&c.Krb5Config, "krb5-config", "", "path to krb5.config file")
	flag.StringVar(&c.LibPath, "gssapi-path", "", "use the specified path to libgssapi shared object")
}

func loadlib(tb testing.TB, verbose bool) *gssapi.Lib {
	o := &gssapi.Options{
		LibPath: c.LibPath,
	}
	if verbose {
		o.Printer = log.New(os.Stderr, "gssapi-sample: ", log.LstdFlags)
	}
	lib, err := gssapi.LoadLib(o)
	if err != nil {
		tb.Fatal(err)
	}
	return lib
}

func TestMain(m *testing.M) {
	if testing.Verbose() {
		log.Printf("gssapi-sample: Config:\n%#v", c)
	}

	if c.Krb5Config != "" {
		err := os.Setenv("KRB5_CONFIG", c.Krb5Config)
		if err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(m.Run())
}
