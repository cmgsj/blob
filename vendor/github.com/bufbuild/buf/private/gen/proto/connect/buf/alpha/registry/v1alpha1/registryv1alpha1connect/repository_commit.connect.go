// Copyright 2020-2024 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: buf/alpha/registry/v1alpha1/repository_commit.proto

package registryv1alpha1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/registry/v1alpha1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// RepositoryCommitServiceName is the fully-qualified name of the RepositoryCommitService service.
	RepositoryCommitServiceName = "buf.alpha.registry.v1alpha1.RepositoryCommitService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// RepositoryCommitServiceListRepositoryCommitsByBranchProcedure is the fully-qualified name of the
	// RepositoryCommitService's ListRepositoryCommitsByBranch RPC.
	RepositoryCommitServiceListRepositoryCommitsByBranchProcedure = "/buf.alpha.registry.v1alpha1.RepositoryCommitService/ListRepositoryCommitsByBranch"
	// RepositoryCommitServiceListRepositoryCommitsByReferenceProcedure is the fully-qualified name of
	// the RepositoryCommitService's ListRepositoryCommitsByReference RPC.
	RepositoryCommitServiceListRepositoryCommitsByReferenceProcedure = "/buf.alpha.registry.v1alpha1.RepositoryCommitService/ListRepositoryCommitsByReference"
	// RepositoryCommitServiceGetRepositoryCommitByReferenceProcedure is the fully-qualified name of the
	// RepositoryCommitService's GetRepositoryCommitByReference RPC.
	RepositoryCommitServiceGetRepositoryCommitByReferenceProcedure = "/buf.alpha.registry.v1alpha1.RepositoryCommitService/GetRepositoryCommitByReference"
	// RepositoryCommitServiceListRepositoryDraftCommitsProcedure is the fully-qualified name of the
	// RepositoryCommitService's ListRepositoryDraftCommits RPC.
	RepositoryCommitServiceListRepositoryDraftCommitsProcedure = "/buf.alpha.registry.v1alpha1.RepositoryCommitService/ListRepositoryDraftCommits"
	// RepositoryCommitServiceDeleteRepositoryDraftCommitProcedure is the fully-qualified name of the
	// RepositoryCommitService's DeleteRepositoryDraftCommit RPC.
	RepositoryCommitServiceDeleteRepositoryDraftCommitProcedure = "/buf.alpha.registry.v1alpha1.RepositoryCommitService/DeleteRepositoryDraftCommit"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	repositoryCommitServiceServiceDescriptor                                = v1alpha1.File_buf_alpha_registry_v1alpha1_repository_commit_proto.Services().ByName("RepositoryCommitService")
	repositoryCommitServiceListRepositoryCommitsByBranchMethodDescriptor    = repositoryCommitServiceServiceDescriptor.Methods().ByName("ListRepositoryCommitsByBranch")
	repositoryCommitServiceListRepositoryCommitsByReferenceMethodDescriptor = repositoryCommitServiceServiceDescriptor.Methods().ByName("ListRepositoryCommitsByReference")
	repositoryCommitServiceGetRepositoryCommitByReferenceMethodDescriptor   = repositoryCommitServiceServiceDescriptor.Methods().ByName("GetRepositoryCommitByReference")
	repositoryCommitServiceListRepositoryDraftCommitsMethodDescriptor       = repositoryCommitServiceServiceDescriptor.Methods().ByName("ListRepositoryDraftCommits")
	repositoryCommitServiceDeleteRepositoryDraftCommitMethodDescriptor      = repositoryCommitServiceServiceDescriptor.Methods().ByName("DeleteRepositoryDraftCommit")
)

