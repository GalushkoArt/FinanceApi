package service

import (
	"context"
	"github.com/galushkoart/finance-api/internal/model"
	"github.com/galushkoart/finance-api/internal/repository"
	"github.com/galushkoart/finance-api/pkg/conpool"
	"github.com/galushkoart/finance-api/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type SymbolService interface {
	Add(ctx context.Context, symbol model.Symbol) error
	GetBySymbol(ctx context.Context, name string) (model.Symbol, error)
	GetAll(ctx context.Context) ([]model.Symbol, error)
	Update(ctx context.Context, symbol model.UpdateSymbol) error
	Delete(ctx context.Context, symbolName string) error
}

type symbolServiceWithRepoAndClient struct {
	repo         repository.SymbolRepository
	pool         *conpool.ConnectionPool
	auditService AuditService
}

func ssLog(c context.Context, e *zerolog.Event) *zerolog.Event {
	return utils.LogRequest(c, e).Str("from", "symbolServiceWithRepoAndClient")
}

func NewSymbolService(repo repository.SymbolRepository, pool *conpool.ConnectionPool, auditService AuditService) SymbolService {
	return &symbolServiceWithRepoAndClient{repo: repo, pool: pool, auditService: auditService}
}

func (s *symbolServiceWithRepoAndClient) Add(ctx context.Context, symbol model.Symbol) error {
	go s.auditService.LogSymbolCreated(ctx, symbol.Symbol)
	return s.repo.Add(ctx, symbol)
}

func (s *symbolServiceWithRepoAndClient) Update(ctx context.Context, symbol model.UpdateSymbol) error {
	go s.auditService.LogSymbolUpdated(ctx, symbol.Symbol)
	return s.repo.Update(ctx, symbol)
}

func (s *symbolServiceWithRepoAndClient) Delete(ctx context.Context, symbolName string) error {
	go s.auditService.LogSymbolDeleted(ctx, symbolName)
	return s.repo.Delete(ctx, symbolName)
}

func (s *symbolServiceWithRepoAndClient) GetBySymbol(ctx context.Context, name string) (model.Symbol, error) {
	symbol, err := s.repo.GetBySymbol(ctx, name)
	if err != nil {
		ssLog(ctx, log.Warn()).Err(err).Msg("Couldn't fetch data from repo! Trying to get data from api")
		timeSeries, err := s.pool.GetHistoricDataForSymbol(ctx, name)
		if err != nil {
			ssLog(ctx, log.Error()).Err(err).Interface("response", timeSeries).Msg("Failed to get data from api!")
			return model.Symbol{}, err
		}
		symbol = timeSeriesToModel(*timeSeries)
		err = s.repo.Add(ctx, symbol)
		if err != nil {
			ssLog(ctx, log.Error()).Err(err).Msgf("Couldn't save symbol!")
			return model.Symbol{}, err
		}
	}
	return symbol, nil
}

func (s *symbolServiceWithRepoAndClient) GetAll(ctx context.Context) ([]model.Symbol, error) {
	return s.repo.GetAll(ctx)
}
