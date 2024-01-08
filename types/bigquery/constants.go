package bigquery

const (
	PrefixBigQuery = "bigquery_"

	PrependQueryRaw          = "raw"
	PrependQueryTransactions = "transactions_"

	PrependLPQueryVaultsBond   = "lp_vaults_bond"
	PrependLPQueryVaultsUnbond = "lp_vaults_unbond"
	PrependLPQueryVaultsClaim  = "lp_vaults_claim"
	PrependLPQueryReport       = "lp_vaults_report"

	PrependCLQueryVaultsDeposit           = "cl_vaults_deposit"
	PrependCLQueryVaultsWithdraw          = "cl_vaults_withdraw"
	PrependCLQueryVaultsClaim             = "cl_vaults_claim"
	PrependCLQueryVaultsDistributeRewards = "cl_vaults_distribute_rewards"
	PrependCLQueryVaultsAPR               = "cl_vaults_apr"
	PrependCLQueryReport                  = "cl_vaults_report"
)
