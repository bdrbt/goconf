package goconf_test

import (
	"testing"

	"github.com/bdrbt/goconf"
)

func TestLoadEnv(t *testing.T) {
	type TestConf struct {
		String  string `env:"SIMPLE_STRING"`
		Integer int    `env:"INT_VALUE"`
		Boolean bool   `env:"BOOLEAN_VALUE"`
	}

	var tc TestConf

	t.Setenv("SIMPLE_STRING", "simple")
	t.Setenv("INT_VALUE", "42")
	t.Setenv("BOOLEAN_VALUE", "true")

	if err := goconf.LoadEnv(&tc); err != nil {
		t.Error("Error while loading environment: ", err)
	}

	if tc.String != "simple" {
		t.Error("String vaue is not set")
	}

	if tc.Integer != 42 {
		t.Errorf("Integer vaue is not set, want:%d got:%d", 42, tc.Integer)
	}

	if tc.Boolean != true {
		t.Errorf("Boolean vaue is not set, want: %t got:%t", true, tc.Boolean)
	}
}
