package service

import (
	"FinanceApi/pkg/utils"
	"context"
	audit "github.com/GalushkoArt/GoAuditService/pkg/proto"
	"github.com/golang/protobuf/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"sync"
)

type MQAuditPublisher struct {
	publishChannel *amqp.Channel
	queueName      string
	wg             *sync.WaitGroup
}

type AuditPublisher interface {
	Publish(ctx context.Context, request *audit.LogRequest) error
}

func NewAuditPublisher(queueName string) *MQAuditPublisher {
	return &MQAuditPublisher{queueName: queueName, wg: &sync.WaitGroup{}}
}

func (p *MQAuditPublisher) InitPublishChannel(enabled bool, brokerUri string) func() {
	if !enabled {
		log.Info().Msg("Audit publisher disabled!")
		return func() {}
	}
	conn, err := amqp.Dial(brokerUri)
	utils.PanicOnError(err)

	p.publishChannel, err = conn.Channel()
	utils.PanicOnError(err)
	log.Info().Msg("Audit publisher started!")
	return func() {
		p.wg.Wait()
		utils.PanicOnError(p.publishChannel.Close())
		utils.PanicOnError(conn.Close())
	}
}

func (p *MQAuditPublisher) Publish(ctx context.Context, request *audit.LogRequest) error {
	p.wg.Add(1)
	data, err := proto.Marshal(request)
	if err != nil {
		p.wg.Done()
		return err
	}

	err = p.publishChannel.PublishWithContext(
		ctx,
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
	p.wg.Done()
	return err
}
