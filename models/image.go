package models

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Image struct {
	PublicURL string `gorm:"size:2048" json:"public_url"`
}

func SaveAsS3(base64Content string, extension string, awsAccessKeyId string, awsSecretAccessKey string, region string) (string, error) {
	base64Decode, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return "", err
	}

	randonFileName := time.Now().Format("20060102150405") + "_" + fmt.Sprint(rand.Intn(1000000)) + "." + extension

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyId, awsSecretAccessKey, ""),
	},
	)

	if err != nil {
		return "", err
	}

	reader := strings.NewReader(string(base64Decode))

	uploader := s3manager.NewUploader(sess)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String("junges-images"),
		Key:         aws.String(randonFileName),
		Body:        reader,
		ContentType: aws.String("image/" + extension),
	})
	if err != nil {
		return "", err
	}

	return result.Location, nil
}
