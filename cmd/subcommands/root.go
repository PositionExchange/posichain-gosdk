package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PositionExchange/posichain-gosdk/pkg/address"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/PositionExchange/posichain-gosdk/pkg/common"
	"github.com/PositionExchange/posichain-gosdk/pkg/rpc"
	rpcEth "github.com/PositionExchange/posichain-gosdk/pkg/rpc/eth"
	rpcV1 "github.com/PositionExchange/posichain-gosdk/pkg/rpc/v1"
	"github.com/PositionExchange/posichain-gosdk/pkg/sharding"
	"github.com/PositionExchange/posichain-gosdk/pkg/store"
	color "github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	verbose         bool
	useLedgerWallet bool
	noLatest        bool
	noPrettyOutput  bool
	node            string
	rpcPrefix       string
	keyStoreDir     string
	givenFilePath   string
	endpoint        = regexp.MustCompile(`https://s[0-9]\..*\.posichain\.org`)
	request         = func(method string, params []interface{}) error {
		if !noLatest {
			params = append(params, "latest")
		}
		success, failure := rpc.Request(method, node, params)
		if failure != nil {
			return failure
		}
		asJSON, _ := json.Marshal(success)
		if noPrettyOutput {
			fmt.Println(string(asJSON))
			return nil
		}
		fmt.Println(common.JSONPrettyFormat(string(asJSON)))
		return nil
	}
	// RootCmd is single entry point of the CLI
	RootCmd = &cobra.Command{
		Use:          "psc",
		Short:        "Posichain",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if verbose {
				common.EnableAllVerbose()
			}
			switch rpcPrefix {
			case "hmy":
				rpc.Method = rpcV1.Method
			case "eth":
				rpc.Method = rpcEth.Method
			default:
				rpc.Method = rpcV1.Method
			}
			if strings.HasPrefix(node, "https://") || strings.HasPrefix(node, "http://") ||
				strings.HasPrefix(node, "ws://") {
				//No op, already has protocol, respect protocol default ports.
			} else {
				return errors.New("node must start with protocol http(s) or ws")
			}

			if targetChain == "" {
				if node == defaultNodeAddr {
					routes, err := sharding.Structure(node)
					if err != nil {
						chainName = chainIDWrapper{chainID: &common.Chain.TestNet}
					} else {
						if len(routes) == 0 {
							return errors.New("empty reply from sharding structure")
						}
						chainName = endpointToChainID(routes[0].HTTP)
					}
				} else if endpoint.Match([]byte(node)) {
					chainName = endpointToChainID(node)
				} else {
					chainName = chainIDWrapper{chainID: &common.Chain.TestNet}
				}
			} else {
				chain, err := common.StringToChainID(targetChain)
				if err != nil {
					return err
				}
				chainName = chainIDWrapper{chainID: chain}
			}

			return nil
		},
		Long: fmt.Sprintf(`
CLI interface to the Posichain

%s`, g("Invoke 'psc cookbook' for examples of the most common, important usages")),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}
)

