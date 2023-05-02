package grpclog

import (
	"context"
	"fmt"
	"logging-service/api/logging"
	"logging-service/cmd/model"
)

type LogServer struct {
	logging.UnimplementedLogServer
	db *model.Model
}

func NewLogServer(db *model.Model) *LogServer {
	return &LogServer{db: db}
}

func (s *LogServer) LogViaGRPC(context context.Context, req *logging.LogRequest) (*logging.LogResponse, error) {
	id, err := s.db.InsertOne(req.Name, req.Data)
	if err != nil {
		return &logging.LogResponse{Message: "logging via grpc failed"}, err
	}
	return &logging.LogResponse{Message: fmt.Sprintf("logging via grpc succeed and the log id is %s", id)}, nil
}
