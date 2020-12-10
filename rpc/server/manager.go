package server

import (
	"context"
	"fmt"

	"github.com/zllovesuki/G14Manager/rpc/protocol"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RequestType int

const (
	RequestCheckState RequestType = iota
	RequestStartController
	RequestStopController
)

type SupervisorRequest struct {
	Request  RequestType
	Response chan SupervisorResponse
}

type SupervisorResponse struct {
	State protocol.ManagerControlResponse_CurrentState
	Error error
}

type ManagerServer struct {
	protocol.UnimplementedManagerControlServer

	control chan SupervisorRequest
}

var _ protocol.ManagerControlServer = &ManagerServer{}

func RegisterManagerServer(s *grpc.Server, ctrl chan SupervisorRequest) *ManagerServer {
	server := &ManagerServer{
		control: ctrl,
	}
	protocol.RegisterManagerControlServer(s, server)
	return server
}

func (m *ManagerServer) GetCurrentState(ctx context.Context, req *emptypb.Empty) (*protocol.ManagerControlResponse, error) {
	respChan := make(chan SupervisorResponse)
	supervisorReq := SupervisorRequest{
		Request:  RequestCheckState,
		Response: respChan,
	}
	m.control <- supervisorReq
	resp := <-respChan
	return &protocol.ManagerControlResponse{
		Success: true,
		State:   resp.State,
	}, nil
}

func (m *ManagerServer) Control(ctx context.Context, req *protocol.ManagerControlRequest) (*protocol.ManagerControlResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("nil request is invalid")
	}
	respChan := make(chan SupervisorResponse)
	supervisorReq := SupervisorRequest{
		Response: respChan,
	}
	if req.GetState() == protocol.ManagerControlRequest_START {
		supervisorReq.Request = RequestStartController
	} else {
		supervisorReq.Request = RequestStopController
	}
	m.control <- supervisorReq
	resp := <-respChan
	if resp.Error != nil {
		return &protocol.ManagerControlResponse{
			Success: false,
			State:   resp.State,
			Message: resp.Error.Error(),
		}, nil
	}
	return &protocol.ManagerControlResponse{
		Success: true,
		State:   resp.State,
	}, nil
}
