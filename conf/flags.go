package conf

import "flag"

type FlagsConfiguration struct {
	IsProduction bool
}

func ParseFlags() (*FlagsConfiguration, error) {
	conf := FlagsConfiguration{}

	flag.BoolVar(&conf.IsProduction, "prod", false, "Set to define that this build is official production server")

	flag.Parse()

	return &conf, nil
}
