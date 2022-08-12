package address

import (
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type T = ethCommon.Address

// MustParse parses the given address as hex.
// Panic if the passing param is not a correct address.
func MustParse(s string) T {
	addr, err := ParseHex(s)
	if err != nil {
		panic(err)
	}
	return addr
}

func ParseHex(s string) (T, error) {
	if !ethCommon.IsHexAddress(s) {
		return T{}, errors.New("address is not in hex format")
	}
	return ethCommon.HexToAddress(s), nil
}
