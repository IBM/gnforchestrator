package config

type PingConfig struct {
	PongAddress string `json:"pongAddress"`
	PongPort    string `json:"pongPort"`
}