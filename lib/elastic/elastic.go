package elastic

import (
	"crypto/tls"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func NewElasticSearchClient(log *zap.SugaredLogger) *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Username:  os.Getenv("ELASTIC_USER"),
		Password:  os.Getenv("ELASTIC_PASSWORD"),
		Addresses: []string{os.Getenv("ELASTIC_HOST")},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("error while connect to elasticsearch: %v", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %v", err)
	}
	defer res.Body.Close()

	return es
}
