package cmd

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/ledger"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/store"
	"strings"

	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/harmony/common/denominations"
	"github.com/harmony-one/harmony/numeric"
	staking "github.com/harmony-one/harmony/staking/types"
	"github.com/spf13/cobra"
	"math/big"
)

const (
	blsPubKeySize       = 48
)

var (
	validatorName             string
	validatorIdentity         string
	validatorWebsite          string
	validatorSecurityContact  string
	validatorDetails          string
	commisionRateStr        string
	commisionMaxRateStr       string
	commisionMaxChangeRateStr string
	minSelfDelegation         float64
	stakingBlsPubKey          string
	stakingAddress            oneAddress
	delegatorAddress          oneAddress
	validatorAddress          oneAddress
	validatorSrcAddress       oneAddress
	validatorDstAddress       oneAddress
	senderAddress             oneAddress
	stakingAmount             float64
)

func getNextNonce(messenger rpc.T) uint64 {
	transactionCountRPCReply, err :=
		messenger.SendRPC(rpc.Method.GetTransactionCount, []interface{}{accounts.ParseAddrH(delegatorAddress.String()), "latest"})

	if err != nil {
		return 0
	}

	transactionCount, _ := transactionCountRPCReply["result"].(string)
	nonce, _ := big.NewInt(0).SetString(transactionCount[2:], 16)
	return nonce.Uint64()
}

func createStakingTransaction(nonce uint64, f staking.StakeMsgFulfiller) (*staking.StakingTransaction, error) {
	gasPrice := big.NewInt(int64(gasPrice))
	gasPrice = gasPrice.Mul(gasPrice, big.NewInt(denominations.Nano))

	//TODO: modify the gas limit calculation algorithm
	gasLimit, err := core.IntrinsicGas(nil, false, true)
	if err != nil {
		return nil, err
	}

	stakingTx, err := staking.NewStakingTransaction(nonce, gasLimit, gasPrice, f)
	return stakingTx, err
}

func handleStakingTransaction(stakingTx *staking.StakingTransaction, networkHandler *rpc.HTTPMessenger, signerAddress oneAddress) error {
	var ks      *keystore.KeyStore
	var acct    *accounts.Account
	var signed  *staking.StakingTransaction
	var err      error

	from := signerAddress.String()

	if useLedgerWallet {
		var signerAddr string
		signed, signerAddr, err = ledger.SignStakingTx(stakingTx,  chainName.chainID.Value)
		if err != nil {
			return err
		}

		if strings.Compare(signerAddr, delegatorAddress.String()) != 0 {
			return errors.New("error : delegator address doesn't match with ledger hardware addresss")
		}
	} else {
		ks, acct, err = store.UnlockedKeystore(from, unlockP)
		if err != nil {
			return err
		}
		signed, err = ks.SignStakingTx(*acct, stakingTx, chainName.chainID.Value)
	}

	if err != nil {
		return err
	}

	enc, err := rlp.EncodeToBytes(signed)
	if err != nil {
		return err
	}

	hexSignature := hexutil.Encode(enc)
	reply, err := networkHandler.SendRPC(rpc.Method.SendRawStakingTransaction, []interface{}{hexSignature})
	if err != nil {
		return err
	}
	r, _ := reply["result"].(string)
	fmt.Println(fmt.Sprintf(`{"transaction-receipt":"%s"}`, r))
	return nil
}

