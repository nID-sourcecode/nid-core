/**
 * @fileoverview gRPC-Web generated client stub for auth
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');



var google_protobuf_descriptor_pb = require('google-protobuf/google/protobuf/descriptor_pb.js')

var google_protobuf_empty_pb = require('google-protobuf/google/protobuf/empty_pb.js')

var google_rpc_code_pb = require('./google/rpc/code_pb.js')

var google_rpc_error_details_pb = require('./google/rpc/error_details_pb.js')

var google_rpc_status_pb = require('./google/rpc/status_pb.js')



var well_known_pb = require('./well_known_pb.js')
const proto = {};
proto.auth = require('./auth_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.auth.AuthClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.auth.AuthPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.auth.AuthorizeRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_Auth_Authorize = new grpc.web.MethodDescriptor(
  '/auth.Auth/Authorize',
  grpc.web.MethodType.UNARY,
  proto.auth.AuthorizeRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.auth.AuthorizeRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.auth.AuthorizeRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_Auth_Authorize = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.auth.AuthorizeRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.auth.AuthorizeRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.auth.AuthClient.prototype.authorize =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/auth.Auth/Authorize',
      request,
      metadata || {},
      methodDescriptor_Auth_Authorize,
      callback);
};


/**
 * @param {!proto.auth.AuthorizeRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.auth.AuthPromiseClient.prototype.authorize =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/auth.Auth/Authorize',
      request,
      metadata || {},
      methodDescriptor_Auth_Authorize);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.auth.AuthorizeHeadlessRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_Auth_AuthorizeHeadless = new grpc.web.MethodDescriptor(
  '/auth.Auth/AuthorizeHeadless',
  grpc.web.MethodType.UNARY,
  proto.auth.AuthorizeHeadlessRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.auth.AuthorizeHeadlessRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.auth.AuthorizeHeadlessRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_Auth_AuthorizeHeadless = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.auth.AuthorizeHeadlessRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.auth.AuthorizeHeadlessRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.auth.AuthClient.prototype.authorizeHeadless =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/auth.Auth/AuthorizeHeadless',
      request,
      metadata || {},
      methodDescriptor_Auth_AuthorizeHeadless,
      callback);
};


/**
 * @param {!proto.auth.AuthorizeHeadlessRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.auth.AuthPromiseClient.prototype.authorizeHeadless =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/auth.Auth/AuthorizeHeadless',
      request,
      metadata || {},
      methodDescriptor_Auth_AuthorizeHeadless);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.auth.SessionRequest,
 *   !proto.auth.SessionResponse>}
 */
const methodDescriptor_Auth_Claim = new grpc.web.MethodDescriptor(
  '/auth.Auth/Claim',
  grpc.web.MethodType.UNARY,
  proto.auth.SessionRequest,
  proto.auth.SessionResponse,
  /**
   * @param {!proto.auth.SessionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.SessionResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.auth.SessionRequest,
 *   !proto.auth.SessionResponse>}
 */
const methodInfo_Auth_Claim = new grpc.web.AbstractClientBase.MethodInfo(
  proto.auth.SessionResponse,
  /**
   * @param {!proto.auth.SessionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.SessionResponse.deserializeBinary
);


/**
 * @param {!proto.auth.SessionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.auth.SessionResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.auth.SessionResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.auth.AuthClient.prototype.claim =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/auth.Auth/Claim',
      request,
      metadata || {},
      methodDescriptor_Auth_Claim,
      callback);
};


/**
 * @param {!proto.auth.SessionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.auth.SessionResponse>}
 *     A native promise that resolves to the response
 */
proto.auth.AuthPromiseClient.prototype.claim =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/auth.Auth/Claim',
      request,
      metadata || {},
      methodDescriptor_Auth_Claim);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.auth.AcceptRequest,
 *   !proto.auth.SessionResponse>}
 */
