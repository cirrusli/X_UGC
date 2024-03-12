package es

import (
	"X_UGC/conf"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
)

var ES *elasticsearch.Client

func New(c *conf.ES) error {
	addr := fmt.Sprintf("http://%s:%d", c.Host, c.Port)

	var err error
	ES, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{addr},
		Username:  c.UserName,
		Password:  c.Password,
	})
	if err != nil {
		return err
	}
	return nil
}