// RepositoryCommitServiceClient is a client for the
// buf.alpha.registry.v1alpha1.RepositoryCommitService service.
type RepositoryCommitServiceClient interface {
	// ListRepositoryCommitsByBranch lists the repository commits associated
	// with a repository branch on a repository, ordered by their create time.
	//
	// Deprecated: do not use.
	ListRepositoryCommitsByBranch(context.Context, *connect.Request[v1alpha1.ListRepositoryCommitsByBranchRequest]) (*connect.Response[v1alpha1.ListRepositoryCommitsByBranchResponse], error)
	// ListRepositoryCommitsByReference returns repository commits up-to and including
	// the provided reference.
	ListRepositoryCommitsByReference(context.Context, *connect.Request[v1alpha1.ListRepositoryCommitsByReferenceRequest]) (*connect.Response[v1alpha1.ListRepositoryCommitsByReferenceResponse], error)
	// GetRepositoryCommitByReference returns the repository commit matching
	// the provided reference, if it exists.
	GetRepositoryCommitByReference(context.Context, *connect.Request[v1alpha1.GetRepositoryCommitByReferenceRequest]) (*connect.Response[v1alpha1.GetRepositoryCommitByReferenceResponse], error)
	// ListRepositoryDraftCommits lists draft commits in a repository.
	ListRepositoryDraftCommits(context.Context, *connect.Request[v1alpha1.ListRepositoryDraftCommitsRequest]) (*connect.Response[v1alpha1.ListRepositoryDraftCommitsResponse], error)
	// DeleteRepositoryDraftCommit deletes a draft.
	DeleteRepositoryDraftCommit(context.Context, *connect.Request[v1alpha1.DeleteRepositoryDraftCommitRequest]) (*connect.Response[v1alpha1.DeleteRepositoryDraftCommitResponse], error)
}

// NewRepositoryCommitServiceClient constructs a client for the
// buf.alpha.registry.v1alpha1.RepositoryCommitService service. By default, it uses the Connect
// protocol with the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed
// requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewRepositoryCommitServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) RepositoryCommitServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &repositoryCommitServiceClient{
		listRepositoryCommitsByBranch: connect.NewClient[v1alpha1.ListRepositoryCommitsByBranchRequest, v1alpha1.ListRepositoryCommitsByBranchResponse](
			httpClient,
			baseURL+RepositoryCommitServiceListRepositoryCommitsByBranchProcedure,
			connect.WithSchema(repositoryCommitServiceListRepositoryCommitsByBranchMethodDescriptor),
			connect.WithIdempotency(connect.IdempotencyNoSideEffects),
			connect.WithClientOptions(opts...),
		),
		listRepositoryCommitsByReference: connect.NewClient[v1alpha1.ListRepositoryCommitsByReferenceRequest, v1alpha1.ListRepositoryCommitsByReferenceResponse](
			httpClient,
			baseURL+RepositoryCommitServiceListRepositoryCommitsByReferenceProcedure,
			connect.WithSchema(repositoryCommitServiceListRepositoryCommitsByReferenceMethodDescriptor),
			connect.WithIdempotency(connect.IdempotencyNoSideEffects),
			connect.WithClientOptions(opts...),
		),
		getRepositoryCommitByReference: connect.NewClient[v1alpha1.GetRepositoryCommitByReferenceRequest, v1alpha1.GetRepositoryCommitByReferenceResponse](
			httpClient,
			baseURL+RepositoryCommitServiceGetRepositoryCommitByReferenceProcedure,
			connect.WithSchema(repositoryCommitServiceGetRepositoryCommitByReferenceMethodDescriptor),
			connect.WithIdempotency(connect.IdempotencyNoSideEffects),
			connect.WithClientOptions(opts...),
		),
		listRepositoryDraftCommits: connect.NewClient[v1alpha1.ListRepositoryDraftCommitsRequest, v1alpha1.ListRepositoryDraftCommitsResponse](
			httpClient,
			baseURL+RepositoryCommitServiceListRepositoryDraftCommitsProcedure,
			connect.WithSchema(repositoryCommitServiceListRepositoryDraftCommitsMethodDescriptor),
			connect.WithIdempotency(connect.IdempotencyNoSideEffects),
			connect.WithClientOptions(opts...),
		),
		deleteRepositoryDraftCommit: connect.NewClient[v1alpha1.DeleteRepositoryDraftCommitRequest, v1alpha1.DeleteRepositoryDraftCommitResponse](
			httpClient,
			baseURL+RepositoryCommitServiceDeleteRepositoryDraftCommitProcedure,
			connect.WithSchema(repositoryCommitServiceDeleteRepositoryDraftCommitMethodDescriptor),
			connect.WithIdempotency(connect.IdempotencyIdempotent),
			connect.WithClientOptions(opts...),
		),
	}
}

