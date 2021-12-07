package conf

import "flag"

type FlagsConfiguration struct {
	IsProduction      bool
	IsAuthorizationOn bool
	IsLoggingOn       bool
	Help              bool
}

func ParseFlags() (*FlagsConfiguration, error) {
	conf := FlagsConfiguration{}

	flag.BoolVar(&conf.IsProduction, "prod", false, "Set to define that this build is official production server")
	flag.BoolVar(&conf.IsAuthorizationOn, "auth", false, "Set to make sure proper endpoints are access only by authorized users")
	flag.BoolVar(&conf.IsLoggingOn, "log", false, "Set to save all logs (logs in console are always printed)")
	flag.BoolVar(&conf.Help, "help", false, "Help flag - for printing this message")
	flag.Parse()

	if conf.Help {
		flag.PrintDefaults()
		panic(0)
	}

	return &conf, nil
}
