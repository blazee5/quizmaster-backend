package elastic

import (
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"net/http"
	"os"
)

type BulkIndexerConfig struct {
	NumWorkers           int `mapstructure:"numWorkers" validate:"required"`
	FlushBytes           int `mapstructure:"flushBytes" validate:"required"`
	FlushIntervalSeconds int `mapstructure:"flushIntervalSeconds" validate:"required"`
	TimeoutMilliseconds  int `mapstructure:"timeoutMilliseconds" validate:"required"`
}

type Config struct {
	Addresses []string
	Username  string
	Password  string

	APIKey        string
	Header        http.Header
	EnableLogging bool
}

func NewElasticSearchClient() *elasticsearch.Client {

	config := elasticsearch.Config{
		Addresses: []string{},
		Username:  os.Getenv("ELASTIC_USER"),
		APIKey:    os.Getenv("ELASTIC_API_KEY"),
		Header:    http.Header{},
		Logger:    &elastictransport.ColorLogger{Output: os.Stdout, EnableRequestBody: true, EnableResponseBody: true},
	}

	client, err := elasticsearch.NewClient(config)

	if err != nil {
		panic(err)
	}

	return client
}
