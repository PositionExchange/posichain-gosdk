package account

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mitchellh/go-homedir"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/PositionExchange/posichain-gosdk/pkg/common"
	"github.com/PositionExchange/posichain-gosdk/pkg/mnemonic"
	"github.com/PositionExchange/posichain-gosdk/pkg/store"
	"github.com/PositionExchange/posichain/accounts/keystore"
	"github.com/btcsuite/btcd/btcec"
	mapset "github.com/deckarep/golang-set"
)

// ImportFromPrivateKey allows import of an ECDSA private key
func ImportFromPrivateKey(privateKey, name, passphrase string) (string, error) {
	privateKey = strings.TrimPrefix(privateKey, "0x")

	if name == "" {
		name = generateName() + "-imported"
		for store.DoesNamedAccountExist(name) {
			name = generateName() + "-imported"
		}
	} else if store.DoesNamedAccountExist(name) {
		return "", fmt.Errorf("account %s already exists", name)
	}

	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	if len(privateKeyBytes) != common.Secp256k1PrivateKeyBytesLength {
		return "", common.ErrBadKeyLength
	}

	// btcec.PrivKeyFromBytes only returns a secret key and public key
	sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes)
	hexAddress := crypto.PubkeyToAddress(sk.PublicKey).Hex()

	if store.FromAddress(hexAddress) != nil {
		return "", fmt.Errorf("address %s already exists", hexAddress)
	}

	ks := store.FromAccountName(name)
	_, err = ks.ImportECDSA(sk.ToECDSA(), passphrase)
	return name, err
}

func generateName() string {
	words := strings.Split(mnemonic.Generate(), " ")
	existingAccounts := mapset.NewSet()
	for a := range store.LocalAccounts() {
		existingAccounts.Add(a)
	}
	foundName := false
	acct := ""
	i := 0
	for {
		if foundName {
			break
		}
		if i == len(words)-1 {
			words = strings.Split(mnemonic.Generate(), " ")
		}
		candidate := words[i]
		if !existingAccounts.Contains(candidate) {
			foundName = true
			acct = candidate
			break
		}
	}
	return acct
}

func writeToFile(path string, data string) error {
	currDir, _ := os.Getwd()
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	os.MkdirAll(filepath.Dir(path), 0777)
	os.Chdir(filepath.Dir(path))
	file, err := os.Create(filepath.Base(path))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	os.Chdir(currDir)
	return file.Sync()
}

// ImportKeyStore imports a keystore along with a password
func ImportKeyStore(keyPath, name, passphrase string) (string, error) {
	keyPath, err := filepath.Abs(keyPath)
	if err != nil {
		return "", err
	}
	keyJSON, readError := ioutil.ReadFile(keyPath)
	if readError != nil {
		return "", readError
	}
	if name == "" {
		name = generateName() + "-imported"
		for store.DoesNamedAccountExist(name) {
			name = generateName() + "-imported"
		}
	} else if store.DoesNamedAccountExist(name) {
		return "", fmt.Errorf("account %s already exists", name)
	}
	key, err := keystore.DecryptKey(keyJSON, passphrase)
	if err != nil {
		return "", err
	}
	hexAddress := key.Address.Hex()
	hasAddress := store.FromAddress(hexAddress) != nil
	if hasAddress {
		return "", fmt.Errorf("address %s already exists in keystore", hexAddress)
	}
	uDir, _ := homedir.Dir()
	newPath := filepath.Join(uDir, common.DefaultConfigDirName, common.DefaultConfigAccountAliasesDirName, name, filepath.Base(keyPath))
	err = writeToFile(newPath, string(keyJSON))
	if err != nil {
		return "", err
	}
	return name, nil
}
