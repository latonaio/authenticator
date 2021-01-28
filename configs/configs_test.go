package configs_test

import (
	"testing"

	"bitbucket.org/latonaio/authenticator/configs"
	"github.com/davecgh/go-spew/spew"
)

func TestConfigs_Load(t *testing.T) {
	cfgs, err := configs.New()
	if err != nil {
		t.Error(err)
	}
	err = cfgs.Load()
	if err != nil {
		t.Error(err)
	}
	spew.Dump(cfgs)
}
