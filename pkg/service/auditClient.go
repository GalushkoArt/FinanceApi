package service

import (
	"context"
	audit "github.com/GalushkoArt/GoAuditService/pkg/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuditClient struct {
	client audit.AuditServiceClient
	conn   *grpc.ClientConn
}

func NewAuditClient(enabled bool, address string) (*AuditClient, error) {
	if !enabled {
		log.Info().Msg("Audit client disabled!")
		return nil, nil
	}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := audit.NewAuditServiceClient(conn)
	log.Info().Msg("Audit client connected!")
	return &AuditClient{client: client, conn: conn}, nil
}

func (a *AuditClient) Close() error {
	if a == nil {
		return nil
	}
	return a.conn.Close()
}

func (a *AuditClient) SendRequest(ctx context.Context, request *audit.LogRequest) (*audit.Response, error) {
	return a.client.Log(ctx, request)
}
