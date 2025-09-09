package grpc

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func Grpc(GRPCAddr string, TLSCertFile string, TLSKeyFile string) *grpc.Server {
	// Charger cfg/env (DB, JWT keys, RPID, Origin, etc.)

	// TLS (prod) : use cert auto (Caddy) devant ou charger ici
	var serverOpts []grpc.ServerOption
	if TLSCertFile != "" {
		creds, err := credentials.NewServerTLSFromFile(TLSCertFile, TLSKeyFile)
		if err != nil {
			log.Fatal(err)
		}
		serverOpts = append(serverOpts, grpc.Creds(creds))
	}
	// Interceptors: logging, tracing, auth, rate-limit

	app := grpc.NewServer(serverOpts...)

	// Services

	// Health + reflection (utile pour debug/evans)
	healthSrv := health.NewServer()
	healthpb.RegisterHealthServer(app, healthSrv)
	reflection.Register(app)

	return app
}
