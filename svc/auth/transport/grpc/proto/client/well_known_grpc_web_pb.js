/**
 * @fileoverview gRPC-Web generated client stub for auth
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');



var google_protobuf_descriptor_pb = require('google-protobuf/google/protobuf/descriptor_pb.js')


const proto = {};
proto.auth = require('./well_known_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.auth.WellKnownClient =
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
proto.auth.WellKnownPromiseClient =
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
 *   !proto.auth.WellKnownRequest,
 *   !proto.auth.WellKnownResponse>}
 */
const methodDescriptor_WellKnown_GetWellKnownOpenIDConfiguration = new grpc.web.MethodDescriptor(
  '/auth.WellKnown/GetWellKnownOpenIDConfiguration',
  grpc.web.MethodType.UNARY,
  proto.auth.WellKnownRequest,
  proto.auth.WellKnownResponse,
  /**
   * @param {!proto.auth.WellKnownRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.WellKnownResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.auth.WellKnownRequest,
 *   !proto.auth.WellKnownResponse>}
 */
const methodInfo_WellKnown_GetWellKnownOpenIDConfiguration = new grpc.web.AbstractClientBase.MethodInfo(
  proto.auth.WellKnownResponse,
  /**
   * @param {!proto.auth.WellKnownRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.WellKnownResponse.deserializeBinary
);


/**
 * @param {!proto.auth.WellKnownRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.auth.WellKnownResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.auth.WellKnownResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.auth.WellKnownClient.prototype.getWellKnownOpenIDConfiguration =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/auth.WellKnown/GetWellKnownOpenIDConfiguration',
      request,
      metadata || {},
      methodDescriptor_WellKnown_GetWellKnownOpenIDConfiguration,
      callback);
};


/**
 * @param {!proto.auth.WellKnownRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.auth.WellKnownResponse>}
 *     A native promise that resolves to the response
 */
proto.auth.WellKnownPromiseClient.prototype.getWellKnownOpenIDConfiguration =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/auth.WellKnown/GetWellKnownOpenIDConfiguration',
      request,
      metadata || {},
      methodDescriptor_WellKnown_GetWellKnownOpenIDConfiguration);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.auth.WellKnownRequest,
 *   !proto.auth.WellKnownResponse>}
 */
const methodDescriptor_WellKnown_GetWellKnownOAuthAuthorizationServer = new grpc.web.MethodDescriptor(
  '/auth.WellKnown/GetWellKnownOAuthAuthorizationServer',
  grpc.web.MethodType.UNARY,
  proto.auth.WellKnownRequest,
  proto.auth.WellKnownResponse,
  /**
   * @param {!proto.auth.WellKnownRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.WellKnownResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.auth.WellKnownRequest,
 *   !proto.auth.WellKnownResponse>}
 */
const methodInfo_WellKnown_GetWellKnownOAuthAuthorizationServer = new grpc.web.AbstractClientBase.MethodInfo(
  proto.auth.WellKnownResponse,
  /**
   * @param {!proto.auth.WellKnownRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.auth.WellKnownResponse.deserializeBinary
);


/**
 * @param {!proto.auth.WellKnownRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.auth.WellKnownResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.auth.WellKnownResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.auth.WellKnownClient.prototype.getWellKnownOAuthAuthorizationServer =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/auth.WellKnown/GetWellKnownOAuthAuthorizationServer',
      request,
      metadata || {},
      methodDescriptor_WellKnown_GetWellKnownOAuthAuthorizationServer,
      callback);
};


/**
 * @param {!proto.auth.WellKnownRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.auth.WellKnownResponse>}
 *     A native promise that resolves to the response
 */
proto.auth.WellKnownPromiseClient.prototype.getWellKnownOAuthAuthorizationServer =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/auth.WellKnown/GetWellKnownOAuthAuthorizationServer',
      request,
      metadata || {},
      methodDescriptor_WellKnown_GetWellKnownOAuthAuthorizationServer);
};


module.exports = proto.auth;

