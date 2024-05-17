package main

import (
	"ETicaret/Database"
	"ETicaret/Handlers"
	"ETicaret/Router"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	database.Connect()
	app := Router.Routes()
	database.ConnectRedis()
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

}

func init() {
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			secretKey,
			"",
		),
	})

	if err != nil {
		panic(err)
	}

	Handlers.Uploader = s3manager.NewUploader(awsSession)
	Handlers.Downloader = s3.New(awsSession)
}

var region string = "eu-north-1"
var accessKey string = ""
var secretKey string = ""
