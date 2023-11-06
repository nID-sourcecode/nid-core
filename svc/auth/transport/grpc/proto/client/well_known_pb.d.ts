import * as jspb from "google-protobuf"

import * as google_protobuf_descriptor_pb from 'google-protobuf/google/protobuf/descriptor_pb';

export class WellKnownRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WellKnownRequest.AsObject;
  static toObject(includeInstance: boolean, msg: WellKnownRequest): WellKnownRequest.AsObject;
  static serializeBinaryToWriter(message: WellKnownRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WellKnownRequest;
  static deserializeBinaryFromReader(message: WellKnownRequest, reader: jspb.BinaryReader): WellKnownRequest;
}

export namespace WellKnownRequest {
  export type AsObject = {
  }
}

export class WellKnownResponse extends jspb.Message {
  getIssuer(): string;
  setIssuer(value: string): void;

  getAuthorizationEndpoint(): string;
  setAuthorizationEndpoint(value: string): void;

  getTokenEndpoint(): string;
  setTokenEndpoint(value: string): void;

  getJwksUri(): string;
  setJwksUri(value: string): void;

  getRegistrationEndpoint(): string;
  setRegistrationEndpoint(value: string): void;

  getScopesSupportedList(): Array<string>;
  setScopesSupportedList(value: Array<string>): void;
  clearScopesSupportedList(): void;
  addScopesSupported(value: string, index?: number): void;

  getResponseTypesSupportedList(): Array<string>;
  setResponseTypesSupportedList(value: Array<string>): void;
  clearResponseTypesSupportedList(): void;
  addResponseTypesSupported(value: string, index?: number): void;

  getResponseModesSupportedList(): Array<string>;
  setResponseModesSupportedList(value: Array<string>): void;
  clearResponseModesSupportedList(): void;
  addResponseModesSupported(value: string, index?: number): void;

  getGrantTypesSupportedList(): Array<string>;
  setGrantTypesSupportedList(value: Array<string>): void;
  clearGrantTypesSupportedList(): void;
  addGrantTypesSupported(value: string, index?: number): void;

  getTokenEndpointAuthMethodsSupportedList(): Array<string>;
  setTokenEndpointAuthMethodsSupportedList(value: Array<string>): void;
  clearTokenEndpointAuthMethodsSupportedList(): void;
  addTokenEndpointAuthMethodsSupported(value: string, index?: number): void;

  getTokenEndpointAuthSigningAlgValuesSupportedList(): Array<string>;
  setTokenEndpointAuthSigningAlgValuesSupportedList(value: Array<string>): void;
  clearTokenEndpointAuthSigningAlgValuesSupportedList(): void;
  addTokenEndpointAuthSigningAlgValuesSupported(value: string, index?: number): void;

  getServiceDocumentation(): string;
  setServiceDocumentation(value: string): void;

  getUiLocalesSupportedList(): Array<string>;
  setUiLocalesSupportedList(value: Array<string>): void;
  clearUiLocalesSupportedList(): void;
  addUiLocalesSupported(value: string, index?: number): void;

  getOpPolicyUri(): string;
  setOpPolicyUri(value: string): void;

  getOpTosUri(): string;
  setOpTosUri(value: string): void;

  getRevocationEndpoint(): string;
  setRevocationEndpoint(value: string): void;

  getRevocationEndpointAuthMethodsSupportedList(): Array<string>;
  setRevocationEndpointAuthMethodsSupportedList(value: Array<string>): void;
  clearRevocationEndpointAuthMethodsSupportedList(): void;
  addRevocationEndpointAuthMethodsSupported(value: string, index?: number): void;

  getRevocationEndpointAuthSigningAlgValuesSupportedList(): Array<string>;
  setRevocationEndpointAuthSigningAlgValuesSupportedList(value: Array<string>): void;
  clearRevocationEndpointAuthSigningAlgValuesSupportedList(): void;
  addRevocationEndpointAuthSigningAlgValuesSupported(value: string, index?: number): void;

  getIntrospectionEndpoint(): string;
  setIntrospectionEndpoint(value: string): void;

  getIntrospectionEndpointAuthMethodsSupportedList(): Array<string>;
  setIntrospectionEndpointAuthMethodsSupportedList(value: Array<string>): void;
  clearIntrospectionEndpointAuthMethodsSupportedList(): void;
  addIntrospectionEndpointAuthMethodsSupported(value: string, index?: number): void;

  getIntrospectionEndpointAuthSigningAlgValuesSupportedList(): Array<string>;
  setIntrospectionEndpointAuthSigningAlgValuesSupportedList(value: Array<string>): void;
  clearIntrospectionEndpointAuthSigningAlgValuesSupportedList(): void;
  addIntrospectionEndpointAuthSigningAlgValuesSupported(value: string, index?: number): void;

