package conf

import "flag"

type Configuration struct {
	IsProduction bool
}

func Init() (*Configuration, error) {

	conf := Configuration{}

	flag.BoolVar(&conf.IsProduction, "prod", false, "Set to define that this build is official production server")

	flag.Parse()

	return &conf, nil
}
