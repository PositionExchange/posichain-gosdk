package common

import (
	"github.com/PositionExchange/posichain/accounts/keystore"
)

func KeyStoreForPath(p string) *keystore.KeyStore {
	return keystore.NewKeyStore(p, ScryptN, ScryptP)
}