  getCodeChallengeMethodsSupportedList(): Array<string>;
  setCodeChallengeMethodsSupportedList(value: Array<string>): void;
  clearCodeChallengeMethodsSupportedList(): void;
  addCodeChallengeMethodsSupported(value: string, index?: number): void;

  getSignedMetadata(): string;
  setSignedMetadata(value: string): void;

  getUserinfoEndpoint(): string;
  setUserinfoEndpoint(value: string): void;

  getAcrValuesSupportedList(): Array<string>;
  setAcrValuesSupportedList(value: Array<string>): void;
  clearAcrValuesSupportedList(): void;
  addAcrValuesSupported(value: string, index?: number): void;

  getSubjectTypesSupportedList(): Array<string>;
  setSubjectTypesSupportedList(value: Array<string>): void;
  clearSubjectTypesSupportedList(): void;
  addSubjectTypesSupported(value: string, index?: number): void;

  getIdTokenSigningAlgValuesSupportedList(): Array<string>;
  setIdTokenSigningAlgValuesSupportedList(value: Array<string>): void;
  clearIdTokenSigningAlgValuesSupportedList(): void;
  addIdTokenSigningAlgValuesSupported(value: string, index?: number): void;

  getIdTokenEncryptionAlgValuesSupportedList(): Array<string>;
  setIdTokenEncryptionAlgValuesSupportedList(value: Array<string>): void;
  clearIdTokenEncryptionAlgValuesSupportedList(): void;
  addIdTokenEncryptionAlgValuesSupported(value: string, index?: number): void;

  getIdTokenEncryptionEncValuesSupportedList(): Array<string>;
  setIdTokenEncryptionEncValuesSupportedList(value: Array<string>): void;
  clearIdTokenEncryptionEncValuesSupportedList(): void;
  addIdTokenEncryptionEncValuesSupported(value: string, index?: number): void;

  getUserinfoSigningAlgValuesSupportedList(): Array<string>;
  setUserinfoSigningAlgValuesSupportedList(value: Array<string>): void;
  clearUserinfoSigningAlgValuesSupportedList(): void;
  addUserinfoSigningAlgValuesSupported(value: string, index?: number): void;

  getUserinfoEncryptionAlgValuesSupportedList(): Array<string>;
  setUserinfoEncryptionAlgValuesSupportedList(value: Array<string>): void;
  clearUserinfoEncryptionAlgValuesSupportedList(): void;
  addUserinfoEncryptionAlgValuesSupported(value: string, index?: number): void;

  getUserinfoEncryptionEncValuesSupportedList(): Array<string>;
  setUserinfoEncryptionEncValuesSupportedList(value: Array<string>): void;
  clearUserinfoEncryptionEncValuesSupportedList(): void;
  addUserinfoEncryptionEncValuesSupported(value: string, index?: number): void;

  getRequestObjectSigningAlgValuesSupportedList(): Array<string>;
  setRequestObjectSigningAlgValuesSupportedList(value: Array<string>): void;
  clearRequestObjectSigningAlgValuesSupportedList(): void;
  addRequestObjectSigningAlgValuesSupported(value: string, index?: number): void;

  getRequestObjectEncryptionAlgValuesSupportedList(): Array<string>;
  setRequestObjectEncryptionAlgValuesSupportedList(value: Array<string>): void;
  clearRequestObjectEncryptionAlgValuesSupportedList(): void;
  addRequestObjectEncryptionAlgValuesSupported(value: string, index?: number): void;

  getRequestObjectEncryptionEncValuesSupportedList(): Array<string>;
  setRequestObjectEncryptionEncValuesSupportedList(value: Array<string>): void;
  clearRequestObjectEncryptionEncValuesSupportedList(): void;
  addRequestObjectEncryptionEncValuesSupported(value: string, index?: number): void;

  getDisplayValuesSupportedList(): Array<string>;
  setDisplayValuesSupportedList(value: Array<string>): void;
  clearDisplayValuesSupportedList(): void;
  addDisplayValuesSupported(value: string, index?: number): void;

  getClaimTypesSupportedList(): Array<string>;
  setClaimTypesSupportedList(value: Array<string>): void;
  clearClaimTypesSupportedList(): void;
  addClaimTypesSupported(value: string, index?: number): void;

  getClaimsSupportedList(): Array<string>;
  setClaimsSupportedList(value: Array<string>): void;
  clearClaimsSupportedList(): void;
  addClaimsSupported(value: string, index?: number): void;

