package service

import (
	"context"
	audit "github.com/GalushkoArt/GoAuditService/pkg/proto"
	"github.com/galushkoart/finance-api/pkg/service"
	"github.com/galushkoart/finance-api/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type auditServiceWithClientAndPublisher struct {
	client           *service.AuditClient
	publisher        service.AuditPublisher
	clientEnabled    bool
	publisherEnabled bool
}

type AuditService interface {
	LogSymbolCreated(ctx context.Context, symbol string)
	LogSymbolUpdated(ctx context.Context, symbol string)
	LogSymbolDeleted(ctx context.Context, symbol string)
	LogUserSignUp(ctx context.Context, userID string)
	LogUserSignIn(ctx context.Context, userID string)
	LogUserRefreshToken(ctx context.Context, userID string)
}

var auditLog zerolog.Logger

func NewAuditService(clientEnabled bool, client *service.AuditClient, publisherEnabled bool, publisher service.AuditPublisher) AuditService {
	auditLog = log.With().Str("from", "auditService").Logger()
	return &auditServiceWithClientAndPublisher{client: client, publisher: publisher, clientEnabled: clientEnabled, publisherEnabled: publisherEnabled}
}

func (s *auditServiceWithClientAndPublisher) sendRequest(ctx context.Context, request *audit.LogRequest) {
	var err error
	if s.publisherEnabled {
		err = s.publisher.Publish(ctx, request)
		if err != nil {
			utils.LogRequest(ctx, auditLog.Error()).Err(err).Interface("request", request).Msg("Fail in audit publisher!")
		}
	}
	if s.clientEnabled && err != nil {
		response, err := s.client.SendRequest(ctx, request)
		if err != nil || response.Answer == audit.Response_ERROR {
			utils.LogRequest(ctx, auditLog.Error()).Err(err).Interface("request", request).Msg("Fail in audit server!")
		}
	}
}

func (s *auditServiceWithClientAndPublisher) LogSymbolCreated(ctx context.Context, symbol string) {
	s.sendRequest(ctx, &audit.LogRequest{
		Action:    audit.LogRequest_CREATE,
		Entity:    audit.LogRequest_SYMBOL,
		EntityId:  symbol,
		Timestamp: timestamppb.Now(),
		RequestId: utils.GetRequestId(ctx),
	})
}

func (s *auditServiceWithClientAndPublisher) LogSymbolUpdated(ctx context.Context, symbol string) {
	s.sendRequest(ctx, &audit.LogRequest{
		Action:    audit.LogRequest_UPDATE,
		Entity:    audit.LogRequest_SYMBOL,
		EntityId:  symbol,
		Timestamp: timestamppb.Now(),
		RequestId: utils.GetRequestId(ctx),
	})
}

func (s *auditServiceWithClientAndPublisher) LogSymbolDeleted(ctx context.Context, symbol string) {
	s.sendRequest(ctx, &audit.LogRequest{
		Action:    audit.LogRequest_DELETE,
		Entity:    audit.LogRequest_SYMBOL,
		EntityId:  symbol,
		Timestamp: timestamppb.Now(),
		RequestId: utils.GetRequestId(ctx),
	})
}

func (s *auditServiceWithClientAndPublisher) LogUserSignUp(ctx context.Context, userID string) {
	s.sendRequest(ctx, &audit.LogRequest{
		Action:    audit.LogRequest_SIGN_UP,
		Entity:    audit.LogRequest_USER,
		EntityId:  userID,
		Timestamp: timestamppb.Now(),
		RequestId: utils.GetRequestId(ctx),
	})
}

func (s *auditServiceWithClientAndPublisher) LogUserSignIn(ctx context.Context, userID string) {
	s.sendRequest(ctx, &audit.LogRequest{
		Action:    audit.LogRequest_SIGN_IN,
		Entity:    audit.LogRequest_USER,
		EntityId:  userID,
		Timestamp: timestamppb.Now(),
		RequestId: utils.GetRequestId(ctx),
	})
}

func (s *auditServiceWithClientAndPublisher) LogUserRefreshToken(ctx context.Context, userID string) {
	s.sendRequest(ctx, &audit.LogRequest{
		Action:    audit.LogRequest_REFRESH,
		Entity:    audit.LogRequest_USER,
		EntityId:  userID,
		Timestamp: timestamppb.Now(),
		RequestId: utils.GetRequestId(ctx),
	})
}
