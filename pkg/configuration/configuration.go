package configuration

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Configuration struct used by GetConfiguration
type Configuration struct {
	Port               string `default:"8080"`
	OauthAudience      string `required:"true" split_words:"true"`
	OauthIssuer        string `required:"true" split_words:"true"`
	GoogleCloudProject string `required:"true" split_words:"true"`
}

var configuration = &Configuration{}

func init() {
	if err := envconfig.Process("", configuration); err != nil {
		log.Fatalf("Failed to process config: %v", err)
	}
}

// GetConfiguration Get the app configuration
func GetConfiguration() *Configuration {
	return configuration
}
