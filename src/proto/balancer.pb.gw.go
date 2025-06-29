// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: src/proto/balancer.proto

/*
Package balancer is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package proto

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var (
	_ codes.Code
	_ io.Reader
	_ status.Status
	_ = errors.New
	_ = runtime.String
	_ = utilities.NewDoubleArray
	_ = metadata.Join
)

func request_VideoBalancer_GetRedirect_0(ctx context.Context, marshaler runtime.Marshaler, client VideoBalancerClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq VideoRequest
		metadata runtime.ServerMetadata
	)
	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && !errors.Is(err, io.EOF) {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	msg, err := client.GetRedirect(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
}

func local_request_VideoBalancer_GetRedirect_0(ctx context.Context, marshaler runtime.Marshaler, server VideoBalancerServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq VideoRequest
		metadata runtime.ServerMetadata
	)
	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && !errors.Is(err, io.EOF) {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	msg, err := server.GetRedirect(ctx, &protoReq)
	return msg, metadata, err
}

// RegisterVideoBalancerHandlerServer registers the http handlers for service VideoBalancer to "mux".
// UnaryRPC     :call VideoBalancerServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterVideoBalancerHandlerFromEndpoint instead.
// GRPC interceptors will not work for this type of registration. To use interceptors, you must use the "runtime.WithMiddlewares" option in the "runtime.NewServeMux" call.
func RegisterVideoBalancerHandlerServer(ctx context.Context, mux *runtime.ServeMux, server VideoBalancerServer) error {
	mux.Handle(http.MethodPost, pattern_VideoBalancer_GetRedirect_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/balancer.VideoBalancer/GetRedirect", runtime.WithHTTPPathPattern("/watch"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_VideoBalancer_GetRedirect_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_VideoBalancer_GetRedirect_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})

	return nil
}

// RegisterVideoBalancerHandlerFromEndpoint is same as RegisterVideoBalancerHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterVideoBalancerHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.NewClient(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()
	return RegisterVideoBalancerHandler(ctx, mux, conn)
}

// RegisterVideoBalancerHandler registers the http handlers for service VideoBalancer to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterVideoBalancerHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterVideoBalancerHandlerClient(ctx, mux, NewVideoBalancerClient(conn))
}

// RegisterVideoBalancerHandlerClient registers the http handlers for service VideoBalancer
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "VideoBalancerClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "VideoBalancerClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "VideoBalancerClient" to call the correct interceptors. This client ignores the HTTP middlewares.
func RegisterVideoBalancerHandlerClient(ctx context.Context, mux *runtime.ServeMux, client VideoBalancerClient) error {
	mux.Handle(http.MethodPost, pattern_VideoBalancer_GetRedirect_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateContext(ctx, mux, req, "/balancer.VideoBalancer/GetRedirect", runtime.WithHTTPPathPattern("/watch"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_VideoBalancer_GetRedirect_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_VideoBalancer_GetRedirect_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	return nil
}

var (
	pattern_VideoBalancer_GetRedirect_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0}, []string{"watch"}, ""))
)

var (
	forward_VideoBalancer_GetRedirect_0 = runtime.ForwardResponseMessage
)