// repositoryCommitServiceClient implements RepositoryCommitServiceClient.
type repositoryCommitServiceClient struct {
	listRepositoryCommitsByBranch    *connect.Client[v1alpha1.ListRepositoryCommitsByBranchRequest, v1alpha1.ListRepositoryCommitsByBranchResponse]
	listRepositoryCommitsByReference *connect.Client[v1alpha1.ListRepositoryCommitsByReferenceRequest, v1alpha1.ListRepositoryCommitsByReferenceResponse]
	getRepositoryCommitByReference   *connect.Client[v1alpha1.GetRepositoryCommitByReferenceRequest, v1alpha1.GetRepositoryCommitByReferenceResponse]
	listRepositoryDraftCommits       *connect.Client[v1alpha1.ListRepositoryDraftCommitsRequest, v1alpha1.ListRepositoryDraftCommitsResponse]
	deleteRepositoryDraftCommit      *connect.Client[v1alpha1.DeleteRepositoryDraftCommitRequest, v1alpha1.DeleteRepositoryDraftCommitResponse]
}

// ListRepositoryCommitsByBranch calls
// buf.alpha.registry.v1alpha1.RepositoryCommitService.ListRepositoryCommitsByBranch.
//
// Deprecated: do not use.
func (c *repositoryCommitServiceClient) ListRepositoryCommitsByBranch(ctx context.Context, req *connect.Request[v1alpha1.ListRepositoryCommitsByBranchRequest]) (*connect.Response[v1alpha1.ListRepositoryCommitsByBranchResponse], error) {
	return c.listRepositoryCommitsByBranch.CallUnary(ctx, req)
}

// ListRepositoryCommitsByReference calls
// buf.alpha.registry.v1alpha1.RepositoryCommitService.ListRepositoryCommitsByReference.
func (c *repositoryCommitServiceClient) ListRepositoryCommitsByReference(ctx context.Context, req *connect.Request[v1alpha1.ListRepositoryCommitsByReferenceRequest]) (*connect.Response[v1alpha1.ListRepositoryCommitsByReferenceResponse], error) {
	return c.listRepositoryCommitsByReference.CallUnary(ctx, req)
}

// GetRepositoryCommitByReference calls
// buf.alpha.registry.v1alpha1.RepositoryCommitService.GetRepositoryCommitByReference.
func (c *repositoryCommitServiceClient) GetRepositoryCommitByReference(ctx context.Context, req *connect.Request[v1alpha1.GetRepositoryCommitByReferenceRequest]) (*connect.Response[v1alpha1.GetRepositoryCommitByReferenceResponse], error) {
	return c.getRepositoryCommitByReference.CallUnary(ctx, req)
}

// ListRepositoryDraftCommits calls
// buf.alpha.registry.v1alpha1.RepositoryCommitService.ListRepositoryDraftCommits.
func (c *repositoryCommitServiceClient) ListRepositoryDraftCommits(ctx context.Context, req *connect.Request[v1alpha1.ListRepositoryDraftCommitsRequest]) (*connect.Response[v1alpha1.ListRepositoryDraftCommitsResponse], error) {
	return c.listRepositoryDraftCommits.CallUnary(ctx, req)
}

