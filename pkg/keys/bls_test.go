package keys

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestBlsKeyGeneration(t *testing.T) {
	valid := true
	passphrase := ""
	folderPath := "TestBlsKeyGeneration"
	absFolderPath, err := filepath.Abs(fmt.Sprintf("./%s", folderPath))
	if err != nil {
		t.Errorf("TestBlsKeyGeneration - failed to convert relative path to absolute path")
	}

	absFilePath := fmt.Sprintf("%s/TestBlsKeyGeneration.key", absFolderPath)

	err = os.MkdirAll(absFolderPath, 0755)
	if err != nil {
		t.Errorf("TestBlsKeyGeneration - failed to make test key folder")
	}

	if err := GenBlsKey(&BlsKey{Passphrase: passphrase, FilePath: absFilePath}); err != nil {
		t.Errorf("TestBlsKeyGeneration - failed to generate bls key using passphrase %s and path %s", passphrase, absFilePath)
	}

	if _, err = os.Stat(absFilePath); err != nil {
		t.Errorf("TestBlsKeyGeneration - failed to check if file %s exists", absFilePath)
	}

	valid = !os.IsNotExist(err)

	if !valid {
		t.Errorf("GenBlsKey - failed to generate a bls key using passphrase %s", "")
	}

	os.RemoveAll(absFolderPath)
}

func TestMultiBlsKeyGeneration(t *testing.T) {
	tests := []struct {
		count      uint32
		shardID    uint32
		shardCount int
		filePath   string
		expected   bool
	}{
		{shardCount: 4, count: 3, shardID: 0, expected: true},
		{shardCount: 4, count: 3, shardID: 4, expected: false},
	}

	for _, test := range tests {
		valid := false
		blsKeys := []*BlsKey{}
		for i := uint32(0); i < test.count; i++ {
			blsKeys = append(blsKeys, &BlsKey{Passphrase: "", FilePath: ""})
		}

		blsKeys, shardCount, err := genBlsKeyForShard(blsKeys, test.shardCount, test.shardID)
		if err != nil {
			valid = false
		}

		successCount := 0

		for _, blsKey := range blsKeys {
			success := (blsKey.ShardPublicKey != nil && blsKeyMatchesShardID(blsKey.ShardPublicKey, test.shardID, shardCount))
			if success {
				successCount++
			}
		}

		valid = (successCount == int(test.count))

		if valid != test.expected {
			t.Errorf("genBlsKeyForShard - failed to generate %d keys for shard %d using node %s", test.count, test.shardID, test.shardCount)
		}
	}
}
