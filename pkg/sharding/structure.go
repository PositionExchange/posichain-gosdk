package sharding

import (
	"bytes"
	"fmt"
	"math/big"
	"encoding/json"

	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/common"
)

// RPCRoutes reflects the RPC endpoints of the target network across shards
type RPCRoutes struct {
	HTTP    string `json:"http"`
	ShardID int    `json:"shardID"`
	WS      string `json:"ws"`
}

// Structure produces a slice of RPCRoutes for the network across shards
func Structure(node string) ([]RPCRoutes, error) {
	type r struct {
		Result []RPCRoutes `json:"result"`
	}
	p, e := rpc.RawRequest(rpc.Method.GetShardingStructure, node, []interface{}{})
	if e != nil {
		return nil, e
	}
	result := r{}
	json.Unmarshal(p, &result)
	return result.Result, nil
}

func CheckAllShards(node, oneAddr string, noPretty bool) (string, error) {
	var out bytes.Buffer
	out.WriteString("[")
	params := []interface{}{oneAddr, "latest"}
	s, err := Structure(node)
	if err != nil {
		return "", err
	}
	count := len(s)
	for i, shard := range s {
		balanceRPCReply, err := rpc.Request(rpc.Method.GetBalance, shard.HTTP, params)
		if err != nil {
			if common.DebugRPC {
				fmt.Printf("NOTE: Route %s failed.", shard.HTTP)
			}
			count--
			continue
		}
		balance, _ := balanceRPCReply["result"].(string)
		bln, _ := big.NewInt(0).SetString(balance[2:], 16)
		out.WriteString(fmt.Sprintf(`{"shard":%d, "amount":%s}`,
			shard.ShardID,
			common.ConvertBalanceIntoReadableFormat(bln),
		))
		if i != count - 1 {
			out.WriteString(",")
		}
	}
	out.WriteString("]")
	if noPretty {
		return out.String(), nil
	}
	return common.JSONPrettyFormat(out.String()), nil
}
