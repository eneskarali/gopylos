package gopylos

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/eneskarali/gopylos/collector"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"log"
)

var (
	Key = ""
)

type pylos struct {
	SendUsage func(string, string, float32) error
}

var (
	client collector.AppUsageCollectorClient
)

func Init(host string) (pylos, error) {
	if len(Key) <= 0 {
		return pylos{}, errors.New("access token should be initialized")
	}

	perRPC := oauth.NewOauthAccess(fetchToken(Key))
	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(credentials.NewTLS(config)), grpc.WithPerRPCCredentials(perRPC))
	if err != nil {
		log.Fatalf("connection error: %v", err)
		return pylos{}, errors.New("connection error")
	}
	client = collector.NewAppUsageCollectorClient(conn)

	pylos := pylos{SendUsage: SendUsage}

	return pylos, nil
}

func SendUsage(planId string, userId string, usage float32) error {
	_, err := client.SendUsage(context.Background(), &collector.UsageRequest{PlanId: planId, UserId: userId, Usage: usage})
	if err != nil {
		log.Fatalf("send event error, %v", err)
		return err
	}
	return nil
}

func fetchToken(key string) *oauth2.Token {
	return &oauth2.Token{
		AccessToken: key,
	}
}