// DeleteRepositoryDraftCommit calls
// buf.alpha.registry.v1alpha1.RepositoryCommitService.DeleteRepositoryDraftCommit.
func (c *repositoryCommitServiceClient) DeleteRepositoryDraftCommit(ctx context.Context, req *connect.Request[v1alpha1.DeleteRepositoryDraftCommitRequest]) (*connect.Response[v1alpha1.DeleteRepositoryDraftCommitResponse], error) {
	return c.deleteRepositoryDraftCommit.CallUnary(ctx, req)
}

// RepositoryCommitServiceHandler is an implementation of the
// buf.alpha.registry.v1alpha1.RepositoryCommitService service.
type RepositoryCommitServiceHandler interface {
	// ListRepositoryCommitsByBranch lists the repository commits associated
	// with a repository branch on a repository, ordered by their create time.
	//
	// Deprecated: do not use.
	ListRepositoryCommitsByBranch(context.Context, *connect.Request[v1alpha1.ListRepositoryCommitsByBranchRequest]) (*connect.Response[v1alpha1.ListRepositoryCommitsByBranchResponse], error)
	// ListRepositoryCommitsByReference returns repository commits up-to and including
	// the provided reference.
	ListRepositoryCommitsByReference(context.Context, *connect.Request[v1alpha1.ListRepositoryCommitsByReferenceRequest]) (*connect.Response[v1alpha1.ListRepositoryCommitsByReferenceResponse], error)
	// GetRepositoryCommitByReference returns the repository commit matching
	// the provided reference, if it exists.
	GetRepositoryCommitByReference(context.Context, *connect.Request[v1alpha1.GetRepositoryCommitByReferenceRequest]) (*connect.Response[v1alpha1.GetRepositoryCommitByReferenceResponse], error)
	// ListRepositoryDraftCommits lists draft commits in a repository.
	ListRepositoryDraftCommits(context.Context, *connect.Request[v1alpha1.ListRepositoryDraftCommitsRequest]) (*connect.Response[v1alpha1.ListRepositoryDraftCommitsResponse], error)
	// DeleteRepositoryDraftCommit deletes a draft.
	DeleteRepositoryDraftCommit(context.Context, *connect.Request[v1alpha1.DeleteRepositoryDraftCommitRequest]) (*connect.Response[v1alpha1.DeleteRepositoryDraftCommitResponse], error)
}

