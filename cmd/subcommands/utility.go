package cmd

import (
	"encoding/json"
	"fmt"
	bls_core "github.com/PositionExchange/bls/ffi/go/bls"
	"github.com/PositionExchange/posichain-gosdk/pkg/rpc"
	"github.com/PositionExchange/posichain/crypto/bls"
	"github.com/spf13/cobra"
	"math/big"
	"strings"
)

func init() {
	cmdUtilities := &cobra.Command{
		Use:   "utility",
		Short: "common posichain utilities",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmdUtilities.AddCommand(&cobra.Command{
		Use:   "metadata",
		Short: "data includes network specific values",
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetNodeMetadata, []interface{}{})
		},
	})

	cmdMetrics := &cobra.Command{
		Use:   "metrics",
		Short: "mostly in-memory fluctuating values",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmdMetrics.AddCommand([]*cobra.Command{{
		Use:   "pending-crosslinks",
		Short: "dump the pending crosslinks in memory of target node",
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetPendingCrosslinks, []interface{}{})
		},
	}, {
		Use:   "pending-cx-receipts",
		Short: "dump the pending cross shard receipts in memory of target node",
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetPendingCXReceipts, []interface{}{})
		},
	},
	}...)

	cmdShardForBls := &cobra.Command{
		// Temp utility that should be removed once resharding implemented
		Use:   "shard-for-bls",
		Args:  cobra.ExactArgs(1),
		Short: "which shard this BLS key would be assigned to",
		RunE: func(cmd *cobra.Command, args []string) error {
			inputKey := strings.TrimPrefix(args[0], "0x")
			key := bls_core.PublicKey{}
			if err := key.DeserializeHexStr(inputKey); err != nil {
				return err
			}
			shardBig := shardCount
			if shardCount <= 0 {
				reply, err := rpc.Request(rpc.Method.GetShardingStructure, node, []interface{}{})
				if err != nil {
					return err
				}
				shardBig = len(reply["result"].([]interface{})) // assume the response is a JSON Array
			}
			wrapper := bls.FromLibBLSPublicKeyUnsafe(&key)
			shardID := int(new(big.Int).Mod(wrapper.Big(), big.NewInt(int64(shardBig))).Int64())
			type t struct {
				ShardID int `json:"shard-id"`
			}
			result, err := json.Marshal(t{shardID})
			if err != nil {
				return err
			}

			fmt.Println(string(result))
			return nil
		},
	}
	cmdShardForBls.Flags().IntVar(&shardCount, "shard-count", 0, "how many shard in total")

	cmdUtilities.AddCommand(cmdMetrics)
	cmdUtilities.AddCommand(cmdShardForBls)
	cmdUtilities.AddCommand([]*cobra.Command{{
		Use:   "committees",
		Short: "current and previous committees",
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetSuperCommmittees, []interface{}{})
		},
	}, {
		Use:   "bad-blocks",
		Short: "bad blocks in memory of the target node",
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetCurrentBadBlocks, []interface{}{})
		},
	}, {
		Use:   "shards",
		Short: "sharding structure and end points",
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetShardingStructure, []interface{}{})
		},
	}, {
		Use:   "last-cross-links",
		Short: "last crosslinks processed",
		RunE: func(cmd *cobra.Command, args []string) error {
			noLatest = true
			return request(rpc.Method.GetLastCrossLinks, []interface{}{})
		},
	}}...)

	RootCmd.AddCommand(cmdUtilities)
}
