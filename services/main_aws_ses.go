package services

import (
	"smart-scene-app-api/config"
	"smart-scene-app-api/pkg"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

func NewAwsSes() (*pkg.AWSSesClient, error) {
	inst, err := pkg.NewAwsSesClient(
		&aws.Config{
			Region:      aws.String(config.Config.AwsSes.Region),
			Credentials: credentials.NewStaticCredentials(config.Config.AwsSes.AccessKeyID, config.Config.AwsSes.SecretAccessKey, ""),
		})
	if err != nil {
		return nil, err
	}
	return inst, nil
}
