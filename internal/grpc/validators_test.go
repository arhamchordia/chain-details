package grpc_test

import (
	"crypto/tls"
	internalgrpc "github.com/arhamchordia/chain-details/internal/grpc"
	grpctypes "github.com/arhamchordia/chain-details/types/grpc"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
	"testing"
)

func TestParseValidators(t *testing.T) {
	testCases := []struct {
		grpcUrl                                                            string
		prefix                                                             string
		dialOptions                                                        []grpc.DialOption
		expectErrorGrpc, expectErrorParseValidators, expectErrorDeleteFile bool
		errorGrpc, errorParseValidators, errorDeleteFile                   string
	}{
		{
			grpcUrl: "grpc-umee-ia.cosmosia.notional.ventures:443",
			prefix:  "umee",
			dialOptions: []grpc.DialOption{
				grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})),
			},
			expectErrorGrpc:            false,
			expectErrorParseValidators: false,
			expectErrorDeleteFile:      false,
		},
		{
			grpcUrl: "grpc-umee-ia.cosmosia.notional.venture",
			prefix:  "umee",
			dialOptions: []grpc.DialOption{
				grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})),
			},
			expectErrorGrpc:            false,
			expectErrorParseValidators: true,
			expectErrorDeleteFile:      true,
			errorParseValidators:       "rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing dial tcp: address grpc-umee-ia.cosmosia.notional.venture: missing port in address\"",
			errorDeleteFile:            "remove validators_info.csv: no such file or directory",
		},
	}

	for _, tc := range testCases {
		testGrpcConn, err := grpc.Dial(
			tc.grpcUrl,
			tc.dialOptions...,
		)
		if tc.expectErrorGrpc {
			require.Equal(t, err.Error(), tc.errorGrpc)
		}
		defer testGrpcConn.Close()

		err = internalgrpc.ParseValidators(testGrpcConn, tc.prefix)
		if tc.expectErrorParseValidators {
			require.Equal(t, err.Error(), tc.errorParseValidators)
		}

		err = os.Remove(grpctypes.PrefixGRPC + grpctypes.ValidatorsInfoFileName + ".csv")
		if tc.expectErrorDeleteFile {
			require.Equal(t, err.Error(), tc.errorDeleteFile)
		}
	}
}
