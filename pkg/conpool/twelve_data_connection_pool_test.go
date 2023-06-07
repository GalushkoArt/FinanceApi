package conpool

import (
	"context"
	"github.com/galushkoart/finance-api/mock"
	"github.com/galushkoart/finance-api/pkg/apiclient"
	"github.com/golang/mock/gomock"
	"strconv"
	"sync"
	"testing"
	"time"
)

//go:generate echo $PWD - $GOFILE
//go:generate mockgen -package mock -destination ../../mock/twelve_data_client_mock.go -source=../apiclient/twelve_data_client.go TwelveDataClient

func TestTwelveDataPool(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockClient := mock.NewMockTwelveDataClient(controller)
	pool := &ConnectionPool{
		numberOfConnections: 10,
		client:              mockClient,
		wg:                  sync.WaitGroup{},
		restoreTime:         1 * time.Second,
	}
	explicitWait := &sync.WaitGroup{}
	pool.init()
	mockClient.EXPECT().GetHistoricDataForSymbol(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, symbol string) (*apiclient.TimeSeries, error) {
		return &apiclient.TimeSeries{Meta: apiclient.Meta{Symbol: symbol}}, nil
	}).Times(pool.numberOfConnections)
	explicitWait.Add(pool.numberOfConnections)
	for i := 0; i < pool.numberOfConnections+2; i++ {
		i := i
		go func() {
			symbolId := strconv.Itoa(i)
			symbol, err := pool.GetHistoricDataForSymbol(context.TODO(), symbolId)
			explicitWait.Done()
			if err != nil {
				t.Errorf("Found unexpected error on api call: %v", err)
				return
			}
			if symbol.Meta.Symbol != symbolId {
				t.Errorf("Found unexpected symbol: %v", symbol.Meta)
				return
			}
		}()
	}
	explicitWait.Wait()
	pool.Stop()
}
