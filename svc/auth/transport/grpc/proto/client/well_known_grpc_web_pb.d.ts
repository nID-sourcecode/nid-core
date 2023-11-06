import * as grpcWeb from 'grpc-web';

import * as google_protobuf_descriptor_pb from 'google-protobuf/google/protobuf/descriptor_pb';

import {
  WellKnownRequest,
  WellKnownResponse} from './well_known_pb';

export class WellKnownClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  getWellKnownOpenIDConfiguration(
    request: WellKnownRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: WellKnownResponse) => void
  ): grpcWeb.ClientReadableStream<WellKnownResponse>;

  getWellKnownOAuthAuthorizationServer(
    request: WellKnownRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: WellKnownResponse) => void
  ): grpcWeb.ClientReadableStream<WellKnownResponse>;

}

export class WellKnownPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  getWellKnownOpenIDConfiguration(
    request: WellKnownRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<WellKnownResponse>;

  getWellKnownOAuthAuthorizationServer(
    request: WellKnownRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<WellKnownResponse>;

}

