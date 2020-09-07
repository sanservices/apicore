package apisettings

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/san-services/apilogger"
	"gopkg.in/yaml.v2"
)

const (
	// EmailKey key to lookup email value in request context
	EmailKey string = "email"

	// UsernNameKey key to lookup username value in request context
	UsernNameKey string = "preferred_username"

	// defines date formats for layouts
	SimpleDateLayout = "2006-01-02"

	// EnrolledProgramsMaxQuantity maximun amount of programs
	// an user can be enrolled at once.
	EnrolledProgramsMaxQuantity = 5
)

type Settings struct {
	Service Service  `yaml:"Service"`
	Cache   Cache    `yaml:"Cache"`
	DB      Database `yaml:"Database"`
}

type Service struct {
	Name       string `yaml:"name"`
	PathPrefix string `yaml:"path_prefix"`
	Version    string `yaml:"version"`
	Port       int    `yaml:"port"`
	Debug      bool   `yaml:"debug"`
}

type Cache struct {
	Enabled           bool          `yaml:"enabled"`
	ExpirationMinutes time.Duration `yaml:"expiration_minutes"`
	PurgeMinutes      time.Duration `yaml:"purge_minutes"`
}

type Database struct {
	Engine   string `yaml:"engine"`
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// GetEnv returns env value, if empty, returns fallback value
func GetEnv(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}

//Get returns main configuration object
func Get(ctx context.Context, filePath string) (*Settings, error) {
	lg := apilogger.New(ctx, "")

	var err error
	var confFile []byte

	config := &Settings{}

	confFile, err = ioutil.ReadFile(filePath)
	if err != nil {
		lg.Error(apilogger.LogCatReadConfig, err)
		return nil, err
	}

	//if file exists use its variables
	if err == nil {
		err = yaml.Unmarshal(confFile, &config)
		if err != nil {
			lg.Error(apilogger.LogCatUnmarshalReq, err)
			return nil, err
		}
	}

	return config, nil
}
