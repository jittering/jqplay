package config

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/joeshaw/envdecode"
	"github.com/owenthereal/jqplay/jq"
)

type Config struct {
	Host      string `env:"HOST,default=0.0.0.0"`
	Port      string `env:"PORT,default=3000"`
	GinMode   string `env:"GIN_MODE,default=debug"`
	AssetHost string `env:"ASSET_HOST"`
	JQVer     string

	NoOpen bool

	JSON string
}

func (c *Config) IsProd() bool {
	return c.GinMode == "release"
}

func Load() (*Config, error) {
	conf := &Config{}
	err := envdecode.Decode(conf)
	if err != nil {
		return nil, err
	}

	conf.JQVer = jq.Version

	flag.BoolVar(&conf.NoOpen, "no-open", false, "Do not open browser on startup")
	flag.Parse()

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		b, err := ioutil.ReadAll(os.Stdin)
		if err == nil {
			conf.JSON = string(b)
		}
	}

	return conf, nil
}
