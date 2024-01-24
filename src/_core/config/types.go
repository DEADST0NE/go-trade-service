package config

type BrokerTransportSS struct {
	Scheme string `json:"scheme"`
	Path   string `json:"path"`
	Host   string `json:"host"`
}

type Broker struct {
	Ss BrokerTransportSS `json:"ss"`
}

type Rabbit struct {
	Url string `json:"url"`
}

type Config struct {
	Broker  Broker   `json:"broker"`
	Symbols []string `json:"symbols"`
	Rabbit  Rabbit   `json:"rabbit"`
}
