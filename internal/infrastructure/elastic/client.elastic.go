package elastic

import (
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/mi6ule/goboy/config"
)

type Elastic struct {
	Client *elasticsearch.Client
}

func NewElasticClient(conf config.ElasticConfig) (*Elastic, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			conf.Url,
		},
		Username:          conf.Username,
		Password:          conf.Pwd,
		EnableDebugLogger: true,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	elastic := Elastic{Client: client}
	elastic.InitIndecies()
	return &elastic, nil
}
