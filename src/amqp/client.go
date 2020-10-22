/**
 * @Time : 2019-06-25 17:52
 * @Author : solacowa@gmail.com
 * @File : amqp
 * @Software: GoLand
 */

package amqp

import (
	"bytes"
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitamqp "github.com/go-kit/kit/transport/amqp"
	"github.com/icowan/config"
	"github.com/streadway/amqp"
	"time"
)

type QueueName string

const (
	CronJobTopic   QueueName = "cronjob.kplcloud"
	BuildTopic     QueueName = "build.kplcloud"
	AlarmTopic     QueueName = "alarm.kplcloud"
	NoticeTopic    QueueName = "notice.kplcloud"
	ProclaimTopic  QueueName = "proclaim.kplcloud"
	MsgWechatTopic QueueName = "msg.wechat.kplcloud"
	HookTopic      QueueName = "hook.kplcloud"
)

func (c QueueName) String() string {
	return string(c)
}

type AmqpClient interface {
	Close() error
	PublishOnQueue(queueName QueueName, data func() []byte) (err error)
	SubscribeToQueue(ctx context.Context, logger log.Logger, call func(ctx context.Context, data string) error)
}

type amqpClient struct {
	conn         *amqp.Connection
	exchange     string
	exchangeType string
	routingKey   string
}

func NewAmqp(cf *config.Config) (AmqpClient, error) {

	conn, err := amqp.Dial(cf.GetString("amqp", "url"))
	if err != nil {
		return nil, err
	}
	return &amqpClient{
		conn:         conn,
		exchange:     cf.GetString("amqp", "exchange_type"),
		exchangeType: cf.GetString("amqp", "exchange"),
		routingKey:   cf.GetString("amqp", "routing_key"),
	}, nil
}

func (c *amqpClient) PublishOnQueue(queueName QueueName, data func() []byte) (err error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return
	}

	defer func() {
		_ = ch.Close()
	}()

	queue, err := ch.QueueDeclare(queueName.String(), false, false, false, false, nil)
	if err != nil {
		return
	}

	if err = ch.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Transient,
		Priority:     0,
		ContentType:  "application/json",
		Body:         data(),
	}); err != nil {
		return
	}

	//ep := makePublishOnQueueEndpoint(ch, &queue)
	//ctx := context.WithValue(context.Background(), kitamqp.ContextKeyPublishKey, queueName.String())
	//ctx = context.WithValue(ctx, kitamqp.ContextKeyExchange, "")
	//
	//errChan := make(chan error, 1)
	//resChan := make(chan interface{}, 1)
	//
	//go func() {
	//	res, err := ep(ctx, data())
	//	if err != nil {
	//		errChan <- err
	//	} else {
	//		resChan <- res
	//	}
	//}()
	//
	//select {
	//case <-resChan:
	//	break
	//
	//case err = <-errChan:
	//	break
	//
	//case <-time.After(5 * time.Second):
	//	return errors.New("timed out waiting for result")
	//}

	return err
}

func (c *amqpClient) SubscribeToQueue(ctx context.Context, logger log.Logger, call func(ctx context.Context, data string) error) {
	ch, err := c.conn.Channel()
	if err != nil {
		_ = level.Error(logger).Log("conn", "Channel", "err", err.Error())
		return
	}

	defer func() {
		_ = ch.Close()
	}()

	topic, _ := ctx.Value(kitamqp.ContextKeyPublishKey).(string)

	// decoder
	sub := kitamqp.NewSubscriber(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// 可以在这里处理入参
		return
	}, func(i context.Context, delivery *amqp.Delivery) (request interface{}, err error) {
		// 处理出参
		return delivery.Body, nil
	}, func(i context.Context, publishing *amqp.Publishing, i2 interface{}) error {
		return nil
	},
		kitamqp.SubscriberErrorEncoder(kitamqp.ReplyErrorEncoder),
		kitamqp.SubscriberBefore(
			kitamqp.SetPublishKey(topic),
			//kitamqp.SetPublishDeliveryMode(Delivery),
			kitamqp.SetContentType("application/json"),
			//kitamqp.SetContentEncoding(contentEncoding),
		))

	outputChan := make(chan amqp.Publishing, 1)

	q, err := ch.QueueDeclare(topic, false, false, false, false, nil)
	if err != nil {
		_ = level.Error(logger).Log("ch", "QueueDeclare", "err", err.Error())
		return
	}

	deliveries, err := ch.Consume(q.Name, topic, true, false, false, false, nil)
	if err != nil {
		_ = level.Error(logger).Log("ch", "Consume", "err", err.Error())
		return
	}

	if err = ch.QueueBind(topic, c.routingKey, c.exchange, false, nil); err != nil {
		_ = level.Error(logger).Log("ch", "QueueBind", "err", err.Error())
		return
	}

	sub.ServeDelivery(ch)(&amqp.Delivery{})

	var msg amqp.Publishing

	select {
	case msg = <-outputChan:
		break

	case <-time.After(100 * time.Millisecond):
		_ = level.Error(logger).Log("time", "after", "err", "Timed out waiting for publishing")
	}

	_ = level.Debug(logger).Log("msgBody", string(msg.Body))

	forever := make(chan bool)
	go func() {
		for d := range deliveries {
			s := BytesToString(&(d.Body))
			if err = call(context.Background(), *s); err != nil {
				_ = level.Error(logger).Log("svc", "ReceiverBuild", "err", err.Error())
			}
		}
	}()

	<-forever
}

func (c *amqpClient) Close() error {
	return c.conn.Close()
}

func BytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}
