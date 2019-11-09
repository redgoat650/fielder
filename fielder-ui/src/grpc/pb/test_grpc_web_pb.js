/**
 * @fileoverview gRPC-Web generated client stub for testgrpc
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!
/* eslint-disable */


const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.testgrpc = require('./test_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.testgrpc.TestServiceClient =
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

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.testgrpc.TestServicePromiseClient =
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

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.testgrpc.Ping,
 *   !proto.testgrpc.Pong>}
 */
const methodDescriptor_TestService_TestRPC = new grpc.web.MethodDescriptor(
  '/testgrpc.TestService/TestRPC',
  grpc.web.MethodType.UNARY,
  proto.testgrpc.Ping,
  proto.testgrpc.Pong,
  /** @param {!proto.testgrpc.Ping} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.testgrpc.Pong.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.testgrpc.Ping,
 *   !proto.testgrpc.Pong>}
 */
const methodInfo_TestService_TestRPC = new grpc.web.AbstractClientBase.MethodInfo(
  proto.testgrpc.Pong,
  /** @param {!proto.testgrpc.Ping} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.testgrpc.Pong.deserializeBinary
);


/**
 * @param {!proto.testgrpc.Ping} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.testgrpc.Pong)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.testgrpc.Pong>|undefined}
 *     The XHR Node Readable Stream
 */
proto.testgrpc.TestServiceClient.prototype.testRPC =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/testgrpc.TestService/TestRPC',
      request,
      metadata || {},
      methodDescriptor_TestService_TestRPC,
      callback);
};


/**
 * @param {!proto.testgrpc.Ping} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.testgrpc.Pong>}
 *     A native promise that resolves to the response
 */
proto.testgrpc.TestServicePromiseClient.prototype.testRPC =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/testgrpc.TestService/TestRPC',
      request,
      metadata || {},
      methodDescriptor_TestService_TestRPC);
};


module.exports = proto.testgrpc;

