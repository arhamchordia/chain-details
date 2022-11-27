package internal

import (
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"
)

func QueryValidatorsData(grpcUrl, accountPrefix string) error {
	// initialise config for grpc connection
	config := &tls.Config{
		InsecureSkipVerify: false,
	}

	// Create a connection to the gRPC server.
	grpcConn, err := grpc.Dial(
		grpcUrl,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
	)
	if err != nil {
		return err
	}
	defer grpcConn.Close()

	// send a query only when connection state is ready
	for {
		// wait for 4 milliseconds for grpc to connect
		time.Sleep(4 * time.Millisecond)

		if grpcConn.GetState().String() == "READY" {
			err = ParseValidators(grpcConn, accountPrefix)
			if err != nil {
				return err
			}
			break
		} else if grpcConn.GetState().String() == "TRANSIENT_FAILURE" {
			break
		}
	}

	return nil
}

func QueryDelegatorsData(grpcUrl string) error {
	// initialise config for grpc connection
	config := &tls.Config{
		InsecureSkipVerify: false,
	}

	// Create a connection to the gRPC server.
	grpcConn, err := grpc.Dial(
		grpcUrl,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
	)
	if err != nil {
		return err
	}
	defer grpcConn.Close()

	// send a query only when connection state is ready
	for {
		// wait for 4 milliseconds for grpc to connect
		time.Sleep(4 * time.Millisecond)

		// trigger action on the basis of state of the connection
		if grpcConn.GetState().String() == "READY" {
			err = ParseDelegators(grpcConn)
			if err != nil {
				return err
			}
			break
		} else if grpcConn.GetState().String() == "TRANSIENT_FAILURE" {
			break
		}
	}

	return nil
}
