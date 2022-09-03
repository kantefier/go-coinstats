package main

type Config struct {
	Waves WavesConfig
}

type WavesConfig struct {
	NodeUrl   string
	Addresses []string
}
