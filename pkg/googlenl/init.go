package googlenl

import (
	"context"

	language "cloud.google.com/go/language/apiv1"
	"google.golang.org/api/kgsearch/v1"
	"google.golang.org/api/option"
)

var NLUClient *language.Client
var KnowledgeGraphClient *kgsearch.Service

func Init(apikey string) error {
	var err error

	NLUClient, err = language.NewRESTClient(context.Background(), option.WithAPIKey(apikey))
	if err != nil {
		return err
	}

	KnowledgeGraphClient, err = kgsearch.NewService(context.Background(), option.WithAPIKey(apikey))
	if err != nil {
		return err
	}

	return nil
}
