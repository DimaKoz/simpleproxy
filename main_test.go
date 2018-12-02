package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	go main()
	ret := m.Run()
	os.Exit(ret)
}
