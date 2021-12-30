package grpcserver

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

// HandleSigGracefulShutdown handle graceful shutdown of the server
func HandleSigGracefulShutdown(server *grpc.Server) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	if sig == syscall.SIGTERM {
		//nolint:forbidigo // We don't have a logger here
		fmt.Println("Received a SIGINT, shutting down immediately!")
		os.Exit(1)
	} else {
		//nolint:forbidigo // We don't have a logger here
		fmt.Println("Received  a SIGINT, shutting down gracefully!")
		server.GracefulStop()
	}
}