const methodDescriptor_Auth_Accept = new grpc.web.MethodDescriptor(
  '/auth.Auth/Accept',
  grpc.web.MethodType.UNARY,
  proto.auth.AcceptRequest,
  proto.auth.SessionResponse,
  /**
   * @param {!proto.auth.AcceptRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.SessionResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.auth.AcceptRequest,
 *   !proto.auth.SessionResponse>}
 */
const methodInfo_Auth_Accept = new grpc.web.AbstractClientBase.MethodInfo(
  proto.auth.SessionResponse,
  /**
   * @param {!proto.auth.AcceptRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.SessionResponse.deserializeBinary
);


/**
 * @param {!proto.auth.AcceptRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.auth.SessionResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.auth.SessionResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.auth.AuthClient.prototype.accept =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/auth.Auth/Accept',
      request,
      metadata || {},
      methodDescriptor_Auth_Accept,
      callback);
};


/**
 * @param {!proto.auth.AcceptRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.auth.SessionResponse>}
 *     A native promise that resolves to the response
 */
proto.auth.AuthPromiseClient.prototype.accept =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/auth.Auth/Accept',
      request,
      metadata || {},
      methodDescriptor_Auth_Accept);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.auth.SessionRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_Auth_Reject = new grpc.web.MethodDescriptor(
  '/auth.Auth/Reject',
  grpc.web.MethodType.UNARY,
  proto.auth.SessionRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.auth.SessionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.auth.SessionRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_Auth_Reject = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.auth.SessionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.auth.SessionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.auth.AuthClient.prototype.reject =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/auth.Auth/Reject',
      request,
      metadata || {},
      methodDescriptor_Auth_Reject,
      callback);
};


/**
 * @param {!proto.auth.SessionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.auth.AuthPromiseClient.prototype.reject =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/auth.Auth/Reject',
      request,
      metadata || {},
      methodDescriptor_Auth_Reject);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.auth.SessionRequest,
 *   !proto.auth.SessionAuthorization>}
 */
const methodDescriptor_Auth_GenerateSessionFinaliseToken = new grpc.web.MethodDescriptor(
  '/auth.Auth/GenerateSessionFinaliseToken',
  grpc.web.MethodType.UNARY,
  proto.auth.SessionRequest,
  proto.auth.SessionAuthorization,
  /**
   * @param {!proto.auth.SessionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.SessionAuthorization.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.auth.SessionRequest,
 *   !proto.auth.SessionAuthorization>}
 */
const methodInfo_Auth_GenerateSessionFinaliseToken = new grpc.web.AbstractClientBase.MethodInfo(
  proto.auth.SessionAuthorization,
  /**
   * @param {!proto.auth.SessionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.SessionAuthorization.deserializeBinary
);


/**
 * @param {!proto.auth.SessionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.auth.SessionAuthorization)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.auth.SessionAuthorization>|undefined}
 *     The XHR Node Readable Stream
 */
proto.auth.AuthClient.prototype.generateSessionFinaliseToken =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/auth.Auth/GenerateSessionFinaliseToken',
      request,
      metadata || {},
      methodDescriptor_Auth_GenerateSessionFinaliseToken,
      callback);
};


/**
 * @param {!proto.auth.SessionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.auth.SessionAuthorization>}
 *     A native promise that resolves to the response
 */
proto.auth.AuthPromiseClient.prototype.generateSessionFinaliseToken =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/auth.Auth/GenerateSessionFinaliseToken',
      request,
      metadata || {},
      methodDescriptor_Auth_GenerateSessionFinaliseToken);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.auth.SessionRequest,
 *   !proto.auth.SessionResponse>}
 */
const methodDescriptor_Auth_GetSessionDetails = new grpc.web.MethodDescriptor(
  '/auth.Auth/GetSessionDetails',
  grpc.web.MethodType.UNARY,
  proto.auth.SessionRequest,
  proto.auth.SessionResponse,
  /**
   * @param {!proto.auth.SessionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.SessionResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.auth.SessionRequest,
 *   !proto.auth.SessionResponse>}
 */
const methodInfo_Auth_GetSessionDetails = new grpc.web.AbstractClientBase.MethodInfo(
  proto.auth.SessionResponse,
  /**
   * @param {!proto.auth.SessionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.SessionResponse.deserializeBinary
);


/**
 * @param {!proto.auth.SessionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.auth.SessionResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.auth.SessionResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.auth.AuthClient.prototype.getSessionDetails =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/auth.Auth/GetSessionDetails',
      request,
      metadata || {},
      methodDescriptor_Auth_GetSessionDetails,
      callback);
};


/**
 * @param {!proto.auth.SessionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.auth.SessionResponse>}
 *     A native promise that resolves to the response
 */
proto.auth.AuthPromiseClient.prototype.getSessionDetails =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/auth.Auth/GetSessionDetails',
      request,
      metadata || {},
      methodDescriptor_Auth_GetSessionDetails);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.auth.SessionRequest,
 *   !proto.auth.StatusResponse>}
 */
