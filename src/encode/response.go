/**
 * @Time: 2020/2/14 13:57
 * @Author: solacowa@gmail.com
 * @File: encode
 * @Software: GoLand
 */

package encode

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
)

type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"message,omitempty"`
}

type Failure interface {
	Failed() error
}

type Errorer interface {
	Error() error
}

func Error(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(err.Error()))
}

func JsonError(ctx context.Context, err error, w http.ResponseWriter) {
	headers, ok := ctx.Value("response-headers").(map[string]string)
	if ok {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
	}

	var errDefined bool
	for k := range ResponseMessage {
		if strings.Contains(err.Error(), k.Error().Error()) {
			errDefined = true
			break
		}
	}

	if !errDefined {
		err = ErrSystem.Error()
	}
	if err == nil {
		err = errors.Wrap(err, ErrSystem.Error().Error())
	}

	_ = kithttp.EncodeJSONResponse(ctx, w, map[string]interface{}{
		"message": err.Error(),
		"code":    ResponseMessage[ResStatus(strings.Split(err.Error(), ":")[0])],
		"success": false,
	})
}

func JsonResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(Failure); ok && f.Failed() != nil {
		JsonError(ctx, f.Failed(), w)
		return nil
	}
	resp := response.(Response)
	if resp.Error == nil {
		resp.Code = 200
		resp.Success = true
	}

	headers, ok := ctx.Value("response-headers").(map[string]string)
	if ok {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
	}

	return kithttp.EncodeJSONResponse(ctx, w, resp)
}

func XmlResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	return nil
}
