import * as jspb from "google-protobuf"

import * as google_protobuf_descriptor_pb from 'google-protobuf/google/protobuf/descriptor_pb';
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';
import * as google_rpc_code_pb from './google/rpc/code_pb';
import * as google_rpc_error_details_pb from './google/rpc/error_details_pb';
import * as google_rpc_status_pb from './google/rpc/status_pb';
import * as well_known_pb from './well_known_pb';

export class AuthorizeRequest extends jspb.Message {
  getScope(): string;
  setScope(value: string): void;

  getResponseType(): string;
  setResponseType(value: string): void;

  getClientId(): string;
  setClientId(value: string): void;

  getRedirectUri(): string;
  setRedirectUri(value: string): void;

  getAudience(): string;
  setAudience(value: string): void;

  getOptionalScopes(): string;
  setOptionalScopes(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthorizeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AuthorizeRequest): AuthorizeRequest.AsObject;
  static serializeBinaryToWriter(message: AuthorizeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthorizeRequest;
  static deserializeBinaryFromReader(message: AuthorizeRequest, reader: jspb.BinaryReader): AuthorizeRequest;
}

export namespace AuthorizeRequest {
  export type AsObject = {
    scope: string,
    responseType: string,
    clientId: string,
    redirectUri: string,
    audience: string,
    optionalScopes: string,
  }
}

export class AuthorizeHeadlessRequest extends jspb.Message {
  getResponseType(): string;
  setResponseType(value: string): void;

  getClientId(): string;
  setClientId(value: string): void;

  getRedirectUri(): string;
  setRedirectUri(value: string): void;

  getAudience(): string;
  setAudience(value: string): void;

  getQueryModelJson(): string;
  setQueryModelJson(value: string): void;

  getQueryModelPath(): string;
  setQueryModelPath(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthorizeHeadlessRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AuthorizeHeadlessRequest): AuthorizeHeadlessRequest.AsObject;
  static serializeBinaryToWriter(message: AuthorizeHeadlessRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthorizeHeadlessRequest;
  static deserializeBinaryFromReader(message: AuthorizeHeadlessRequest, reader: jspb.BinaryReader): AuthorizeHeadlessRequest;
}

export namespace AuthorizeHeadlessRequest {
  export type AsObject = {
    responseType: string,
    clientId: string,
    redirectUri: string,
    audience: string,
    queryModelJson: string,
    queryModelPath: string,
  }
}

export class SessionRequest extends jspb.Message {
  getSessionId(): string;
  setSessionId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SessionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SessionRequest): SessionRequest.AsObject;
  static serializeBinaryToWriter(message: SessionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SessionRequest;
  static deserializeBinaryFromReader(message: SessionRequest, reader: jspb.BinaryReader): SessionRequest;
}

export namespace SessionRequest {
  export type AsObject = {
    sessionId: string,
  }
}

export class SessionResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getState(): SessionState;
  setState(value: SessionState): void;

  getClient(): Client | undefined;
  setClient(value?: Client): void;
  hasClient(): boolean;
  clearClient(): void;

  getAudience(): Audience | undefined;
  setAudience(value?: Audience): void;
  hasAudience(): boolean;
  clearAudience(): void;

  getRequiredAccessModelsList(): Array<AccessModel>;
  setRequiredAccessModelsList(value: Array<AccessModel>): void;
  clearRequiredAccessModelsList(): void;
  addRequiredAccessModels(value?: AccessModel, index?: number): AccessModel;

  getOptionalAccessModelsList(): Array<AccessModel>;
  setOptionalAccessModelsList(value: Array<AccessModel>): void;
  clearOptionalAccessModelsList(): void;
  addOptionalAccessModels(value?: AccessModel, index?: number): AccessModel;

  getAcceptedAccessModelsList(): Array<AccessModel>;
  setAcceptedAccessModelsList(value: Array<AccessModel>): void;
  clearAcceptedAccessModelsList(): void;
  addAcceptedAccessModels(value?: AccessModel, index?: number): AccessModel;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SessionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SessionResponse): SessionResponse.AsObject;
  static serializeBinaryToWriter(message: SessionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SessionResponse;
  static deserializeBinaryFromReader(message: SessionResponse, reader: jspb.BinaryReader): SessionResponse;
}

export namespace SessionResponse {
  export type AsObject = {
    id: string,
    state: SessionState,
    client?: Client.AsObject,
    audience?: Audience.AsObject,
    requiredAccessModelsList: Array<AccessModel.AsObject>,
    optionalAccessModelsList: Array<AccessModel.AsObject>,
    acceptedAccessModelsList: Array<AccessModel.AsObject>,
  }
}

export class AcceptRequest extends jspb.Message {
  getSessionId(): string;
  setSessionId(value: string): void;

