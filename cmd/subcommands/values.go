package cmd

import (
	"fmt"

	"github.com/fatih/color"
)

const (
	hmyDocsDir             = "psc-docs"
	defaultNodeAddr        = "http://localhost:9500"
	defaultRpcPrefix       = "hmy"
	defaultMainnetEndpoint = "https://api.s0.posichain.org/"
)

var (
	g           = color.New(color.FgGreen).SprintFunc()
	cookbookDoc = fmt.Sprintf(`
Cookbook of Usage

Note:

1) Every subcommand recognizes a '--help' flag
2) If a passphrase is used by a subcommand, one can enter their own passphrase interactively
   with the --passphrase option. Alternatively, one can pass their own passphrase via a file
   using the --passphrase-file option. If no passphrase option is selected, the default
   passphrase of '' is used.
3) These examples use Shard 0 of [NETWORK] as argument for --node

Examples:

%s
./psc --node=[NODE] balances <SOME_HEX_ADDRESS>

%s
./psc --node=[NODE] blockchain transaction-by-hash <SOME_TX_HASH>

%s
./psc keys list

%s
./psc --node=[NODE] transfer \
    --from <SOME_HEX_ADDRESS> --to <SOME_HEX_ADDRESS> \
    --from-shard 0 --to-shard 1 --amount 200 --passphrase

%s
./psc --node=[NODE] transfer --file <PATH_TO_JSON_FILE>
Check README for details on json file format.

%s
./psc --node=[NODE] blockchain transaction-receipt <SOME_TX_HASH>

%s
./psc keys recover-from-mnemonic <ACCOUNT_NAME> --passphrase

%s
./psc keys import-ks <PATH_TO_KEYSTORE_JSON>

%s
./psc keys import-private-key <secp256k1_PRIVATE_KEY>

%s
./psc keys export-private-key <ACCOUNT_ADDRESS> --passphrase

%s
./psc keys generate-bls-key --bls-file-path <PATH_FOR_BLS_KEY_FILE>

%s
./psc --node=[NODE] staking create-validator --amount 10 --validator-addr <SOME_HEX_ADDRESS> \
    --bls-pubkeys <BLS_KEY_1>,<BLS_KEY_2>,<BLS_KEY_3> \
    --identity foo --details bar --name baz --max-change-rate 0.1 --max-rate 0.1 --max-total-delegation 10 \
    --min-self-delegation 10 --rate 0.1 --security-contact Leo  --website posichain.org --passphrase

%s
./psc --node=[NODE] staking edit-validator \
    --validator-addr <SOME_HEX_ADDRESS> --identity foo --details bar \
    --name baz --security-contact EK --website posichain.org \
    --min-self-delegation 0 --max-total-delegation 10 --rate 0.1\
    --add-bls-key <SOME_BLS_KEY> --remove-bls-key <OTHER_BLS_KEY> --passphrase

%s
./psc --node=[NODE] staking delegate \
    --delegator-addr <SOME_HEX_ADDRESS> --validator-addr <VALIDATOR_HEX_ADDRESS> \
    --amount 10 --passphrase

%s
./psc --node=[NODE] staking undelegate \
    --delegator-addr <SOME_HEX_ADDRESS> --validator-addr <VALIDATOR_HEX_ADDRESS> \
    --amount 10 --passphrase

%s
./psc --node=[NODE] staking collect-rewards \
    --delegator-addr <SOME_HEX_ADDRESS> --passphrase

%s
./psc --node=[NODE] blockchain validator elected

%s
./psc --node=[NODE] blockchain utility-metrics

%s
./psc --node=[NODE] failures staking

%s
./psc --node=[NODE] utility shard-for-bls <BLS_PUBLIC_KEY>

%s
./psc governance list-space

%s
./psc governance list-proposal --space=[space key, example: staking-testnet]

%s
./psc governance view-proposal --proposal=[proposal hash]

%s
./psc governance new-proposal --proposal-yaml=[file path] --key=[account address]
PS: key must first use (psc keys import-private-key) to import
Yaml example(time is in UTC timezone):
space: staking-testnet
start: 2020-04-16 21:45:12
end: 2020-04-21 21:45:12
choices:
  - yes
  - no
title: this is title
body: |
  this is body,
  you can write mutli line

%s
./psc governance vote-proposal --proposal=[proposal hash] --choice=[your choise text, eg: yes] --key=[account address]
PS: key must first use (psc keys import-private-key) to import
`,
		g("1.  Check account balance on given chain"),
		g("2.  Check sent transaction"),
		g("3.  List local account keys"),
		g("4.  Sending a transaction (waits 40 seconds for transaction confirmation)"),
		g("5.  Sending a batch of transactions as dictated from a file (the `--dry-run` options still apply)"),
		g("6.  Check a completed transaction receipt"),
		g("7.  Import an account using the mnemonic. Prompts the user to give the mnemonic."),
		g("8.  Import an existing keystore file"),
		g("9.  Import a keystore file using a secp256k1 private key"),
		g("10. Export a keystore file's secp256k1 private key"),
		g("11. Generate a BLS key then encrypt and save the private key to the specified location."),
		g("12. Create a new validator with a list of BLS keys"),
		g("13. Edit an existing validator"),
		g("14. Delegate an amount to a validator"),
		g("15. Undelegate to a validator"),
		g("16. Collect block rewards as a delegator"),
		g("17. Check elected validators"),
		g("18. Get current staking utility metrics"),
		g("19. Check in-memory record of failed staking transactions"),
		g("20. Check which shard your BLS public key would be assigned to as a validator"),
		g("21. List Space In Governance"),
		g("22. List Proposal In Space Of Governance"),
		g("23. View Proposal In Governance"),
		g("24. New Proposal In Space Of Governance"),
		g("25. Vote Proposal In Space Of Governance"),
	)
)
