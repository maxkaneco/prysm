package node

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/prysmaticlabs/prysm/cmd/validator/flags"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
	"github.com/prysmaticlabs/prysm/testing/require"
	"github.com/prysmaticlabs/prysm/validator/accounts"
	"github.com/prysmaticlabs/prysm/validator/accounts/wallet"
	"github.com/prysmaticlabs/prysm/validator/keymanager"
	remote_web3signer "github.com/prysmaticlabs/prysm/validator/keymanager/remote-web3signer"
	logTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/urfave/cli/v2"
)

// Test that the sharding node can build with default flag values.
func TestNode_Builds(t *testing.T) {
	app := cli.App{}
	set := flag.NewFlagSet("test", 0)
	set.String("datadir", t.TempDir()+"/datadir", "the node data directory")
	dir := t.TempDir() + "/walletpath"
	passwordDir := t.TempDir() + "/password"
	require.NoError(t, os.MkdirAll(passwordDir, os.ModePerm))
	passwordFile := filepath.Join(passwordDir, "password.txt")
	walletPassword := "$$Passw0rdz2$$"
	require.NoError(t, ioutil.WriteFile(
		passwordFile,
		[]byte(walletPassword),
		os.ModePerm,
	))
	set.String("wallet-dir", dir, "path to wallet")
	set.String("wallet-password-file", passwordFile, "path to wallet password")
	set.String("keymanager-kind", "imported", "keymanager kind")
	set.String("verbosity", "debug", "log verbosity")
	require.NoError(t, set.Set(flags.WalletPasswordFileFlag.Name, passwordFile))
	context := cli.NewContext(&app, set, nil)
	_, err := accounts.CreateWalletWithKeymanager(context.Context, &accounts.CreateWalletConfig{
		WalletCfg: &wallet.Config{
			WalletDir:      dir,
			KeymanagerKind: keymanager.Local,
			WalletPassword: walletPassword,
		},
	})
	require.NoError(t, err)

	valClient, err := NewValidatorClient(context)
	require.NoError(t, err, "Failed to create ValidatorClient")
	err = valClient.db.Close()
	require.NoError(t, err)
}

// TestClearDB tests clearing the database
func TestClearDB(t *testing.T) {
	hook := logTest.NewGlobal()
	tmp := filepath.Join(t.TempDir(), "datadirtest")
	require.NoError(t, clearDB(context.Background(), tmp, true))
	require.LogsContain(t, hook, "Removing database")
}

// TestWeb3SignerConfig tests the web3 signer config returns the correct values.
func TestWeb3SignerConfig(t *testing.T) {
	pubkey1decoded, err := hexutil.Decode("0xa99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c")
	require.NoError(t, err)
	bytepubkey1 := bytesutil.ToBytes48(pubkey1decoded)

	pubkey2decoded, err := hexutil.Decode("0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b")
	require.NoError(t, err)
	bytepubkey2 := bytesutil.ToBytes48(pubkey2decoded)

	type args struct {
		baseURL         string
		publicKeysOrURL string
	}
	tests := []struct {
		name       string
		args       args
		want       *remote_web3signer.SetupConfig
		wantErrMsg string
	}{
		{
			name: "happy path with public keys",
			args: args{
				baseURL: "http://localhost:8545",
				publicKeysOrURL: "0xa99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c," +
					"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b",
			},
			want: &remote_web3signer.SetupConfig{
				BaseEndpoint:          "http://localhost:8545",
				GenesisValidatorsRoot: nil,
				PublicKeysURL:         "",
				ProvidedPublicKeys: [][48]byte{
					bytepubkey1,
					bytepubkey2,
				},
			},
		},
		{
			name: "happy path with external url",
			args: args{
				baseURL:         "http://localhost:8545",
				publicKeysOrURL: "http://localhost:8545/api/v1/eth2/publicKeys",
			},
			want: &remote_web3signer.SetupConfig{
				BaseEndpoint:          "http://localhost:8545",
				GenesisValidatorsRoot: nil,
				PublicKeysURL:         "http://localhost:8545/api/v1/eth2/publicKeys",
				ProvidedPublicKeys:    nil,
			},
		},
		{
			name: "Bad base URL",
			args: args{
				baseURL: "0xa99a76ed7796f7be22d5b7e85deeb7c5677e88,",
				publicKeysOrURL: "0xa99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c," +
					"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b",
			},
			want:       nil,
			wantErrMsg: "web3signer url 0xa99a76ed7796f7be22d5b7e85deeb7c5677e88, is invalid: parse \"0xa99a76ed7796f7be22d5b7e85deeb7c5677e88,\": invalid URI for request",
		},
		{
			name: "Bad publicKeys",
			args: args{
				baseURL: "http://localhost:8545",
				publicKeysOrURL: "0xa99a76ed7796f7be22c," +
					"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b",
			},
			want:       nil,
			wantErrMsg: "could not decode public key for web3signer: 0xa99a76ed7796f7be22c: hex string of odd length",
		},
		{
			name: "Bad publicKeysURL",
			args: args{
				baseURL:         "http://localhost:8545",
				publicKeysOrURL: "localhost",
			},
			want:       nil,
			wantErrMsg: "could not decode public key for web3signer: localhost: hex string without 0x prefix",
		},
		{
			name: "Base URL missing scheme or host",
			args: args{
				baseURL:         "localhost:8545",
				publicKeysOrURL: "localhost",
			},
			want:       nil,
			wantErrMsg: "web3signer url must be in the format of http(s)://host:port url used: localhost:8545",
		},
		{
			name: "Public Keys URL missing scheme or host",
			args: args{
				baseURL:         "http://localhost:8545",
				publicKeysOrURL: "localhost:8545",
			},
			want:       nil,
			wantErrMsg: "could not decode public key for web3signer: localhost:8545: hex string without 0x prefix",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := web3SignerConfig(newWeb3SignerCli(t, tt.args.baseURL, tt.args.publicKeysOrURL))
			if (tt.wantErrMsg != "") && (tt.wantErrMsg != fmt.Sprintf("%v", err)) {
				t.Errorf("web3SignerConfig error = %v, wantErrMsg = %v", err, tt.wantErrMsg)
				return
			}
			require.DeepEqual(t, got, tt.want)
		})
	}
}

func newWeb3SignerCli(t *testing.T, baseUrl string, publicKeysOrURL string) *cli.Context {
	app := cli.App{}
	set := flag.NewFlagSet("test", 0)
	set.String("validators-external-signer-url", baseUrl, "baseUrl")
	set.String("validators-external-signer-public-keys", publicKeysOrURL, "publicKeys or URL")
	require.NoError(t, set.Set(flags.Web3SignerURLFlag.Name, baseUrl))
	require.NoError(t, set.Set(flags.Web3SignerPublicValidatorKeysFlag.Name, publicKeysOrURL))
	return cli.NewContext(&app, set, nil)
}
