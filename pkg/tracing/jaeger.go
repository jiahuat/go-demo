package tracing

import (
	"bytes"
	"context"
	"io/ioutil"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// Jaeger 通过 middleware 将 tracer 和 ctx 注入到 gin.Context 中
func Jaeger() gin.HandlerFunc {
	return func(c *gin.Context) {
		var parentSpan opentracing.Span
		tracer, closer := NewTracer("wscp-restful-demo")
		defer closer.Close()
		// 直接从 c.Request.Header 中提取 span,如果没有就新建一个
		spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			parentSpan = tracer.StartSpan(c.Request.URL.Path)
			defer parentSpan.Finish()
		} else {
			parentSpan = opentracing.StartSpan(
				c.Request.URL.Path,
				opentracing.ChildOf(spCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				ext.SpanKindRPCServer,
				ext.RPCServerOption(spCtx),
			)
			defer parentSpan.Finish()
		}

		// 记录请求 Url
		ext.HTTPUrl.Set(parentSpan, c.Request.URL.Path)
		// Http Method
		ext.HTTPMethod.Set(parentSpan, c.Request.Method)
		// 记录组件名称
		ext.Component.Set(parentSpan, "Gin-Http")
		// 自定义 Tag X-Forwarded-For
		opentracing.Tag{Key: "http.headers.x-forwarded-for", Value: c.Request.Header.Get("X-Forwarded-For")}.Set(parentSpan)
		// 自定义 Tag User-Agent
		opentracing.Tag{Key: "http.headers.user-agent", Value: c.Request.Header.Get("User-Agent")}.Set(parentSpan)
		// 自定义 Tag Request-Time
		opentracing.Tag{Key: "request.time", Value: time.Now().Format(time.RFC3339)}.Set(parentSpan)
		// 自定义 Tag Server-Mode
		opentracing.Tag{Key: "http.server.mode", Value: gin.Mode()}.Set(parentSpan)

		// read http request body
		body, err := c.GetRawData() // 读取 request body 的内容
		if err != nil {
			body = []byte("failed to get request body")
		} else {
			// 自定义 Tag Request-Body
			opentracing.Tag{Key: "http.request_body", Value: string(body)}.Set(parentSpan)
		}
		// set request body back
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), parentSpan))
		// 然后存到 g.ctx 中 供后续使用
		c.Set("tracer", tracer)
		c.Set("ctx", opentracing.ContextWithSpan(context.Background(), parentSpan))

		c.Next()

		// resp, err := ioutil.ReadAll(c.Request.Response.Body)
		// if err == nil {
		// 	opentracing.Tag{Key: "http.response_body", Value: string(resp)}.Set(parentSpan)
		// }

		if gin.Mode() == gin.DebugMode {
			// 自定义 Tag StackTrace
			opentracing.Tag{Key: "debug.trace", Value: string(debug.Stack())}.Set(parentSpan)
		}

		ext.HTTPStatusCode.Set(parentSpan, uint16(c.Writer.Status()))
		opentracing.Tag{Key: "request.errors", Value: c.Errors.String()}.Set(parentSpan)
	}
}
