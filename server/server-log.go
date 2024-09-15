package server

import (
	"github.com/RodrigoCF25/0231637_SistemasDistribuidos/log"

	api "github.com/RodrigoCF25/0231637_SistemasDistribuidos/api/v1"
)

var _ api.LogServer = (*grpcServer)(nil)

type grpcServer struct {
	api.UnimplementedLogServer
	*log.Log
}

func newgrpcServer(commitlog *log.Log) (srv *grpcServer, err error) {
	srv = &grpcServer{
		Log: commitlog,
	}
	return srv, nil
}
