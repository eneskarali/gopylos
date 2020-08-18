package github.com/eneskarali/gopylos

import (
	"context"
	"crypto/tls"
	"errors"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"github.com/eneskarali/gopylos/gopylos/collector"
	"log"
)

var (
	Key = ""
)

var (
	client collector.EventCollectorClient
)

func Start() error {
	if len(Key) <= 0 {
		return errors.New("access token should be initialized")
	}

	perRPC := oauth.NewOauthAccess(fetchToken(Key))
	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := grpc.Dial("localhost:8086", grpc.WithTransportCredentials(credentials.NewTLS(config)), grpc.WithPerRPCCredentials(perRPC))
	if err != nil {
		log.Fatalf("connection error: %v", err)
		return errors.New("connection error")
	}

	client = collector.NewEventCollectorClient(conn)
	return nil
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
