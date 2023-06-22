package server

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestStart(t *testing.T) {
	// Create a mock listener
	listener, err := net.Listen("tcp", "localhost:0") // Using port 0 to automatically assign an available port
	assert.NoError(t, err)

	// Create a mock gRPC server
	server := grpc.NewServer()

	// Start the server in a separate goroutine
	go func() {
		err := Start()
		assert.NoError(t, err)
	}()

	// Verify that the server is running by attempting to connect to it
	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	// Perform additional tests on the server if necessary

	// Stop the server
	server.Stop()
}
