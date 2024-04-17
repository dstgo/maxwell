package conf

import (
	"fmt"
	"github.com/dstgo/maxwell/pkg/cfgx"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestLoadConf(t *testing.T) {
	var conf AppConf
	err := cfgx.LoadConfigAndMapTo("./conf.yaml", &conf)
	assert.Nil(t, err)
	marshal, err := yaml.Marshal(conf)
	assert.Nil(t, err)
	fmt.Println(string(marshal))
}
