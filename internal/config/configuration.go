package config

import "flag"

type Configuration struct {
	Addr *string
	Dsn  *string
}

func NewConfiguration() *Configuration {
	configObject := new(Configuration)
	configObject.Addr = flag.String("addr", ":8888", "customHttp listen port")
	configObject.Dsn = flag.String("dsn", "neo4j://localhost:7687", "Sqllite3 data source name")
	flag.Parse()
	return configObject
}
