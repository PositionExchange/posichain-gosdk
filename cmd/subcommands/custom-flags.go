package cmd

import (
	"github.com/PositionExchange/posichain-gosdk/pkg/common"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type oneAddress struct {
	address string
}

func (oneAddress oneAddress) String() string {
	return oneAddress.address
}

func (oneAddress *oneAddress) Set(s string) error {
	if !ethCommon.IsHexAddress(s) {
		return errors.New("address is not in hex format")
	}
	oneAddress.address = s
	return nil
}

func (oneAddress oneAddress) Type() string {
	return "address"
}

type chainIDWrapper struct {
	chainID *common.ChainID
}

func (chainIDWrapper chainIDWrapper) String() string {
	return chainIDWrapper.chainID.Name
}

func (chainIDWrapper *chainIDWrapper) Set(s string) error {
	chain, err := common.StringToChainID(s)
	chainIDWrapper.chainID = chain
	if err != nil {
		return err
	}
	return nil
}

func (chainIDWrapper chainIDWrapper) Type() string {
	return "chain-id"
}
