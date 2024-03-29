package common

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
)

// ChainID is a wrapper around the human name for a chain and the actual Big.Int used
type ChainID struct {
	Name  string   `json:"-"`
	Value *big.Int `json:"chain-as-number"`
}

type chainIDList struct {
	MainNet   ChainID `json:"mainnet"`
	TestNet   ChainID `json:"testnet"`
	DevNet    ChainID `json:"devnet"`
	StressNet ChainID `json:"stress"`
	DockerNet ChainID `json:"dockernet"`
}

// Chain is an enumeration of the known Chain-IDs
var Chain = chainIDList{
	MainNet:   ChainID{"mainnet", big.NewInt(1)},
	TestNet:   ChainID{"testnet", big.NewInt(2)},
	DevNet:    ChainID{"devnet", big.NewInt(3)},
	StressNet: ChainID{"stress", big.NewInt(5)},
	DockerNet: ChainID{"dockernet", big.NewInt(8)},
}

func (c chainIDList) String() string {
	s, _ := json.MarshalIndent(c, "", "  ")
	return string(s)
}

// StringToChainID returns the ChainID wrapper for the given human name of a chain-id
func StringToChainID(name string) (*ChainID, error) {
	switch name {
	case "mainnet":
		return &Chain.MainNet, nil
	case "testnet":
		return &Chain.TestNet, nil
	case "devnet":
		return &Chain.DevNet, nil
	case "dockernet":
		return &Chain.DockerNet, nil
	case "stressnet":
		return &Chain.StressNet, nil
	case "dryrun":
		return &Chain.MainNet, nil
	default:
		if chainID, err := strconv.Atoi(name); err == nil && chainID >= 0 {
			return &ChainID{Name: fmt.Sprintf("%d", chainID), Value: big.NewInt(int64(chainID))}, nil
		}
		return nil, fmt.Errorf("unknown chain-id: %s", name)
	}
}
