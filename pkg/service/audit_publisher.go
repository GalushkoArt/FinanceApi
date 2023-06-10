package service

import (
	"context"
	audit "github.com/GalushkoArt/GoAuditService/pkg/proto"
	"github.com/galushkoart/finance-api/pkg/utils"
	"github.com/golang/protobuf/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"io"
	"sync"
)

type MQAuditPublisher struct {
	publishChannel PublishChannel
	queueName      string
	connection     io.Closer
	wg             *sync.WaitGroup
}

type AuditPublisher interface {
	Publish(ctx context.Context, request *audit.LogRequest) error
}

type PublishChannel interface {
	PublishWithContext(ctx context.Context, destination string, msg proto.Message) error
	io.Closer
}

type rabbitQueue struct {
	publishChannel *amqp.Channel
}

func newRabbitQueue(conn *amqp.Connection) (*rabbitQueue, error) {
	channel, err := conn.Channel()
	return &rabbitQueue{publishChannel: channel}, err
}

func (q *rabbitQueue) PublishWithContext(ctx context.Context, destination string, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	return q.publishChannel.PublishWithContext(
		ctx,
		"",
		destination,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
}

func (q *rabbitQueue) Close() error {
	return q.publishChannel.Close()
}

func NewAuditPublisher(queueName string) *MQAuditPublisher {
	return &MQAuditPublisher{queueName: queueName, wg: &sync.WaitGroup{}}
}

func (p *MQAuditPublisher) InitPublishChannel(enabled bool, brokerUri string) func() error {
	if !enabled {
		log.Info().Msg("Audit publisher disabled!")
		return func() error {
			return nil
		}
	}
	conn, err := amqp.Dial(brokerUri)
	utils.PanicOnError(err)

	p.publishChannel, err = newRabbitQueue(conn)
	utils.PanicOnError(err)
	p.connection = conn
	log.Info().Msg("Audit publisher started!")
	return p.Close
}

func (p *MQAuditPublisher) Close() error {
	p.wg.Wait()
	err := p.publishChannel.Close()
	if err != nil {
		return err
	}
	return p.connection.Close()
}

func (p *MQAuditPublisher) Publish(ctx context.Context, request *audit.LogRequest) error {
	p.wg.Add(1)
	err := p.publishChannel.PublishWithContext(ctx, p.queueName, request)
	p.wg.Done()
	return err
}