const methodDescriptor_Auth_Status = new grpc.web.MethodDescriptor(
  '/auth.Auth/Status',
  grpc.web.MethodType.UNARY,
  proto.auth.SessionRequest,
  proto.auth.StatusResponse,
  /**
   * @param {!proto.auth.SessionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.StatusResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.auth.SessionRequest,
 *   !proto.auth.StatusResponse>}
 */
const methodInfo_Auth_Status = new grpc.web.AbstractClientBase.MethodInfo(
  proto.auth.StatusResponse,
  /**
   * @param {!proto.auth.SessionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.StatusResponse.deserializeBinary
);


/**
 * @param {!proto.auth.SessionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.auth.StatusResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.auth.StatusResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.auth.AuthClient.prototype.status =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/auth.Auth/Status',
      request,
      metadata || {},
      methodDescriptor_Auth_Status,
      callback);
};


/**
 * @param {!proto.auth.SessionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.auth.StatusResponse>}
 *     A native promise that resolves to the response
 */
proto.auth.AuthPromiseClient.prototype.status =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/auth.Auth/Status',
      request,
      metadata || {},
      methodDescriptor_Auth_Status);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.auth.FinaliseRequest,
 *   !proto.auth.FinaliseResponse>}
 */
const methodDescriptor_Auth_Finalise = new grpc.web.MethodDescriptor(
  '/auth.Auth/Finalise',
  grpc.web.MethodType.UNARY,
  proto.auth.FinaliseRequest,
  proto.auth.FinaliseResponse,
  /**
   * @param {!proto.auth.FinaliseRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.FinaliseResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.auth.FinaliseRequest,
 *   !proto.auth.FinaliseResponse>}
 */
const methodInfo_Auth_Finalise = new grpc.web.AbstractClientBase.MethodInfo(
  proto.auth.FinaliseResponse,
  /**
   * @param {!proto.auth.FinaliseRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.FinaliseResponse.deserializeBinary
);


/**
 * @param {!proto.auth.FinaliseRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.auth.FinaliseResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.auth.FinaliseResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.auth.AuthClient.prototype.finalise =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/auth.Auth/Finalise',
      request,
      metadata || {},
      methodDescriptor_Auth_Finalise,
      callback);
};


/**
 * @param {!proto.auth.FinaliseRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.auth.FinaliseResponse>}
 *     A native promise that resolves to the response
 */
proto.auth.AuthPromiseClient.prototype.finalise =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/auth.Auth/Finalise',
      request,
      metadata || {},
      methodDescriptor_Auth_Finalise);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.auth.TokenRequest,
 *   !proto.auth.TokenResponse>}
 */
const methodDescriptor_Auth_Token = new grpc.web.MethodDescriptor(
  '/auth.Auth/Token',
  grpc.web.MethodType.UNARY,
  proto.auth.TokenRequest,
  proto.auth.TokenResponse,
  /**
   * @param {!proto.auth.TokenRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.TokenResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.auth.TokenRequest,
 *   !proto.auth.TokenResponse>}
 */
const methodInfo_Auth_Token = new grpc.web.AbstractClientBase.MethodInfo(
  proto.auth.TokenResponse,
  /**
   * @param {!proto.auth.TokenRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.TokenResponse.deserializeBinary
);


/**
 * @param {!proto.auth.TokenRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.auth.TokenResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.auth.TokenResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.auth.AuthClient.prototype.token =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/auth.Auth/Token',
      request,
      metadata || {},
      methodDescriptor_Auth_Token,
      callback);
};


/**
 * @param {!proto.auth.TokenRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.auth.TokenResponse>}
 *     A native promise that resolves to the response
 */
proto.auth.AuthPromiseClient.prototype.token =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/auth.Auth/Token',
      request,
      metadata || {},
      methodDescriptor_Auth_Token);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.auth.TokenClientFlowRequest,
 *   !proto.auth.TokenResponse>}
 */
