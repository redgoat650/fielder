const {Ping, Pong} = require('../pb/test_pb.js');
const {TestServiceClient} = require('./test_grpc_web_pb.js');

var client = new TestServiceClient('http://localhost:8080');

var request = new Ping();
request.setMessage('Fuck me?');

client.testRPC(request, {}, (err, response) => {
  console.log(response.getMessage());
});