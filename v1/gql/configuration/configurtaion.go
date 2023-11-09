package configuration

import "os"

const defaultPort = "8080"
const defaultServerName = "no name"
const defaultNatsURL = "nats://localhost:4222"

type Config struct {
	ServerName    string
	ServerPort    string
	NatsURL       string
	VersionRef    string
	VersionCommit string
}

func FromEnvironment() (Config, error) {
	port := defaultPort
	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}

	natsURL := defaultNatsURL
	if url, ok := os.LookupEnv("NATS_URL"); ok {
		natsURL = url
	}

	serverName := defaultServerName
	if name, ok := os.LookupEnv("SERVER_NAME"); ok {
		serverName = name
	}

	return Config{
		ServerName: serverName,
		ServerPort: port,
		NatsURL:    natsURL,
	}, nil
}