  getAccessModelIdsList(): Array<string>;
  setAccessModelIdsList(value: Array<string>): void;
  clearAccessModelIdsList(): void;
  addAccessModelIds(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AcceptRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AcceptRequest): AcceptRequest.AsObject;
  static serializeBinaryToWriter(message: AcceptRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AcceptRequest;
  static deserializeBinaryFromReader(message: AcceptRequest, reader: jspb.BinaryReader): AcceptRequest;
}

export namespace AcceptRequest {
  export type AsObject = {
    sessionId: string,
    accessModelIdsList: Array<string>,
  }
}

export class StatusResponse extends jspb.Message {
  getState(): SessionState;
  setState(value: SessionState): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StatusResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StatusResponse): StatusResponse.AsObject;
  static serializeBinaryToWriter(message: StatusResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StatusResponse;
  static deserializeBinaryFromReader(message: StatusResponse, reader: jspb.BinaryReader): StatusResponse;
}

export namespace StatusResponse {
  export type AsObject = {
    state: SessionState,
  }
}

export class AccessModel extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getHash(): string;
  setHash(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccessModel.AsObject;
  static toObject(includeInstance: boolean, msg: AccessModel): AccessModel.AsObject;
  static serializeBinaryToWriter(message: AccessModel, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AccessModel;
  static deserializeBinaryFromReader(message: AccessModel, reader: jspb.BinaryReader): AccessModel;
}

export namespace AccessModel {
  export type AsObject = {
    id: string,
    name: string,
    hash: string,
    description: string,
  }
}

export class Client extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getLogo(): string;
  setLogo(value: string): void;

  getIcon(): string;
  setIcon(value: string): void;

  getColor(): string;
  setColor(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Client.AsObject;
  static toObject(includeInstance: boolean, msg: Client): Client.AsObject;
  static serializeBinaryToWriter(message: Client, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Client;
  static deserializeBinaryFromReader(message: Client, reader: jspb.BinaryReader): Client;
}

export namespace Client {
  export type AsObject = {
    id: string,
    name: string,
    logo: string,
    icon: string,
    color: string,
  }
}

export class Audience extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getAudience(): string;
  setAudience(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Audience.AsObject;
  static toObject(includeInstance: boolean, msg: Audience): Audience.AsObject;
  static serializeBinaryToWriter(message: Audience, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Audience;
  static deserializeBinaryFromReader(message: Audience, reader: jspb.BinaryReader): Audience;
}

export namespace Audience {
  export type AsObject = {
    id: string,
    audience: string,
    namespace: string,
  }
}

export class TokenClientFlowRequest extends jspb.Message {
  getGrantType(): string;
  setGrantType(value: string): void;

  getScope(): string;
  setScope(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TokenClientFlowRequest.AsObject;
  static toObject(includeInstance: boolean, msg: TokenClientFlowRequest): TokenClientFlowRequest.AsObject;
  static serializeBinaryToWriter(message: TokenClientFlowRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TokenClientFlowRequest;
  static deserializeBinaryFromReader(message: TokenClientFlowRequest, reader: jspb.BinaryReader): TokenClientFlowRequest;
}

export namespace TokenClientFlowRequest {
  export type AsObject = {
    grantType: string,
    scope: string,
  }
}

export class TokenRequest extends jspb.Message {
  getGrantType(): string;
  setGrantType(value: string): void;

  getAuthorizationCode(): string;
  setAuthorizationCode(value: string): void;

  getRefreshToken(): string;
  setRefreshToken(value: string): void;

  getTypeValueCase(): TokenRequest.TypeValueCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TokenRequest.AsObject;
  static toObject(includeInstance: boolean, msg: TokenRequest): TokenRequest.AsObject;
  static serializeBinaryToWriter(message: TokenRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TokenRequest;
  static deserializeBinaryFromReader(message: TokenRequest, reader: jspb.BinaryReader): TokenRequest;
}

export namespace TokenRequest {
  export type AsObject = {
    grantType: string,
    authorizationCode: string,
    refreshToken: string,
  }

  export enum TypeValueCase { 
    TYPE_VALUE_NOT_SET = 0,
    AUTHORIZATION_CODE = 2,
    REFRESH_TOKEN = 3,
  }
}

export class SessionAuthorization extends jspb.Message {
  getFinaliseToken(): string;
  setFinaliseToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SessionAuthorization.AsObject;
  static toObject(includeInstance: boolean, msg: SessionAuthorization): SessionAuthorization.AsObject;
  static serializeBinaryToWriter(message: SessionAuthorization, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SessionAuthorization;
  static deserializeBinaryFromReader(message: SessionAuthorization, reader: jspb.BinaryReader): SessionAuthorization;
}

export namespace SessionAuthorization {
  export type AsObject = {
    finaliseToken: string,
  }
}

export class FinaliseRequest extends jspb.Message {
  getSessionId(): string;
  setSessionId(value: string): void;

  getSessionFinaliseToken(): string;
  setSessionFinaliseToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FinaliseRequest.AsObject;
  static toObject(includeInstance: boolean, msg: FinaliseRequest): FinaliseRequest.AsObject;
  static serializeBinaryToWriter(message: FinaliseRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FinaliseRequest;
  static deserializeBinaryFromReader(message: FinaliseRequest, reader: jspb.BinaryReader): FinaliseRequest;
}

export namespace FinaliseRequest {
  export type AsObject = {
    sessionId: string,
    sessionFinaliseToken: string,
  }
}

export class FinaliseResponse extends jspb.Message {
  getRedirectLocation(): string;
  setRedirectLocation(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FinaliseResponse.AsObject;
  static toObject(includeInstance: boolean, msg: FinaliseResponse): FinaliseResponse.AsObject;
  static serializeBinaryToWriter(message: FinaliseResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FinaliseResponse;
  static deserializeBinaryFromReader(message: FinaliseResponse, reader: jspb.BinaryReader): FinaliseResponse;
}

export namespace FinaliseResponse {
  export type AsObject = {
    redirectLocation: string,
  }
}

export class TokenResponse extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): void;

  getRefreshToken(): string;
  setRefreshToken(value: string): void;

  getTokenType(): string;
  setTokenType(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TokenResponse.AsObject;
  static toObject(includeInstance: boolean, msg: TokenResponse): TokenResponse.AsObject;
  static serializeBinaryToWriter(message: TokenResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TokenResponse;
  static deserializeBinaryFromReader(message: TokenResponse, reader: jspb.BinaryReader): TokenResponse;
}

export namespace TokenResponse {
  export type AsObject = {
    accessToken: string,
    refreshToken: string,
    tokenType: string,
  }
}

export class AccessModelRequest extends jspb.Message {
  getAudience(): string;
  setAudience(value: string): void;

  getQueryModelJson(): string;
  setQueryModelJson(value: string): void;

  getScopeName(): string;
  setScopeName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccessModelRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AccessModelRequest): AccessModelRequest.AsObject;
  static serializeBinaryToWriter(message: AccessModelRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AccessModelRequest;
  static deserializeBinaryFromReader(message: AccessModelRequest, reader: jspb.BinaryReader): AccessModelRequest;
}

export namespace AccessModelRequest {
  export type AsObject = {
    audience: string,
    queryModelJson: string,
    scopeName: string,
    description: string,
  }
}

export class SwapTokenRequest extends jspb.Message {
  getCurrentToken(): string;
  setCurrentToken(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getAudience(): string;
  setAudience(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SwapTokenRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SwapTokenRequest): SwapTokenRequest.AsObject;
  static serializeBinaryToWriter(message: SwapTokenRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SwapTokenRequest;
  static deserializeBinaryFromReader(message: SwapTokenRequest, reader: jspb.BinaryReader): SwapTokenRequest;
}

export namespace SwapTokenRequest {
  export type AsObject = {
    currentToken: string,
    query: string,
    audience: string,
  }
}

export enum SessionState { 
  UNSPECIFIED = 0,
  UNCLAIMED = 1,
  CLAIMED = 2,
  ACCEPTED = 3,
  REJECTED = 4,
  CODE_GRANTED = 5,
  TOKEN_GRANTED = 6,
}
