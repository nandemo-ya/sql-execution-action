package config

type Query struct {
	SQL    string `yaml:"sql"`
	Params []any  `yaml:"params"`
}
