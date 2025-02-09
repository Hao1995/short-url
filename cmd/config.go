package main

import (
	"log"

	"github.com/caarlos0/env/v11"
)

var cfg Config

func init() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("failed to parse config", err)
	}
	log.Print("cfg: ", cfg)
}

type Config struct {
	App   App   `envPrefix:"APP_"`
	MySQL MySQL `envPrefix:"MYSQL_"`
	Redis Redis `envPrefix:"REDIS_"`
	Cache Cache `envPrefix:"CACHE_"`
}

type App struct {
	Name string `env:"NAME,required" envDefault:"short_url"`
	Port string `env:"PORT,required" envDefault:"8080"`
	Env  string `env:"ENV,required" envDefault:"dev"`
}

type MySQL struct {
	Host     string `env:"HOST,required" envDefault:"mysql"`
	Port     string `env:"PORT,required" envDefault:"3306"`
	User     string `env:"USER,required" envDefault:"root"`
	Password string `env:"PASSWORD,required" envDefault:"root"`
	DB       string `env:"DB,required" envDefault:"short_url"`
}

type Redis struct {
	Addrs map[string]string `env:"ADDRS" envSeparator:"-" envKeyValSeparator:"|" envDefault:"server1|redis:6379"`
}

type Cache struct {
	Size      int `env:"SIZE,required" envDefault:"10000"`
	LocalTTL  int `env:"LOCAL_TTL,required" envDefault:"600"`
	SharedTTL int `env:"SHARED_TTL,required" envDefault:"3600"`
}