  getClaimsLocalesSupportedList(): Array<string>;
  setClaimsLocalesSupportedList(value: Array<string>): void;
  clearClaimsLocalesSupportedList(): void;
  addClaimsLocalesSupported(value: string, index?: number): void;

  getClaimsParameterSupported(): boolean;
  setClaimsParameterSupported(value: boolean): void;

  getRequestParameterSupported(): boolean;
  setRequestParameterSupported(value: boolean): void;

  getRequestUriParameterSupported(): boolean;
  setRequestUriParameterSupported(value: boolean): void;

  getRequireRequestUriRegistration(): boolean;
  setRequireRequestUriRegistration(value: boolean): void;

  getCheckSessionIframe(): string;
  setCheckSessionIframe(value: string): void;

  getFrontchannelLogoutSupported(): boolean;
  setFrontchannelLogoutSupported(value: boolean): void;

  getFrontchannelLogoutSessionSupported(): boolean;
  setFrontchannelLogoutSessionSupported(value: boolean): void;

  getBackchannelLogoutSupported(): boolean;
  setBackchannelLogoutSupported(value: boolean): void;

  getBackchannelLogoutSessionSupported(): boolean;
  setBackchannelLogoutSessionSupported(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WellKnownResponse.AsObject;
  static toObject(includeInstance: boolean, msg: WellKnownResponse): WellKnownResponse.AsObject;
  static serializeBinaryToWriter(message: WellKnownResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WellKnownResponse;
  static deserializeBinaryFromReader(message: WellKnownResponse, reader: jspb.BinaryReader): WellKnownResponse;
}

export namespace WellKnownResponse {
  export type AsObject = {
    issuer: string,
    authorizationEndpoint: string,
    tokenEndpoint: string,
    jwksUri: string,
    registrationEndpoint: string,
    scopesSupportedList: Array<string>,
    responseTypesSupportedList: Array<string>,
    responseModesSupportedList: Array<string>,
    grantTypesSupportedList: Array<string>,
    tokenEndpointAuthMethodsSupportedList: Array<string>,
    tokenEndpointAuthSigningAlgValuesSupportedList: Array<string>,
    serviceDocumentation: string,
    uiLocalesSupportedList: Array<string>,
    opPolicyUri: string,
    opTosUri: string,
    revocationEndpoint: string,
    revocationEndpointAuthMethodsSupportedList: Array<string>,
    revocationEndpointAuthSigningAlgValuesSupportedList: Array<string>,
    introspectionEndpoint: string,
    introspectionEndpointAuthMethodsSupportedList: Array<string>,
    introspectionEndpointAuthSigningAlgValuesSupportedList: Array<string>,
    codeChallengeMethodsSupportedList: Array<string>,
    signedMetadata: string,
    userinfoEndpoint: string,
    acrValuesSupportedList: Array<string>,
    subjectTypesSupportedList: Array<string>,
    idTokenSigningAlgValuesSupportedList: Array<string>,
    idTokenEncryptionAlgValuesSupportedList: Array<string>,
    idTokenEncryptionEncValuesSupportedList: Array<string>,
    userinfoSigningAlgValuesSupportedList: Array<string>,
    userinfoEncryptionAlgValuesSupportedList: Array<string>,
    userinfoEncryptionEncValuesSupportedList: Array<string>,
    requestObjectSigningAlgValuesSupportedList: Array<string>,
    requestObjectEncryptionAlgValuesSupportedList: Array<string>,
    requestObjectEncryptionEncValuesSupportedList: Array<string>,
    displayValuesSupportedList: Array<string>,
    claimTypesSupportedList: Array<string>,
    claimsSupportedList: Array<string>,
    claimsLocalesSupportedList: Array<string>,
    claimsParameterSupported: boolean,
    requestParameterSupported: boolean,
    requestUriParameterSupported: boolean,
    requireRequestUriRegistration: boolean,
    checkSessionIframe: string,
    frontchannelLogoutSupported: boolean,
    frontchannelLogoutSessionSupported: boolean,
    backchannelLogoutSupported: boolean,
    backchannelLogoutSessionSupported: boolean,
  }
}

export enum WellKnownType { 
  WELLKNOWN_TYPE_UNSPECIFIED = 0,
  AUTHORIZATION_ENDPOINT = 1,
  TOKEN_ENDPOINT = 2,
  JWKS_URI = 3,
  REGISTRATION_ENDPOINT = 4,
  SERVICE_DOCUMENTATION = 5,
  OP_POLICY_URI = 6,
  OP_TOS_URI = 7,
  REVOCATION_ENDPOINT = 8,
  INTROSPECTION_ENDPOINT = 9,
  USERINFO_ENDPOINT = 10,
  CHECK_SESSION_IFRAME = 11,
}
