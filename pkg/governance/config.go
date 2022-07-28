package governance

type governanceApi string

const backendAddress = "https://snapshot.hmny.io/api/" //TODO replace by Posichain

const (
	_ governanceApi = ""

	urlListSpace                           = backendAddress + "spaces"
	urlListProposalsBySpace                = backendAddress + "%s/proposals"
	urlListProposalsVoteBySpaceAndProposal = backendAddress + "%s/proposal/%s"
	urlMessage                             = backendAddress + "message"
	urlGetValidatorsInTestNet              = "https://apex-testnet.posichain.org/staking/networks/devnet/validators"
	urlGetValidatorsInMainNet              = "https://apex.posichain.org/staking/networks/devnet/validators"
	urlGetProposalInfo                     = "https://gateway.ipfs.io/ipfs/%s"
)