func stakingSubCommands() []*cobra.Command {

	subCmdNewValidator := &cobra.Command{
		Use:   "newvalidator",
		Short: "create a new validator",
		Long: `
Create a new validator"
`,
		RunE: func(cmd *cobra.Command, args []string)  error {
			networkHandler, err := handlerForShard(0, node)
			if err != nil {
				return err
			}

			commisionRate, err := numeric.NewDecFromStr(commisionRateStr)
			if err != nil {
				return err
			}

			commisionMaxRate, err := numeric.NewDecFromStr(commisionMaxRateStr)
			if err != nil {
				return err
			}

			commisionMaxChangeRate, err := numeric.NewDecFromStr(commisionMaxChangeRateStr)
			if err != nil {
				return err
			}

			if len(stakingBlsPubKey) != blsPubKeySize {
				return errors.New("staking BLS pubkey key size should be 48 bytes")
			}

			delegateStakePayloadMaker := func() (staking.Directive, interface{}) {
				amountBigInt := big.NewInt(int64(stakingAmount * denominations.Nano))
				amt := amountBigInt.Mul(amountBigInt, big.NewInt(denominations.Nano))

				minSelfDelegationBigInt := big.NewInt(int64(minSelfDelegation * denominations.Nano))
				minSelfDel := minSelfDelegationBigInt.Mul(minSelfDelegationBigInt, big.NewInt(denominations.Nano))

				blsPubKey := [48]byte{}
				copy(blsPubKey[:], []byte(stakingBlsPubKey))

				return staking.DirectiveNewValidator, staking.NewValidator{
					staking.Description {
						validatorName,
						validatorIdentity,
						validatorWebsite,
						validatorSecurityContact,
						validatorDetails,
					},
					staking.CommissionRates {
						commisionRate,
						commisionMaxRate,
						commisionMaxChangeRate },
					minSelfDel,
					accounts.ParseAddrH(stakingAddress.String()),
					blsPubKey,
					amt,
				}

			}

			stakingTx, err := createStakingTransaction(getNextNonce(networkHandler), delegateStakePayloadMaker)
			if err != nil {
				return err
			}

			err = handleStakingTransaction(stakingTx, networkHandler, stakingAddress)
			if err != nil {
				return err
			}
			return nil
		},
	}


	subCmdNewValidator.Flags().StringVar(&validatorName, "name", "","validator's name")
	subCmdNewValidator.Flags().StringVar(&validatorIdentity, "identity", "", "validator's identity")
	subCmdNewValidator.Flags().StringVar(&validatorWebsite, "website", "", "validator's website")
	subCmdNewValidator.Flags().StringVar(&validatorSecurityContact, "security-contact", "","validator's security contact")
	subCmdNewValidator.Flags().StringVar(&validatorDetails, "details", "", "validator's details")
	subCmdNewValidator.Flags().StringVar(&commisionRateStr, "rate",  "","commission rate")
	subCmdNewValidator.Flags().StringVar(&commisionMaxRateStr, "max-rate","","commision max rate")
	subCmdNewValidator.Flags().StringVar(&commisionMaxChangeRateStr, "max-change-rate","","commission max change amount")
	subCmdNewValidator.Flags().Float64Var(&minSelfDelegation, "min-self-delegation", 0.0, "minimal self delegation")
	subCmdNewValidator.Flags().Var(&stakingAddress, "staking-addr", "validator's staking address")
	subCmdNewValidator.Flags().StringVar(&stakingBlsPubKey, "pubkey", "","validator's public BLS key address")
	subCmdNewValidator.Flags().Float64Var(&stakingAmount, "amount", 0.0, "staking amount")
	subCmdNewValidator.Flags().Float64Var(&gasPrice, "gas-price", 0.0, "gas price to pay")
	subCmdNewValidator.Flags().Var(&chainName, "chain-id", "what chain ID to target")
	subCmdNewValidator.Flags().StringVar(&unlockP,
		"passphrase", common.DefaultPassphrase,
		"passphrase to unlock delegator's keystore",
	)

	for _, flagName := range [...]string{"name", "identity", "website", "security-contact", "details", "rate", "max-rate",
		"max-change-rate", "min-self-delegation","staking-addr", "pubkey", "amount", } {
		subCmdNewValidator.MarkFlagRequired(flagName)
	}

	subCmdEditValidator := &cobra.Command{
		Use:   "editvalidator",
		Short: "edit a validator",
		Long: `
Edit an existing validator"
`,
		RunE: func(cmd *cobra.Command, args []string)  error {
			networkHandler, err := handlerForShard(0, node)
			if err != nil {
				return err
			}

			commisionRate, err := numeric.NewDecFromStr(commisionRateStr)
			if err != nil {
				return err
			}

			delegateStakePayloadMaker := func() (staking.Directive, interface{}) {
				minSelfDelegationBigInt := big.NewInt(int64(minSelfDelegation * denominations.Nano))
				minSelfDel := minSelfDelegationBigInt.Mul(minSelfDelegationBigInt, big.NewInt(denominations.Nano))

				return staking.DirectiveEditValidator, staking.EditValidator{
					staking.Description {
						validatorName,
						validatorIdentity,
						validatorWebsite,
						validatorSecurityContact,
						validatorDetails,
					},
					accounts.ParseAddrH(stakingAddress.String()),
					commisionRate,
					minSelfDel,
				}

			}

			stakingTx, err := createStakingTransaction(getNextNonce(networkHandler), delegateStakePayloadMaker)
			if err != nil {
				return err
			}

			err = handleStakingTransaction(stakingTx, networkHandler, stakingAddress)
			if err != nil {
				return err
			}
			return nil
		},
	}


	subCmdEditValidator.Flags().StringVar(&validatorName, "name", "","validator's name")
	subCmdEditValidator.Flags().StringVar(&validatorIdentity, "identity", "", "validator's identity")
	subCmdEditValidator.Flags().StringVar(&validatorWebsite, "website", "", "validator's website")
	subCmdEditValidator.Flags().StringVar(&validatorSecurityContact, "security-contact", "","validator's security contact")
	subCmdEditValidator.Flags().StringVar(&validatorDetails, "details", "", "validator's details")
	subCmdEditValidator.Flags().StringVar(&commisionRateStr, "rate",  "","commission rate")
	subCmdEditValidator.Flags().Float64Var(&minSelfDelegation, "min-self-delegation", 0.0, "minimal self delegation")
	subCmdEditValidator.Flags().Var(&stakingAddress, "staking-addr", "validator's staking address")
	subCmdEditValidator.Flags().Float64Var(&gasPrice, "gas-price", 0.0, "gas price to pay")
	subCmdEditValidator.Flags().Var(&chainName, "chain-id", "what chain ID to target")
	subCmdEditValidator.Flags().StringVar(&unlockP,
		"passphrase", common.DefaultPassphrase,
		"passphrase to unlock delegator's keystore",
	)

	for _, flagName := range [...]string{"name", "identity", "website", "security-contact", "details", "rate",
		"min-self-delegation","staking-addr", } {
		subCmdEditValidator.MarkFlagRequired(flagName)
	}

	subCmdDelegate := &cobra.Command{
		Use:   "delegate",
		Short: "delegate staking",
		Long: `
Delegating to a validator
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			networkHandler, err := handlerForShard(0, node)
			if err != nil {
				return err
			}

			delegateStakePayloadMaker := func() (staking.Directive, interface{}) {
				amountBigInt := big.NewInt(int64(stakingAmount * denominations.Nano))
				amt := amountBigInt.Mul(amountBigInt, big.NewInt(denominations.Nano))

				return staking.DirectiveDelegate, staking.Delegate{
					accounts.ParseAddrH(delegatorAddress.String()),
					accounts.ParseAddrH(validatorAddress.String()),
					amt,
				}
			}

			stakingTx, err := createStakingTransaction(getNextNonce(networkHandler), delegateStakePayloadMaker)
			if err != nil {
				return err
			}

			err = handleStakingTransaction(stakingTx, networkHandler, delegatorAddress)
			if err != nil {
				return err
			}
			return nil
		},
	}

	subCmdDelegate.Flags().Var(&delegatorAddress, "delegator", "delegator's address")
	subCmdDelegate.Flags().Var(&validatorAddress, "validator", "validator's address")
	subCmdDelegate.Flags().Float64Var(&stakingAmount, "amount", 0.0, "staking amount")
	subCmdDelegate.Flags().Float64Var(&gasPrice, "gas-price", 0.0, "gas price to pay")
	subCmdDelegate.Flags().Var(&chainName, "chain-id", "what chain ID to target")
	subCmdDelegate.Flags().StringVar(&unlockP,
		"passphrase", common.DefaultPassphrase,
		"passphrase to unlock delegator's keystore",
	)

	for _, flagName := range [...]string{"delegator", "validator", "amount"} {
		subCmdDelegate.MarkFlagRequired(flagName)
	}

	subCmdUnDelegate := &cobra.Command{
		Use:   "undelegate",
		Short: "un-delegate staking",
		Long: `
Remove delegating to a validator
`,
		RunE: func(cmd *cobra.Command, args []string) error  {
			networkHandler, err := handlerForShard(0, node)
			if err != nil {
				return err
			}

			delegateStakePayloadMaker := func() (staking.Directive, interface{}) {
				amountBigInt := big.NewInt(int64(stakingAmount * denominations.Nano))
				amt := amountBigInt.Mul(amountBigInt, big.NewInt(denominations.Nano))

				return staking.DirectiveUndelegate, staking.Undelegate{
					accounts.ParseAddrH(delegatorAddress.String()),
					accounts.ParseAddrH(validatorAddress.String()),
					amt,
				}
			}

			stakingTx, err := createStakingTransaction(getNextNonce(networkHandler), delegateStakePayloadMaker)
			if err != nil {
				return err
			}

			err = handleStakingTransaction(stakingTx, networkHandler, delegatorAddress)
			if err != nil {
				return err
			}
			return nil
		},
	}

	subCmdUnDelegate.Flags().Var(&delegatorAddress, "delegator", "delegator's address")
	subCmdUnDelegate.Flags().Var(&validatorAddress, "validator", "source validator's address")
	subCmdUnDelegate.Flags().Float64Var(&stakingAmount, "amount", 0.0, "staking amount")
	subCmdUnDelegate.Flags().Float64Var(&gasPrice, "gas-price", 0.0, "gas price to pay")
	subCmdUnDelegate.Flags().Var(&chainName, "chain-id", "what chain ID to target")
	subCmdUnDelegate.Flags().StringVar(&unlockP,
		"passphrase", common.DefaultPassphrase,
		"passphrase to unlock delegator's keystore",
	)

	for _, flagName := range [...]string{"delegator", "validator", "amount"} {
		subCmdUnDelegate.MarkFlagRequired(flagName)
	}

	subCmdReDelegate := &cobra.Command{
		Use:   "redelegate",
		Short: "re-delegate staking",
		Long: `
Re-delegating to a validator
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			networkHandler, err := handlerForShard(0, node)
			if err != nil {
				return err
			}

			delegateStakePayloadMaker := func() (staking.Directive, interface{}) {
				amountBigInt := big.NewInt(int64(stakingAmount * denominations.Nano))
				amt := amountBigInt.Mul(amountBigInt, big.NewInt(denominations.Nano))

				return staking.DirectiveUndelegate, staking.Redelegate{
					accounts.ParseAddrH(delegatorAddress.String()),
					accounts.ParseAddrH(validatorSrcAddress.String()),
					accounts.ParseAddrH(validatorDstAddress.String()),
					amt,
				}
			}

			stakingTx, err := createStakingTransaction(getNextNonce(networkHandler), delegateStakePayloadMaker)
			if err != nil {
				return err
			}

			err = handleStakingTransaction(stakingTx, networkHandler, delegatorAddress)
			if err != nil {
				return err
			}
			return nil
		},
	}

	subCmdReDelegate.Flags().Var(&delegatorAddress, "delegator", "delegator's address")
	subCmdReDelegate.Flags().Float64Var(&stakingAmount, "amount", 0.0, "staking amount")
	subCmdReDelegate.Flags().Var(&validatorSrcAddress, "src-validator", "source validator's address")
	subCmdReDelegate.Flags().Var(&validatorDstAddress, "dest-validator", "destination validator's address")
	subCmdReDelegate.Flags().Float64Var(&gasPrice, "gas-price", 0.0, "gas price to pay")
	subCmdReDelegate.Flags().Var(&chainName, "chain-id", "what chain ID to target")
	subCmdReDelegate.Flags().StringVar(&unlockP,
		"passphrase", common.DefaultPassphrase,
		"passphrase to unlock delegator's keystore",
	)

	for _, flagName := range [...]string{"delegator", "src-validator", "dest-validator", "amount"} {
		subCmdReDelegate.MarkFlagRequired(flagName)
	}

	return []*cobra.Command{
		subCmdNewValidator,
		subCmdEditValidator,
		subCmdDelegate,
		subCmdUnDelegate,
		subCmdReDelegate,
	}
}

func init() {
	cmdStaking := &cobra.Command{
		Use:   "staking",
		Short: "newvalidator, editvalidator, delegate, undelegate or redelegate",
		Long: `
Create a staking transaction, sign it, and send off to the Harmony blockchain
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmdStaking.AddCommand(stakingSubCommands()...)
	RootCmd.AddCommand(cmdStaking)
}
