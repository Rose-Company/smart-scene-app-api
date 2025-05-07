package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type AWSSesClient struct {
	Client *ses.SES
}

func NewAwsSesClient(config *aws.Config) (*AWSSesClient, error) {
	sesSession, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	svc := ses.New(sesSession)
	i := AWSSesClient{
		Client: svc,
	}

	return &i, nil
}
