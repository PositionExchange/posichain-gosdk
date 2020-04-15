package account

import (
	"fmt"
	"path/filepath"

	"github.com/harmony-one/go-sdk/pkg/store"
	"github.com/harmony-one/harmony/accounts"
)

func ExportPrivateKey(address, passphrase string) error {
	ks := store.FromAddress(address)
	allAccounts := ks.Accounts()
	for _, account := range allAccounts {
		_, key, err := ks.GetDecryptedKey(accounts.Account{Address: account.Address}, passphrase)
		if err != nil {
			return err
		}
		fmt.Printf("%064x\n", key.PrivateKey.D)
	}
	return nil
}

func ExportKeystore(address, path, passphrase string) (string, error) {
	ks := store.FromAddress(address)
	allAccounts := ks.Accounts()
	dirPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	outFile := filepath.Join(dirPath, fmt.Sprintf("%s.key", address))
	for _, account := range allAccounts {
		keyFile, err := ks.Export(accounts.Account{Address: account.Address}, passphrase, passphrase)
		if err != nil {
			return "", err
		}
		e := writeToFile(outFile, string(keyFile))
		if e != nil {
			return "", e
		}
	}
	return outFile, nil
}
