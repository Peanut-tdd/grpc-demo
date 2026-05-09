package middleware

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

func TimeoutStreamInterceptor(timeout time.Duration) grpc.StreamServerInterceptor {

	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx := ss.Context()

		// 如果客户端已经设置了 deadline，则不再覆盖
		if _, ok := ctx.Deadline(); ok {
			return handler(srv, ss)
		}

		// 设置服务端默认超时
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// 包装 ServerStream，注入新的 context
		wrapped := &wrappedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		return handler(srv, wrapped)
	}
}

// 包装 ServerStream
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

func TimeoutUnaryInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// 如果客户端已经设置了 deadline，就不覆盖
		if _, ok := ctx.Deadline(); ok {
			return handler(ctx, req)
		}

		// 设置服务端默认超时
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return handler(ctx, req)
	}
}
