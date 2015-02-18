// Copyright 2013-2015 Apcera Inc. All rights reserved.

package test

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/apcera/gssapi"
)

type Context struct {
	DebugLog       bool
	RunAsService   bool
	ServiceName    string
	ServiceAddress string

	gssapi.Options

	*gssapi.Lib `json:"-"`
	loadonce    sync.Once
}

var c = &Context{}

func init() {
	flag.BoolVar(&c.DebugLog, "debug", false, "Output debug log")
	flag.BoolVar(&c.RunAsService, "service", false, "Stay running as sample service after executing the tests")
	flag.StringVar(&c.ServiceName, "service-name", "SampleService", "service name")
	flag.StringVar(&c.ServiceAddress, "service-address", ":8080", "service address hostname:port")
	flag.StringVar(&c.Options.LibPath, "gssapi-path", "", "use the specified path to libgssapi shared object")
	flag.StringVar(&c.Options.Krb5Ktname, "krb5-ktname", "", "path to the keytab file")
	flag.StringVar(&c.Options.Krb5Config, "krb5-config", "", "path to krb5.config file")
}

func loadlib(debug bool, prefix string) (*gssapi.Lib, error) {
	max := gssapi.Err + 1
	if debug {
		max = gssapi.MaxSeverity
	}
	pp := make([]gssapi.Printer, 0, max)
	for i := gssapi.Severity(0); i < max; i++ {
		p := log.New(os.Stderr,
			fmt.Sprintf("%s: %s\t", prefix, i),
			log.LstdFlags)
		pp = append(pp, p)
	}
	c.Options.Printers = pp

	lib, err := gssapi.Load(&c.Options)
	if err != nil {
		return nil, err
	}
	return lib, nil
}

func TestMain(m *testing.M) {
	flag.Parse()
	prefix := "go-gssapi-test-client"
	if c.RunAsService {
		prefix = "go-gssapi-test-service"
	}
	lib, err := loadlib(c.DebugLog, prefix)
	if err != nil {
		log.Fatal(err)
	}
	c.Lib = lib

	j, _ := json.MarshalIndent(c, "", "  ")
	c.Debug(fmt.Sprintf("Config: %s", string(j)))

	code := m.Run()
	if code != 0 {
		os.Exit(code)
	}

	if c.RunAsService {
		log.Fatal(Service(c))
	}
}
