package config

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"gopkg.in/yaml.v2"
)

type APICfg struct {
	Key    string
	Secret string
}

type Buyin struct {
	Ticker string
	Num    float64
	Price  float64
}

type Cfg struct {
	API    APICfg
	Buyins []Buyin
}

func getConfigPath() (string, error) {
	// Read $HOME/.amirich.yaml
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	dir := usr.HomeDir
	return path.Join(dir, ".amirich.yaml"), nil
}

func generateConfig(configPath string) (Cfg, error) {
	cfg := Cfg{}
	dat, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(dat, &cfg)
	if err != nil {
		return Cfg{}, err
	}

	return cfg, nil
}

func GetConfig() (Cfg, error) {
	cfgPath, err := getConfigPath()
	if err != nil {
		return Cfg{}, fmt.Errorf("could not determine current user. Like. What. %+v", err)
	}

	cfg, err := generateConfig(cfgPath)
	if err != nil {
		return Cfg{}, fmt.Errorf("couldn't read configuration file. %+v", err)
	}

	return cfg, nil
}
