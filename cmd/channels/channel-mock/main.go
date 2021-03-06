package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pjoc-team/pay-gateway/internal/service"
	"github.com/pjoc-team/pay-gateway/pkg/channels/mock"
	"github.com/pjoc-team/pay-gateway/pkg/config"
	_ "github.com/pjoc-team/pay-gateway/pkg/config/file"
	"github.com/pjoc-team/pay-gateway/pkg/discovery"
	pay "github.com/pjoc-team/pay-proto/go"
	"github.com/pjoc-team/tracing/logger"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
)

const serviceName = discovery.ServiceName("channel-mock")

var (
	configURL string
)

func flagSet() *pflag.FlagSet {
	set := pflag.NewFlagSet(serviceName.String(), pflag.ExitOnError)
	set.StringVarP(
		&configURL, "config-url", "c",
		"file://./conf/biz/channel/mock.yaml", "config url",
	)
	return set
}

func main() {
	log := logger.Log()

	s, err := service.NewServer(serviceName.String())
	if err != nil {
		log.Fatal(err.Error())
	}

	if err != nil {
		log.Fatal(err.Error())
	}

	err = s.Init(service.WithFlagSet(flagSet()))
	if err != nil {
		log.Fatal(err.Error())
	}
	if configURL == "" {
		log.Fatal("config url is nill")
	}

	cs, err := config.InitConfigServer(configURL)
	if err != nil {
		log.Fatalf("illegal configs, error: %v", err.Error())
	}

	server, err := mock.NewServer(cs)
	if err != nil {
		log.Fatalf("failed init server, error: %v", err.Error())
	}
	grpcInfo := &service.GrpcInfo{
		RegisterGrpcFunc: func(ctx context.Context, gs *grpc.Server) error {
			pay.RegisterPayChannelServer(gs, server)
			return nil
		},
		RegisterGatewayFunc: func(ctx context.Context, mux *runtime.ServeMux) error {
			err := pay.RegisterPayChannelHandlerServer(ctx, mux, server)
			return err
		},
		Name: serviceName.String(),
	}
	s.Start(service.WithGrpc(grpcInfo))
}
