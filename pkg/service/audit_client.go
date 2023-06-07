package service

import (
	"context"
	audit "github.com/GalushkoArt/GoAuditService/pkg/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"sync"
)

type AuditClient struct {
	client audit.AuditServiceClient
	conn   io.Closer
	wg     *sync.WaitGroup
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
	return &AuditClient{client: client, conn: conn, wg: &sync.WaitGroup{}}, nil
}

func (a *AuditClient) Close() error {
	if a == nil {
		return nil
	}
	a.wg.Wait()
	return a.conn.Close()
}

func (a *AuditClient) SendRequest(ctx context.Context, request *audit.LogRequest) (*audit.Response, error) {
	a.wg.Add(1)
	response, err := a.client.Log(ctx, request)
	a.wg.Done()
	return response, err
}
