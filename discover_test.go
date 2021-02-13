package freebox

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestDiscoverHTTP(t *testing.T) {
	devices, err := Discover(DiscoverProtocolHTTP)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	fmt.Printf("%+v\n", devices)
}
