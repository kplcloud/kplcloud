/**
 * @Time : 2019-07-10 13:56
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package amqp

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kitamqp "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	"reflect"
	"time"
)

type BuildPublishData struct {
	Name           string `json:"name"`             // name
	Namespace      string `json:"namespace"`        // namespace
	JenkinJobName  string `json:"jenkin_job_name"`  // jenkins job name
	BuildId        int64  `json:"build_id"`         // build id
	JenkinsBuildId int64  `json:"jenkins_build_id"` // jenkins build id
	Builder        int64  `json:"builder"`          // 操作人
	Timestamp      int64  `json:"timestamp"`        //创建时间的时间戳
	Version        string `json:"version"`
}

type PublishData struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func makePublishOnQueueEndpoint(ch *amqp.Channel, queue *amqp.Queue) endpoint.Endpoint {
	cid := "correlation"

	return kitamqp.NewPublisher(ch, queue, func(ctx context.Context, publishing *amqp.Publishing, request interface{}) error {
		var body []byte
		switch reflect.TypeOf(request).String() {
		case "string":
			body = []byte(request.(string))
		case "[]byte":
			body = request.([]byte)
		case "[]uint8":
			for _, b := range request.([]uint8) {
				body = append(body, byte(b))
			}
		default:
			body, _ = json.Marshal(request)
		}

		ctx = context.WithValue(ctx, kitamqp.ContextKeyPublishKey, queue.Name)
		ctx = context.WithValue(ctx, kitamqp.ContextKeyExchange, "")
		publishing.DeliveryMode = amqp.Transient
		publishing.Priority = 0
		publishing.ContentType = "application/json"
		publishing.Body = body
		return nil
	}, func(i context.Context, delivery *amqp.Delivery) (response interface{}, err error) {
		return delivery.Body, nil
	},
		kitamqp.PublisherTimeout(5*time.Second),
		kitamqp.PublisherBefore(kitamqp.SetCorrelationID(cid)),
	).Endpoint()
}
