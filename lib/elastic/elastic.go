package elastic

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
)

func NewElasticSearchClient(log *zap.SugaredLogger) *elasticsearch.Client {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		fmt.Println("error while connect to elasticsearch")
	}

	//res, err := es.Info()
	//if err != nil {
	//	fmt.Printf("Error getting response: %s", err)
	//}
	//defer res.Body.Close()

	return es
}
