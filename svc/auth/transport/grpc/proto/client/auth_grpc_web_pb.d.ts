import * as grpcWeb from 'grpc-web';

import * as google_protobuf_descriptor_pb from 'google-protobuf/google/protobuf/descriptor_pb';
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';
import * as google_rpc_code_pb from './google/rpc/code_pb';
import * as google_rpc_error_details_pb from './google/rpc/error_details_pb';
import * as google_rpc_status_pb from './google/rpc/status_pb';
import * as well_known_pb from './well_known_pb';

import {
  AcceptRequest,
  AccessModelRequest,
  AuthorizeHeadlessRequest,
  AuthorizeRequest,
  FinaliseRequest,
  FinaliseResponse,
  SessionAuthorization,
  SessionRequest,
  SessionResponse,
  StatusResponse,
  SwapTokenRequest,
  TokenClientFlowRequest,
  TokenRequest,
  TokenResponse} from './auth_pb';

export class AuthClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  authorize(
    request: AuthorizeRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  authorizeHeadless(
    request: AuthorizeHeadlessRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  claim(
    request: SessionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SessionResponse) => void
  ): grpcWeb.ClientReadableStream<SessionResponse>;

  accept(
    request: AcceptRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SessionResponse) => void
  ): grpcWeb.ClientReadableStream<SessionResponse>;

  reject(
    request: SessionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  generateSessionFinaliseToken(
    request: SessionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SessionAuthorization) => void
  ): grpcWeb.ClientReadableStream<SessionAuthorization>;

  getSessionDetails(
    request: SessionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SessionResponse) => void
  ): grpcWeb.ClientReadableStream<SessionResponse>;

  status(
    request: SessionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: StatusResponse) => void
  ): grpcWeb.ClientReadableStream<StatusResponse>;

  finalise(
    request: FinaliseRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: FinaliseResponse) => void
  ): grpcWeb.ClientReadableStream<FinaliseResponse>;

  token(
    request: TokenRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: TokenResponse) => void
  ): grpcWeb.ClientReadableStream<TokenResponse>;

  tokenClientFlow(
    request: TokenClientFlowRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: TokenResponse) => void
  ): grpcWeb.ClientReadableStream<TokenResponse>;

  registerAccessModel(
    request: AccessModelRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  swapToken(
    request: SwapTokenRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: TokenResponse) => void
  ): grpcWeb.ClientReadableStream<TokenResponse>;

}

export class AuthPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  authorize(
    request: AuthorizeRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  authorizeHeadless(
    request: AuthorizeHeadlessRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  claim(
    request: SessionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SessionResponse>;

  accept(
    request: AcceptRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SessionResponse>;

  reject(
    request: SessionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  generateSessionFinaliseToken(
    request: SessionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SessionAuthorization>;

  getSessionDetails(
    request: SessionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SessionResponse>;

  status(
    request: SessionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<StatusResponse>;

  finalise(
    request: FinaliseRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<FinaliseResponse>;

  token(
    request: TokenRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<TokenResponse>;

  tokenClientFlow(
    request: TokenClientFlowRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<TokenResponse>;

  registerAccessModel(
    request: AccessModelRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  swapToken(
    request: SwapTokenRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<TokenResponse>;

}

