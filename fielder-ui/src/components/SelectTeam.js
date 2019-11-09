import React, { Component } from "react";
import { connect } from "react-redux";
import { Button } from "react-bootstrap";

export class SelectTeam extends Component {
    handleAddTeam() {
        console.log("ADD TEAM!");
        
        const {Ping, Pong} = require('../grpc/pb/test_pb.js');
        const {TestServiceClient} = require('../grpc/pb/test_grpc_web_pb.js');

        var client = new TestServiceClient('http://localhost:8080');

        var request = new Ping();
        request.setMessage('Fuck me?');

        client.testRPC(request, {}, (err, response) => {
            console.log(response.getMessage());
        });
    }
    render() {
        return (
            <div>
                <Button onClick={this.handleAddTeam}>Add team</Button>
            </div>
        );
    }
}

const mapStateToProps = state => ({});

const mapDispatchToProps = {};

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(SelectTeam);