// NewRepositoryCommitServiceHandler builds an HTTP handler from the service implementation. It
// returns the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewRepositoryCommitServiceHandler(svc RepositoryCommitServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	repositoryCommitServiceListRepositoryCommitsByBranchHandler := connect.NewUnaryHandler(
		RepositoryCommitServiceListRepositoryCommitsByBranchProcedure,
		svc.ListRepositoryCommitsByBranch,
		connect.WithSchema(repositoryCommitServiceListRepositoryCommitsByBranchMethodDescriptor),
		connect.WithIdempotency(connect.IdempotencyNoSideEffects),
		connect.WithHandlerOptions(opts...),
	)
	repositoryCommitServiceListRepositoryCommitsByReferenceHandler := connect.NewUnaryHandler(
		RepositoryCommitServiceListRepositoryCommitsByReferenceProcedure,
		svc.ListRepositoryCommitsByReference,
		connect.WithSchema(repositoryCommitServiceListRepositoryCommitsByReferenceMethodDescriptor),
		connect.WithIdempotency(connect.IdempotencyNoSideEffects),
		connect.WithHandlerOptions(opts...),
	)
	repositoryCommitServiceGetRepositoryCommitByReferenceHandler := connect.NewUnaryHandler(
		RepositoryCommitServiceGetRepositoryCommitByReferenceProcedure,
		svc.GetRepositoryCommitByReference,
		connect.WithSchema(repositoryCommitServiceGetRepositoryCommitByReferenceMethodDescriptor),
		connect.WithIdempotency(connect.IdempotencyNoSideEffects),
		connect.WithHandlerOptions(opts...),
	)
	repositoryCommitServiceListRepositoryDraftCommitsHandler := connect.NewUnaryHandler(
		RepositoryCommitServiceListRepositoryDraftCommitsProcedure,
		svc.ListRepositoryDraftCommits,
		connect.WithSchema(repositoryCommitServiceListRepositoryDraftCommitsMethodDescriptor),
		connect.WithIdempotency(connect.IdempotencyNoSideEffects),
		connect.WithHandlerOptions(opts...),
	)
	repositoryCommitServiceDeleteRepositoryDraftCommitHandler := connect.NewUnaryHandler(
		RepositoryCommitServiceDeleteRepositoryDraftCommitProcedure,
		svc.DeleteRepositoryDraftCommit,
		connect.WithSchema(repositoryCommitServiceDeleteRepositoryDraftCommitMethodDescriptor),
		connect.WithIdempotency(connect.IdempotencyIdempotent),
		connect.WithHandlerOptions(opts...),
	)
	return "/buf.alpha.registry.v1alpha1.RepositoryCommitService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case RepositoryCommitServiceListRepositoryCommitsByBranchProcedure:
			repositoryCommitServiceListRepositoryCommitsByBranchHandler.ServeHTTP(w, r)
		case RepositoryCommitServiceListRepositoryCommitsByReferenceProcedure:
			repositoryCommitServiceListRepositoryCommitsByReferenceHandler.ServeHTTP(w, r)
		case RepositoryCommitServiceGetRepositoryCommitByReferenceProcedure:
			repositoryCommitServiceGetRepositoryCommitByReferenceHandler.ServeHTTP(w, r)
		case RepositoryCommitServiceListRepositoryDraftCommitsProcedure:
			repositoryCommitServiceListRepositoryDraftCommitsHandler.ServeHTTP(w, r)
		case RepositoryCommitServiceDeleteRepositoryDraftCommitProcedure:
			repositoryCommitServiceDeleteRepositoryDraftCommitHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedRepositoryCommitServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedRepositoryCommitServiceHandler struct{}

func (UnimplementedRepositoryCommitServiceHandler) ListRepositoryCommitsByBranch(context.Context, *connect.Request[v1alpha1.ListRepositoryCommitsByBranchRequest]) (*connect.Response[v1alpha1.ListRepositoryCommitsByBranchResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.RepositoryCommitService.ListRepositoryCommitsByBranch is not implemented"))
}

func (UnimplementedRepositoryCommitServiceHandler) ListRepositoryCommitsByReference(context.Context, *connect.Request[v1alpha1.ListRepositoryCommitsByReferenceRequest]) (*connect.Response[v1alpha1.ListRepositoryCommitsByReferenceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.RepositoryCommitService.ListRepositoryCommitsByReference is not implemented"))
}

func (UnimplementedRepositoryCommitServiceHandler) GetRepositoryCommitByReference(context.Context, *connect.Request[v1alpha1.GetRepositoryCommitByReferenceRequest]) (*connect.Response[v1alpha1.GetRepositoryCommitByReferenceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.RepositoryCommitService.GetRepositoryCommitByReference is not implemented"))
}

func (UnimplementedRepositoryCommitServiceHandler) ListRepositoryDraftCommits(context.Context, *connect.Request[v1alpha1.ListRepositoryDraftCommitsRequest]) (*connect.Response[v1alpha1.ListRepositoryDraftCommitsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.RepositoryCommitService.ListRepositoryDraftCommits is not implemented"))
}

func (UnimplementedRepositoryCommitServiceHandler) DeleteRepositoryDraftCommit(context.Context, *connect.Request[v1alpha1.DeleteRepositoryDraftCommitRequest]) (*connect.Response[v1alpha1.DeleteRepositoryDraftCommitResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.RepositoryCommitService.DeleteRepositoryDraftCommit is not implemented"))
}
