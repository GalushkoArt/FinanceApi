package service

import (
	"context"
	audit "github.com/GalushkoArt/GoAuditService/pkg/proto"
	"github.com/galushkoart/finance-api/mock"
	"github.com/golang/mock/gomock"
	"strconv"
	"sync"
	"testing"
)

//go:generate echo $PWD - $GOFILE
//go:generate mockgen -package mock -destination ../../mock/audit_service_client_mock.go github.com/GalushkoArt/GoAuditService/pkg/proto AuditServiceClient
//go:generate mockgen -package mock -destination ../../mock/closer_mock.go io Closer

func TestAuditClient(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockService := mock.NewMockAuditServiceClient(controller)
	mockConn := mock.NewMockCloser(controller)
	mockConn.EXPECT().Close().Return(nil)
	explicitWait := &sync.WaitGroup{}
	auditClient := AuditClient{client: mockService, conn: mockConn, wg: &sync.WaitGroup{}}
	testData := make([]*audit.LogRequest, 0, 50)
	for i := 0; i < cap(testData); i++ {
		request := &audit.LogRequest{RequestId: strconv.Itoa(i)}
		testData = append(testData, request)
		mockService.EXPECT().Log(gomock.Any(), request, gomock.Any()).Return(&audit.Response{Answer: audit.Response_SUCCESS}, nil)
	}
	for _, request := range testData {
		request := request
		explicitWait.Add(1)
		go func() {
			resp, err := auditClient.SendRequest(context.Background(), request)
			explicitWait.Done()
			if err != nil {
				t.Errorf("Found unexpected error on api call: %v", err)
				return
			}
			if resp.Answer != audit.Response_SUCCESS {
				t.Errorf("Found unexpected response: %v", resp.Answer)
				return
			}
		}()
	}
	explicitWait.Wait()
	if err := auditClient.Close(); err != nil {
		t.Fatalf("Found unexpected error on close: %v", err)
	}
}
