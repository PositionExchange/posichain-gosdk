package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/PositionExchange/posichain-gosdk/pkg/address"
	"github.com/PositionExchange/posichain-gosdk/pkg/common"
	c "github.com/PositionExchange/posichain-gosdk/pkg/common"
	"github.com/PositionExchange/posichain/accounts"
	"github.com/PositionExchange/posichain/accounts/keystore"
	"github.com/pkg/errors"

	"github.com/mitchellh/go-homedir"
)

func init() {
	uDir, _ := homedir.Dir()
	hmyCLIDir := path.Join(uDir, common.DefaultConfigDirName, common.DefaultConfigAccountAliasesDirName)
	if _, err := os.Stat(hmyCLIDir); os.IsNotExist(err) {
		os.MkdirAll(hmyCLIDir, 0700)
	}
}

// LocalAccounts returns a slice of local account alias names
func LocalAccounts() []string {
	uDir, _ := homedir.Dir()
	files, _ := ioutil.ReadDir(path.Join(
		uDir,
		common.DefaultConfigDirName,
		common.DefaultConfigAccountAliasesDirName,
	))
	var accList []string
	for _, node := range files {
		if node.IsDir() {
			accList = append(accList, path.Base(node.Name()))
		}
	}
	return accList
}

var (
	describe              = fmt.Sprintf("%-24s\t\t%23s\n", "NAME", "ADDRESS")
	NoUnlockBadPassphrase = errors.New("could not unlock wallet with given passphrase")
)

// DescribeLocalAccounts will display all the account alias name and their corresponding hex address
func DescribeLocalAccounts() {
	fmt.Println(describe)
	for _, name := range LocalAccounts() {
		ks := FromAccountName(name)
		allAccounts := ks.Accounts()
		for _, account := range allAccounts {
			fmt.Printf("%-48s\t%s\n", name, account.Address.Hex())
		}
	}
}

// DoesNamedAccountExist return true if the given string name is an alias account already define,
// and return false otherwise
func DoesNamedAccountExist(name string) bool {
	for _, account := range LocalAccounts() {
		if account == name {
			return true
		}
	}
	return false
}

// AddressFromAccountName Returns hex address for account name if exists
func AddressFromAccountName(name string) (string, error) {
	ks := FromAccountName(name)
	// FIXME: Assume 1 account per keystore for now
	for _, account := range ks.Accounts() {
		return account.Address.Hex(), nil
	}
	return "", fmt.Errorf("keystore not found")
}

// FromAddress will return nil if the hex address is not found in the imported accounts
func FromAddress(hexAddress string) *keystore.KeyStore {
	for _, name := range LocalAccounts() {
		ks := FromAccountName(name)
		allAccounts := ks.Accounts()
		for _, account := range allAccounts {
			if hexAddress == account.Address.Hex() {
				return ks
			}
		}
	}
	return nil
}

func FromAccountName(name string) *keystore.KeyStore {
	uDir, _ := homedir.Dir()
	p := path.Join(uDir, c.DefaultConfigDirName, c.DefaultConfigAccountAliasesDirName, name)
	return common.KeyStoreForPath(p)
}

func DefaultLocation() string {
	uDir, _ := homedir.Dir()
	return path.Join(uDir, c.DefaultConfigDirName, c.DefaultConfigAccountAliasesDirName)
}

func UnlockedKeystore(from, passphrase string) (*keystore.KeyStore, *accounts.Account, error) {
	return UnlockedKeystoreTimeLimit(from, passphrase, 0)
}

func LockKeystore(from string) (*keystore.KeyStore, *accounts.Account, error) {
	sender := address.MustParse(from)
	ks := FromAddress(sender.Hex())
	if ks == nil {
		return nil, nil, fmt.Errorf("could not open local keystore for %s", from)
	}
	account, lookupErr := ks.Find(accounts.Account{Address: sender})
	if lookupErr != nil {
		return nil, nil, fmt.Errorf("could not find %s in keystore", from)
	}
	if lockError := ks.Lock(account.Address); lockError != nil {
		return nil, nil, lockError
	}
	return ks, &account, nil
}

func UnlockedKeystoreTimeLimit(from, passphrase string, time time.Duration) (*keystore.KeyStore, *accounts.Account, error) {
	sender := address.MustParse(from)
	ks := FromAddress(sender.Hex())
	if ks == nil {
		return nil, nil, fmt.Errorf("could not open local keystore for %s", from)
	}
	account, lookupErr := ks.Find(accounts.Account{Address: sender})
	if lookupErr != nil {
		return nil, nil, fmt.Errorf("could not find %s in keystore", from)
	}
	if unlockError := ks.TimedUnlock(account, passphrase, time); unlockError != nil {
		return nil, nil, errors.Wrap(NoUnlockBadPassphrase, unlockError.Error())
	}
	return ks, &account, nil
}
