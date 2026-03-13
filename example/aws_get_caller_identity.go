package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/mister-webhooks/sts-assume-role-proxy/client"
)

const SOCKET_PATH = "/tmp/sts-assume-role-proxy.sock"

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithCredentialsProvider(
		client.NewProxyCredentialsProvider(SOCKET_PATH),
	))

	if err != nil {
		log.Fatal(err)
	}

	stsService := sts.NewFromConfig(cfg)

	result, err := stsService.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ARN of credentials provided by sts-assume-role-proxy: %s\n", *result.Arn)
}
