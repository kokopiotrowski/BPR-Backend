package conf

import (
	"flag"
	"sync"
)

type FlagsConfiguration struct {
	IsProduction      bool
	IsAuthorizationOn bool
	IsLoggingOn       bool
	Help              bool
}

var (
	FlagConf *FlagsConfiguration
	lock     sync.Mutex
)

func ParseFlags() (*FlagsConfiguration, error) {
	lock.Lock()
	defer lock.Unlock()

	if FlagConf != nil {
		return FlagConf, nil
	}

	FlagConf = &FlagsConfiguration{}

	flag.BoolVar(&FlagConf.IsProduction, "prod", false, "Set to define that this build is official production server")
	flag.BoolVar(&FlagConf.IsAuthorizationOn, "auth", false, "Set to make sure proper endpoints are access only by authorized users")
	flag.BoolVar(&FlagConf.IsLoggingOn, "log", false, "Set to save all logs in file (logs in console are always printed)")
	flag.BoolVar(&FlagConf.Help, "help", false, "Help flag - for printing this message")
	flag.Parse()

	if FlagConf.Help {
		flag.PrintDefaults()
		panic(0)
	}

	return FlagConf, nil
}
