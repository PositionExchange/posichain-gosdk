package account

import (
	"fmt"
	"os"
	"path"

	"github.com/PositionExchange/posichain-gosdk/pkg/common"
	"github.com/PositionExchange/posichain-gosdk/pkg/store"
	"github.com/mitchellh/go-homedir"
)

// RemoveAccount - removes an account from the keystore
func RemoveAccount(name string) error {
	accountExists := store.DoesNamedAccountExist(name)

	if !accountExists {
		return fmt.Errorf("account %s doesn't exist", name)
	}

	uDir, _ := homedir.Dir()
	hmyCLIDir := path.Join(uDir, common.DefaultConfigDirName, common.DefaultConfigAccountAliasesDirName)
	accountDir := fmt.Sprintf("%s/%s", hmyCLIDir, name)
	os.RemoveAll(accountDir)

	return nil
}
