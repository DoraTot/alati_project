package configuration

import "os"

type Configuration struct {
	Address       string
	JaegerAddress string
}

func GetConfiguration() Configuration {
	return Configuration{
		Address:       os.Getenv("SERVICE_ADDRESS"),
		JaegerAddress: os.Getenv("JAEGER_ADDRESS"),
	}
}
