/**
 * @Time : 2019-07-23 18:48
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package tools

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"strings"
	"time"
)

type FakeTimeMethod string

const (
	FakeTimeAdd   FakeTimeMethod = "add"
	FakeTimeClean FakeTimeMethod = "del"
)

type duplicationRequest struct {
	SourceNamespace      string `json:"source_namespace"`
	SourceAppName        string `json:"source_app_name"`
	DestinationNamespace string `json:"destination_namespace"`
}

type fakeTimeRequest struct {
	FakeTime Time           `json:"fake_time"`
	Method   FakeTimeMethod `json:"method"`
}

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), "\"")
	fakeTime, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
	if err != nil {
		return err
	}
	*t = Time{fakeTime}
	return nil
}

func makeDuplicationEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(duplicationRequest)
		err := s.Duplication(ctx, req.SourceNamespace, req.SourceAppName, req.DestinationNamespace)
		return encode.Response{Err: err}, err
	}
}

func makeFakeTimeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(fakeTimeRequest)
		err := s.FakeTime(ctx, req.FakeTime.Time, req.Method)
		return encode.Response{Err: err}, err
	}
}
