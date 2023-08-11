package elastic

import (
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
)

func NewElasticClient(conf config.ElasticConfig) (*elasticsearch.Client, error) {
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
	return client, nil
}

func NewTypedElasticClient(conf config.ElasticConfig) (*elasticsearch.TypedClient, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			conf.Url,
		},
		Username:          conf.Username,
		Password:          conf.Pwd,
		EnableDebugLogger: true,
	}

	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, err
	}
	return client, nil
}
