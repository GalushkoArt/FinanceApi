package connectionPool

import (
	"FinanceApi/pkg/apiClient"
	"context"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type ConnectionPool struct {
	numberOfConnections int
	client              *apiClient.TwelveDataClient
	wg                  sync.WaitGroup
	connections         chan connection
	stopped             bool
}

type connection struct {
	id int
}

func NewTwelveDataPool(apiKey string, apiHost string, clientTimout time.Duration, connectionNumber int) *ConnectionPool {
	pool := &ConnectionPool{
		numberOfConnections: connectionNumber,
		client:              apiClient.NewTwelveDataClient(apiKey, apiHost, clientTimout),
		wg:                  sync.WaitGroup{},
	}
	pool.init()
	log.Info().Msgf("Connection pool initialized with %d connections", connectionNumber)
	return pool
}

func (p *ConnectionPool) init() {
	p.connections = make(chan connection, p.numberOfConnections)
	for i := 0; i < p.numberOfConnections; i++ {
		p.connections <- connection{i + 1}
	}
}

func (p *ConnectionPool) Stop() {
	p.stopped = true
	p.wg.Wait()
}

func (p *ConnectionPool) GetHistoricDataForSymbol(ctx context.Context, symbol string) (*apiClient.TimeSeries, error) {
	con := <-p.connections
	p.wg.Add(1)
	result, err := p.client.GetHistoricDataForSymbol(ctx, symbol)
	p.wg.Done()
	go p.restoreConnection(con)
	return result, err
}

func (p *ConnectionPool) restoreConnection(c connection) {
	time.Sleep(1 * time.Minute)
	p.connections <- c
	log.Info().Msgf("Connection #%d restored", c.id)
}
