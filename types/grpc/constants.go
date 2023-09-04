package grpc

const (
	PrefixGRPC = "grpc_"

	ValidatorsInfoFileName             = "validators_info"
	DelegatorDelegationEntriesFileName = "delegator_delegation_entries"
	DelegatorSharesFileName            = "delegator_shares"
	GenesisAccountAnalysisFileName     = "genesis_accounts"
	GenesisPostGenesisValidators       = "genesis_post_genesis_validators"
	DelegationAnalysis                 = "delegation_analysis"

	HeaderMoniker              = "Moniker"
	HeaderPercentageWeight     = "Percentage Weight"
	HeaderSelfDelegation       = "Self Delegation"
	HeaderTotalDelegations     = "Total Delegations"
	HeaderDelegator            = "Delegator"
	HeaderValidator            = "Validator"
	HeaderShares               = "Shares"
	HeaderAddress              = "Address"
	HeaderVestingEndTime       = "Vesting End Time"
	HeaderOriginalVesting      = "Original Vesting"
	HeaderVestingStartTime     = "Vesting Start Time"
	HeaderTokensFreeEveryBlock = "Tokens Free Every Block"
	HeaderTokensFreeEveryDay   = "Tokens Free Every Day"

	HeaderGenesisType       = "Genesis Type"
	HeaderOperatorAddress   = "Operator Address"
	HeaderConsensusPubkey   = "Consensus Pubkey"
	HeaderStatus            = "Status"
	HeaderTokens            = "Tokens"
	HeaderDelegatorShares   = "Delegator Shares"
	HeaderDescription       = "Description"
	HeaderUnbondingHeight   = "Unbonding Height"
	HeaderUnbondingTime     = "Unbonding Time"
	HeaderCommission        = "Commission"
	HeaderMinSelfDelegation = "Min Self Delegation"
	HeaderJailed            = "Jailed"

	Message           = "message"
	Wasm              = "wasm"
	BondID            = "bond_id"
	CoinSpent         = "coin_spent"
	CoinReceived      = "coin_received"
	ContractAddress   = "_contract_address"
	LockID            = "lock_id"
	LockedTokens      = "locked_tokens"
	Action            = "action"
	CallbackInfo      = "callback-info"
	ReplyMsgID        = "reply-msg-id"
	ReplyResult       = "reply-result"
	User              = "user"
	VaultTokenBalance = "vault_token_balance"
	Websocket         = "/websocket"

	VaultAddress      = "quasar18a2u6az6dzw528rptepfg6n49ak6hdzkf8ewf0n5r0nwju7gtdgqamr7qu"
	PrimitiveAddress1 = "quasar1kj8q8g2pmhnagmfepp9jh9g2mda7gzd0m5zdq0s08ulvac8ck4dq9ykfps"
	PrimitiveAddress2 = "quasar1ma0g752dl0yujasnfs9yrk6uew7d0a2zrgvg62cfnlfftu2y0egqx8e7fv"
	PrimitiveAddress3 = "quasar1ery8l6jquynn9a4cz2pff6khg8c68f7urt33l5n9dng2cwzz4c4qxhm6a2"

	IdentifierDelayedVestingAccount    = "/cosmos.vesting.v1beta1.DelayedVestingAccount"
	IdentifierContinuousVestingAccount = "/cosmos.vesting.v1beta1.ContinuousVestingAccount"
	IdentifierPermanentLockedAccount   = "/cosmos.vesting.v1beta1.PermanentLockedAccount"
	IdentifierPeriodicVestingAccount   = "/cosmos.vesting.v1beta1.PeriodicVestingAccount"
	IdentifierMsgExecuteContract       = "/cosmwasm.wasm.v1.MsgExecuteContract"
	IdentifierMsgUpdateClient          = "/ibc.core.client.v1.MsgUpdateClient"
	IdentifierMsgAcknowledgement       = "/ibc.core.channel.v1.MsgAcknowledgement"

	ValidatorsLimit  = 50000
	DelegatorsLimit  = 1000000000
	AverageBlockTime = 5
	SecondsInADay    = 86400
)
