package beat

type PhpfpmConfig struct {
	Period *int64
	URLs   []string
}

type ConfigSettings struct {
	Input PhpfpmConfig
}
