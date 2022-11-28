package internal_test

import (
	"crypto/tls"
	"github.com/arhamchordia/chain-details/internal"
	"github.com/arhamchordia/chain-details/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
	"testing"
)

func assertNoErr(t *testing.T, e error) {
	if e != nil {
		t.Error(e)
	}
}

func TestParseDelegators(t *testing.T) {
	testCases := []struct {
		grpcUrl                                                                        string
		dialOptions                                                                    []grpc.DialOption
		expectErrorGrpc, expectErrorParseDelegators, expectErrorDeleteFile             bool
		errorGrpc, errorParseDelegators, errorDeleteEntriesFile, errorDeleteSharesFile string
	}{
		{
			grpcUrl: "grpc-umee-ia.cosmosia.notional.ventures:443",
			dialOptions: []grpc.DialOption{
				grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})),
			},
			expectErrorGrpc:            false,
			expectErrorParseDelegators: false,
			expectErrorDeleteFile:      false,
		},
		{
			grpcUrl: "grpc-umee-ia.cosmosia.notional.ventures",
			dialOptions: []grpc.DialOption{
				grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})),
			},
			expectErrorGrpc:            false,
			expectErrorParseDelegators: true,
			expectErrorDeleteFile:      true,
			errorParseDelegators:       "rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing dial tcp: address grpc-umee-ia.cosmosia.notional.ventures: missing port in address\"",
			errorDeleteEntriesFile:     "remove delegator_delegation_entries.csv: no such file or directory",
			errorDeleteSharesFile:      "remove delegator_shares.csv: no such file or directory",
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

		err = internal.ParseDelegators(testGrpcConn)
		if tc.expectErrorParseDelegators {
			require.Equal(t, err.Error(), tc.errorParseDelegators)
		}

		err = os.Remove(types.DelegatorDelegationEntriesFileName + ".csv")
		if tc.expectErrorDeleteFile {
			require.Equal(t, err.Error(), tc.errorDeleteEntriesFile)
		}

		err = os.Remove(types.DelegatorSharesFileName + ".csv")
		if tc.expectErrorDeleteFile {
			require.Equal(t, err.Error(), tc.errorDeleteSharesFile)
		}
	}
}