const methodDescriptor_Auth_TokenClientFlow = new grpc.web.MethodDescriptor(
  '/auth.Auth/TokenClientFlow',
  grpc.web.MethodType.UNARY,
  proto.auth.TokenClientFlowRequest,
  proto.auth.TokenResponse,
  /**
   * @param {!proto.auth.TokenClientFlowRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.TokenResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.auth.TokenClientFlowRequest,
 *   !proto.auth.TokenResponse>}
 */
const methodInfo_Auth_TokenClientFlow = new grpc.web.AbstractClientBase.MethodInfo(
  proto.auth.TokenResponse,
  /**
   * @param {!proto.auth.TokenClientFlowRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.TokenResponse.deserializeBinary
);


/**
 * @param {!proto.auth.TokenClientFlowRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.auth.TokenResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.auth.TokenResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.auth.AuthClient.prototype.tokenClientFlow =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/auth.Auth/TokenClientFlow',
      request,
      metadata || {},
      methodDescriptor_Auth_TokenClientFlow,
      callback);
};


/**
 * @param {!proto.auth.TokenClientFlowRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.auth.TokenResponse>}
 *     A native promise that resolves to the response
 */
proto.auth.AuthPromiseClient.prototype.tokenClientFlow =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/auth.Auth/TokenClientFlow',
      request,
      metadata || {},
      methodDescriptor_Auth_TokenClientFlow);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.auth.AccessModelRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_Auth_RegisterAccessModel = new grpc.web.MethodDescriptor(
  '/auth.Auth/RegisterAccessModel',
  grpc.web.MethodType.UNARY,
  proto.auth.AccessModelRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.auth.AccessModelRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.auth.AccessModelRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_Auth_RegisterAccessModel = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.auth.AccessModelRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.auth.AccessModelRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.auth.AuthClient.prototype.registerAccessModel =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/auth.Auth/RegisterAccessModel',
      request,
      metadata || {},
      methodDescriptor_Auth_RegisterAccessModel,
      callback);
};


/**
 * @param {!proto.auth.AccessModelRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.auth.AuthPromiseClient.prototype.registerAccessModel =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/auth.Auth/RegisterAccessModel',
      request,
      metadata || {},
      methodDescriptor_Auth_RegisterAccessModel);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.auth.SwapTokenRequest,
 *   !proto.auth.TokenResponse>}
 */
const methodDescriptor_Auth_SwapToken = new grpc.web.MethodDescriptor(
  '/auth.Auth/SwapToken',
  grpc.web.MethodType.UNARY,
  proto.auth.SwapTokenRequest,
  proto.auth.TokenResponse,
  /**
   * @param {!proto.auth.SwapTokenRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.TokenResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.auth.SwapTokenRequest,
 *   !proto.auth.TokenResponse>}
 */
const methodInfo_Auth_SwapToken = new grpc.web.AbstractClientBase.MethodInfo(
  proto.auth.TokenResponse,
  /**
   * @param {!proto.auth.SwapTokenRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.TokenResponse.deserializeBinary
);


/**
 * @param {!proto.auth.SwapTokenRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.auth.TokenResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.auth.TokenResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.auth.AuthClient.prototype.swapToken =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/auth.Auth/SwapToken',
      request,
      metadata || {},
      methodDescriptor_Auth_SwapToken,
      callback);
};


/**
 * @param {!proto.auth.SwapTokenRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.auth.TokenResponse>}
 *     A native promise that resolves to the response
 */
proto.auth.AuthPromiseClient.prototype.swapToken =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/auth.Auth/SwapToken',
      request,
      metadata || {},
      methodDescriptor_Auth_SwapToken);
};


module.exports = proto.auth;

