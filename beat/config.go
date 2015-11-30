package beat

type FpmConfig struct {
	Period *int64
	URLs   []string
}

type ConfigSettings struct {
	Input FpmConfig
}
