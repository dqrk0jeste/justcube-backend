package s3_bucket

import (
	"context"

	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploaderMaker() (*manager.Uploader, error) {
	awsConfig, err := aws_config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsConfig)
	uploader := manager.NewUploader(client)
	return uploader, nil
}
