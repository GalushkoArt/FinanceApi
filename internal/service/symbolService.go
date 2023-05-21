package service

import (
	"FinanceApi/internal/model"
	"FinanceApi/internal/repository"
	"FinanceApi/pkg/connectionPool"
	"FinanceApi/pkg/utils"
	"context"
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
	repo repository.SymbolRepository
	pool *connectionPool.ConnectionPool
}

func ssLog(c context.Context, e *zerolog.Event) *zerolog.Event {
	return utils.LogRequest(c, e).Str("from", "symbolServiceWithRepoAndClient")
}

func NewSymbolService(repo repository.SymbolRepository, pool *connectionPool.ConnectionPool) SymbolService {
	return &symbolServiceWithRepoAndClient{repo: repo, pool: pool}
}

func (s *symbolServiceWithRepoAndClient) Add(ctx context.Context, symbol model.Symbol) error {
	return s.repo.Add(ctx, symbol)
}

func (s *symbolServiceWithRepoAndClient) Update(ctx context.Context, symbol model.UpdateSymbol) error {
	return s.repo.Update(ctx, symbol)
}

func (s *symbolServiceWithRepoAndClient) Delete(ctx context.Context, symbolName string) error {
	return s.repo.Delete(ctx, symbolName)
}

func (s *symbolServiceWithRepoAndClient) GetBySymbol(ctx context.Context, name string) (model.Symbol, error) {
	symbol, err := s.repo.GetBySymbol(ctx, name)
	if err != nil {
		ssLog(ctx, log.Warn()).Err(err).Msg("Couldn't fetch data from repo! Trying to get data from api")
		timeSeries, err := s.pool.GetHistoricDataForSymbol(ctx, name)
		if err != nil {
			return model.Symbol{}, err
		}
		symbol = timeSeriesToModel(*timeSeries)
		err = s.repo.Add(ctx, symbol)
		if err != nil {
			ssLog(ctx, log.Error()).Err(err).Msgf("Couldn't save symbol!")
		}
	}
	return symbol, nil
}

func (s *symbolServiceWithRepoAndClient) GetAll(ctx context.Context) ([]model.Symbol, error) {
	return s.repo.GetAll(ctx)
}
