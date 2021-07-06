package main

import (
	"fmt"
	"github.com/dalconoid/kiddy-lp/api"
	"github.com/dalconoid/kiddy-lp/pkg/grpc_srv"
	"github.com/dalconoid/kiddy-lp/pkg/http_srv"
	"github.com/dalconoid/kiddy-lp/pkg/storage"
	"github.com/dalconoid/kiddy-lp/pkg/works"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type envVars struct {
	HTTPPort          string  `short:"p" long:"http_port" env:"HTTP_PORT" description:"http server port" default:"8080"`
	GRPCPort          string  `short:"g" long:"grpc_port" env:"GRPC_PORT" description:"grpc server port" default:"8081"`
	LinesProviderAddr string  `short:"a" long:"lp_addr" env:"LP_ADDRESS" description:"lines provider address" default:"http://localhost:8000"`
	BTime             float64 `long:"btime" env:"B_TIME" description:"baseball rates refresh time" default:"1"`
	FTime             float64 `long:"ftime" env:"F_TIME" description:"football rates refresh time" default:"1"`
	STime             float64 `long:"stime" env:"S_TIME" description:"soccer rates refresh time" default:"1"`
	StorageAddress    string  `short:"s" long:"storage" env:"STORAGE" description:"storage address" default:"localhost:6379"`
	StoragePass       string  `long:"password" env:"STORAGE_PASSWORD" description:"storage password" default:""`
}

func getEnvVars() (*envVars, error) {
	env := &envVars{}
	if _, err := flags.Parse(env); err != nil {
		return nil, err
	}
	return env, nil
}

func main() {
	log.SetLevel(log.DebugLevel)

	env, err := getEnvVars()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v %v %v %v %v %v %v %v\n", env.BTime, env.FTime, env.STime, env.StorageAddress, env.HTTPPort, env.GRPCPort, env.LinesProviderAddr, env.StoragePass)

	log.Infof("Connect to redis on [%s] with pass [%s]", env.StorageAddress, env.StoragePass)
	st, err := storage.NewRedisStorage(env.StorageAddress, env.StoragePass)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Connect to lines provider on [%s]", env.LinesProviderAddr)
	wm := works.NewWorkManager(st, env.LinesProviderAddr)
	wm.StartWorks(env.BTime, env.FTime, env.STime)

	log.Infof("Create GRPC listener on port [%s]", env.GRPCPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", env.GRPCPort))
	if err != nil {
		log.Fatal(err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	kiddy := grpc_srv.New(st)
	api.RegisterKiddyServer(grpcServer, kiddy)
	go grpcServer.Serve(lis)

	httpServer := http_srv.New()
	httpServer.ConfigureRouter(st)
	log.Infof("Start HTTP server on port [%s]", env.HTTPPort)
	log.Fatal(httpServer.Start(fmt.Sprintf(":%s", env.HTTPPort)))
}
