package connectionPool

import (
	"FinanceApi/pkg/apiClient"
	"FinanceApi/pkg/config"
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type ConnectionPool struct {
	context             context.Context
	numberOfConnections int
	client              *apiClient.TwelveDataClient
	wg                  sync.WaitGroup
	connections         chan connection
	stopped             bool
}

type connection struct {
	id int
}

func NewTwelveDataPool(context context.Context) *ConnectionPool {
	pool := &ConnectionPool{
		context:             context,
		numberOfConnections: config.Conf.API.TwelveData.RateLimit,
		client:              apiClient.NewTwelveDataClient(),
		wg:                  sync.WaitGroup{},
	}
	pool.init()
	go gracefulShutdown(pool)
	return pool
}

func (p *ConnectionPool) init() {
	p.connections = make(chan connection, p.numberOfConnections)
	for i := 0; i < p.numberOfConnections; i++ {
		p.connections <- connection{i + 1}
	}
}

func gracefulShutdown(pool *ConnectionPool) {
	stopWg := pool.context.Value("stopWg").(*sync.WaitGroup)
	stopWg.Add(1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	pool.stop()
	stopWg.Done()
}

func (p *ConnectionPool) stop() {
	p.stopped = true
	p.wg.Wait()
}

func (p *ConnectionPool) GetHistoricDataForSymbol(symbol string) (*apiClient.TimeSeries, error) {
	con := <-p.connections
	p.wg.Add(1)
	result, err := p.client.GetHistoricDataForSymbol(symbol)
	p.wg.Done()
	go p.restoreConnection(con)
	return result, err
}

func (p *ConnectionPool) restoreConnection(c connection) {
	time.Sleep(1 * time.Minute)
	p.connections <- c
	log.Printf("Connection #%d restored\n", c.id)
}
