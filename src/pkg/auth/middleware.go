/**
 * @Time: 2019-06-29 09:59
 * @Author: solacowa@gmail.com
 * @File: middleware
 * @Software: GoLand
 */

package auth

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	captcha "github.com/icowan/kit-captcha"
	"github.com/kplcloud/kplcloud/src/encode"
)

// 验证图形验证码
func checkCaptchaMiddleware(captcha captcha.Service) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			var captchaId, captchaCode string
			switch request.(type) {
			case authRequest:
				req, ok := request.(authRequest)
				if !ok {
					return nil, encode.ErrSystem.Error()
				}
				captchaId = req.CaptchaId
				captchaCode = req.CaptchaCode
			default:
				return next(ctx, request)
			}

			if captchaId == "" || captchaCode == "" {
				return nil, encode.ErrAuthCheckCaptchaNotnull.Error()
			}
			if !captcha.VerifyCaptcha(ctx, captchaId, captchaCode) {
				return nil, encode.ErrAuthCheckCaptchaCode.Error()
			}

			return next(ctx, request)
		}
	}
}
