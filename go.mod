module github.com/PositionExchange/posichain-gosdk

go 1.14

require (
	github.com/PositionExchange/bls v0.0.0-20210728190118-2b7e49894c0f
	github.com/PositionExchange/posichain v0.0.12
	github.com/aristanetworks/goarista v0.0.0-20191023202215-f096da5361bb // indirect
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	github.com/cosmos/cosmos-sdk v0.37.0
	github.com/deckarep/golang-set v1.7.1
	github.com/dop251/goja v0.0.0-20210427212725-462d53687b0d
	github.com/ethereum/go-ethereum v1.9.25
	github.com/fatih/color v1.10.0
	github.com/golang/snappy v0.0.2-0.20200707131729-196ae77b8a26 // indirect
	github.com/hashicorp/go-version v1.2.0 // indirect
	github.com/karalabe/hid v1.0.0
	github.com/libp2p/go-libp2p-core v0.8.6
	github.com/mattn/go-colorable v0.1.9
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/peterh/liner v1.1.1-0.20190123174540-a2c9a5303de7
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v0.0.5
	github.com/tyler-smith/go-bip39 v1.0.2
	github.com/valyala/fasthttp v1.2.0
	github.com/valyala/fastjson v1.6.3
	golang.org/x/crypto v0.0.0-20210506145944-38f3c27a63bf
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/ethereum/go-ethereum => github.com/ethereum/go-ethereum v1.9.9

replace github.com/fatih/color => github.com/fatih/color v1.13.0
