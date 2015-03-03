package app

import (
	"bytes"
	. "fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"

	. "kofalt.com/unce/def"
)

var (
	ConfigFile = []string{ ".config", "unce", "unce.toml" }
)

func GetConfigFilename() string {
	home, err := homedir.Dir()
	if err != nil { Println("Could not find homedir", err); os.Exit(1) }

	path := []string{}
	path = append(path, home)
	path = append(path, ConfigFile...)

	return filepath.Join(path...)
}

func Create() *Config {
	os.Mkdir(filepath.Dir(GetConfigFilename()), 0777)

	config := &Config {
		Producers: &Producers {
			Github: &GithubConfig {
				AccessToken: "example",
			},
		},
	}

	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	encoder.Indent = "	"
	if err := encoder.Encode(config); err != nil {
		Println("Could not encode config:", err); os.Exit(1)
	}
	err := ioutil.WriteFile(GetConfigFilename(), buf.Bytes(), 0600)

	if err != nil {
		Println("Could not write config:", err); os.Exit(1)
	} else {
		Println("New config file written to ", GetConfigFilename())
	}

	return config
}

func LoadorCreate() *Config {
	configContent, err := ioutil.ReadFile(GetConfigFilename())
	var config *Config

	if os.IsNotExist(err) {
		config = Create()
	} else if err != nil {
		Println("Could not open config:", err); os.Exit(1)
	} else {
		if _, err := toml.Decode(string(configContent), &config); err != nil {
			Println("Could not parse config:", err); os.Exit(1)
		}
	}

	return config
}