func init() {
	vS := "dump out debug information, same as env var PSC_ALL_DEBUG=true"
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, vS)
	RootCmd.PersistentFlags().StringVarP(&node, "node", "n", defaultNodeAddr, "<host>")
	RootCmd.PersistentFlags().StringVarP(&rpcPrefix, "rpc-prefix", "r", defaultRpcPrefix, "<rpc>")
	RootCmd.PersistentFlags().BoolVar(
		&noLatest, "no-latest", false, "Do not add 'latest' to RPC params",
	)
	RootCmd.PersistentFlags().BoolVar(
		&noPrettyOutput, "no-pretty", false, "Disable pretty print JSON outputs",
	)
	RootCmd.AddCommand(&cobra.Command{
		Use:   "cookbook",
		Short: "Example usages of the most important, frequently used commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			var docNode, docNet string
			if node == defaultNodeAddr || chainName.chainID == &common.Chain.MainNet {
				docNode = `https://s0.posichain.org`
				docNet = `Mainnet`
			} else if chainName.chainID == &common.Chain.TestNet {
				docNode = `https://s0.t.posichain.org`
				docNet = `Long-Running Testnet`
			} else if chainName.chainID == &common.Chain.DevNet {
				docNode = `https://s0.d.posichain.org`
				docNet = `Long-Running Devnet`
			} else if chainName.chainID == &common.Chain.StressNet {
				docNode = `https://s0.s.posichain.org`
				docNet = `Stress Testing Network`
			}
			fmt.Print(strings.ReplaceAll(strings.ReplaceAll(cookbookDoc, `[NODE]`, docNode), `[NETWORK]`, docNet))
			return nil
		},
	})
	RootCmd.PersistentFlags().BoolVarP(&useLedgerWallet, "ledger", "e", false, "Use ledger hardware wallet")
	RootCmd.PersistentFlags().StringVar(&givenFilePath, "file", "", "Path to file for given command when applicable")
	RootCmd.AddCommand(&cobra.Command{
		Use:   "docs",
		Short: fmt.Sprintf("Generate docs to a local %s directory", hmyDocsDir),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, _ := os.Getwd()
			docDir := path.Join(cwd, hmyDocsDir)
			os.Mkdir(docDir, 0700)
			doc.GenMarkdownTree(RootCmd, docDir)
			return nil
		},
	})
}

var (
	// VersionWrapDump meant to be set from main.go
	VersionWrapDump = ""
	cookbook        = color.GreenString("psc cookbook")
	versionLink     = "https://harmony.one/hmycli_ver"
	versionFormat   = regexp.MustCompile("v[0-9]+-[a-z0-9]{7}")
)

// Execute kicks off the hmy CLI
func Execute() {
	RootCmd.SilenceErrors = true
	if err := RootCmd.Execute(); err != nil {
		resp, httpErr := http.Get(versionLink)
		if httpErr != nil {
			return
		}
		defer resp.Body.Close()
		// If error, no op
		if resp != nil && resp.StatusCode == 200 {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)

			currentVersion := versionFormat.FindAllString(buf.String(), 1)
			if currentVersion != nil && currentVersion[0] != VersionWrapDump {
				warnMsg := fmt.Sprintf("Warning: Using outdated version. Redownload to upgrade to %s\n", currentVersion[0])
				fmt.Fprintf(os.Stderr, color.RedString(warnMsg))
			}
		}
		errMsg := errors.Wrapf(err, "commit: %s, error", VersionWrapDump).Error()
		fmt.Fprintf(os.Stderr, errMsg+"\n")
		fmt.Fprintf(os.Stderr, "check "+cookbook+" for valid examples or try adding a `--help` flag\n")
		os.Exit(1)
	}
}

func endpointToChainID(nodeAddr string) chainIDWrapper {
	if strings.Contains(nodeAddr, ".t.") {
		return chainIDWrapper{chainID: &common.Chain.TestNet}
	} else if strings.Contains(nodeAddr, ".d.") {
		return chainIDWrapper{chainID: &common.Chain.DevNet}
	} else if strings.Contains(nodeAddr, ".s.") {
		return chainIDWrapper{chainID: &common.Chain.StressNet}
	}
	return chainIDWrapper{chainID: &common.Chain.MainNet}
}

func validateAddress(cmd *cobra.Command, args []string) error {
	// Check if input valid one address
	tmpAddr := oneAddress{}
	if err := tmpAddr.Set(args[0]); err != nil {
		// Check if input is valid account name
		if acc, err := store.AddressFromAccountName(args[0]); err == nil {
			addr = oneAddress{acc}
			return nil
		}

		bech32Addr := address.ToBech32(address.Parse(args[0]))
		if bech32Addr == "" {
			return fmt.Errorf("Invalid one address/Invalid account name: %s", args[0])
		}

		tmpAddr = oneAddress{bech32Addr}
	}
	addr = tmpAddr
	return nil
}
