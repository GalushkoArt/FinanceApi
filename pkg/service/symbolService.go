package service

import (
	"FinanceApi/pkg/connectionPool"
	"FinanceApi/pkg/log"
	"FinanceApi/pkg/model"
	"FinanceApi/pkg/repository"
)

type SymbolService interface {
	Add(symbol model.Symbol) error
	GetBySymbol(name string) (model.Symbol, error)
	GetAll() ([]model.Symbol, error)
	Update(symbol model.Symbol) error
	Delete(symbolName string) error
}

type symbolServiceWithRepoAndClient struct {
	repo repository.SymbolRepository
	pool *connectionPool.ConnectionPool
}

func NewSymbolService(repo repository.SymbolRepository, pool *connectionPool.ConnectionPool) SymbolService {
	return &symbolServiceWithRepoAndClient{repo: repo, pool: pool}
}

func (s *symbolServiceWithRepoAndClient) Add(symbol model.Symbol) error {
	return s.repo.Add(symbol)
}

func (s *symbolServiceWithRepoAndClient) Update(symbol model.Symbol) error {
	return s.repo.Update(symbol)
}

func (s *symbolServiceWithRepoAndClient) Delete(symbolName string) error {
	return s.repo.Delete(symbolName)
}

func (s *symbolServiceWithRepoAndClient) GetBySymbol(name string) (model.Symbol, error) {
	symbol, err := s.repo.GetBySymbol(name)
	if err != nil {
		log.Warn("Couldn't fetch data from repo! Trying to get data from api", err)
		timeSeries, err := s.pool.GetHistoricDataForSymbol(name)
		if err != nil {
			return model.Symbol{}, err
		}
		symbol = timeSeriesToModel(*timeSeries)
		err = s.repo.Add(symbol)
		if err != nil {
			log.ErrorF("Couldn't save symbol!\n%+v\n%v", symbol, err)
		}
	}
	return symbol, nil
}

func (s *symbolServiceWithRepoAndClient) GetAll() ([]model.Symbol, error) {
	return s.repo.GetAll()
}
