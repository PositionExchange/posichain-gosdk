package keys

import (
	"fmt"
	p2p_crypto "github.com/libp2p/go-libp2p-core/crypto"
)

// P2PKeyStore is used to persist private key to/from file
type P2PKeyStore struct {
	Key string `json:"key"`
}

// GenKeyP2PRand generates a pair of RSA keys used in libp2p host, using random seed
func GenKeyP2PRand() (p2p_crypto.PrivKey, p2p_crypto.PubKey, error) {
	return p2p_crypto.GenerateKeyPair(p2p_crypto.RSA, 2048)
}

// MakeP2PKeyStore save private key to keyfile
func MakeP2PKeyStore(key p2p_crypto.PrivKey) (*P2PKeyStore, error) {
	str, err := convertP2PKeyToString(key)
	if err != nil {
		return nil, err
	}

	keyStruct := P2PKeyStore{Key: str}
	return &keyStruct, nil
}

// convertP2PKeyToString convert the PrivKey to base64 format and return string
func convertP2PKeyToString(key p2p_crypto.PrivKey) (string, error) {
	if key != nil {
		b, err := p2p_crypto.MarshalPrivateKey(key)
		if err != nil {
			return "", fmt.Errorf("failed to marshal private key: %v", err)
		}
		str := p2p_crypto.ConfigEncodeKey(b)
		return str, nil
	}
	return "", fmt.Errorf("key is nil")
}
