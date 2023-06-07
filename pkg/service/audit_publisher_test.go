package service

import (
	"context"
	audit "github.com/GalushkoArt/GoAuditService/pkg/proto"
	"github.com/galushkoart/finance-api/mock"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"strconv"
	"sync"
	"testing"
)

//go:generate echo $PWD - $GOFILE
//go:generate mockgen -package mock -destination ../../mock/publish_channel_mock.go -source=audit_publisher.go PublishChannel

func TestMQAuditPublisher(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockConnection := mock.NewMockCloser(controller)
	mockConnection.EXPECT().Close().Return(nil)
	mockPublishChannel := mock.NewMockPublishChannel(controller)
	mockPublishChannel.EXPECT().Close().Return(nil)
	publisher := MQAuditPublisher{
		queueName:      "test",
		connection:     mockConnection,
		publishChannel: mockPublishChannel,
		wg:             &sync.WaitGroup{},
	}
	testData := make([]*audit.LogRequest, 0, 50)
	for i := 0; i < cap(testData); i++ {
		request := &audit.LogRequest{RequestId: strconv.Itoa(i)}
		protoRequest, err := proto.Marshal(request)
		if err != nil {
			t.Fatalf("Found unexpected error on proto marshal: %v", err)
		}
		testData = append(testData, request)
		mockPublishChannel.EXPECT().PublishWithContext(gomock.Any(), "", "test", false, false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        protoRequest,
			}).Return(nil)
	}
	explicitWait := &sync.WaitGroup{}
	for _, request := range testData {
		request := request
		explicitWait.Add(1)
		go func() {
			err := publisher.Publish(context.TODO(), request)
			explicitWait.Done()
			if err != nil {
				t.Errorf("Found unexpected error on api call: %v", err)
				return
			}
		}()
	}
	explicitWait.Wait()
	if err := publisher.Close(); err != nil {
		t.Fatalf("Found unexpected error on close: %v", err)
	}
}
