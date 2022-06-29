package governance

type governanceApi string

const backendAddress = "https://snapshot.hmny.io/api/" //TODO replace by Posichain

const (
	_ governanceApi = ""

	urlListSpace                           = backendAddress + "spaces"
	urlListProposalsBySpace                = backendAddress + "%s/proposals"
	urlListProposalsVoteBySpaceAndProposal = backendAddress + "%s/proposal/%s"
	urlMessage                             = backendAddress + "message"
	urlGetValidatorsInTestNet              = "https://api.stake.hmny.io/networks/testnet/validators" //TODO replace by Posichain
	urlGetValidatorsInMainNet              = "https://api.stake.hmny.io/networks/mainnet/validators" //TODO replace by Posichain
	urlGetProposalInfo                     = "https://gateway.ipfs.io/ipfs/%s"
)
