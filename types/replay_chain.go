package types

type DepositorDetailsBond struct {
	Address      string `json:"address"`
	BlockHeight  int64  `json:"block_height"`
	Amount       string `json:"amount"`
	VaultAddress string `json:"primitive_address"`
	BondID       int64  `json:"bond_id"`
}

type DepositorDetailsUnbond struct {
	Address      string `json:"address"`
	BlockHeight  int64  `json:"block_height"`
	VaultAddress string `json:"primitive_address"`
	BurntShares  string `json:"burnt_shares"`
	UnbondID     int64  `json:"unbond_id"`
}

type LockDetailsByHeight struct {
	Height          int64             `json:"height"`
	ContractDetails []ContractDetails `json:"contract_details"`
}

type ContractDetails struct {
	Address                 string `json:"address"`
	LockID                  int64  `json:"lock_id"`
	LockedTokensProtoString string `json:"locked_tokens_proto_string"`
	Action                  string `json:"action"`
	CallbackInfo            string `json:"callback_info"`
	ReplyMessageID          string `json:"reply_message_id"`
	ReplyResult             string `json:"reply_result"`
}

type Test struct {
	Address           string   `json:"address"`
	Shares            []string `json:"shares"`
	LastUpdatedHeight []int64  `json:"last_updated_height"`
}

type AddressSharesInIncentiveContract struct {
	Shares            []string `json:"shares"`
	LastUpdatedHeight []int64  `json:"last_updated_height"`
}

type CallBackInfoWithHeight struct {
	Height        int64          `json:"height"`
	CallBackInfos []CallBackInfo `json:"callBackInfos"`
}

type CallBackInfo struct {
	ContractAddress    string `json:"contract_address"`
	Action             string `json:"action"`
	CallBackInfoString string `json:"call_back_info"`
	ReplyMsgID         string `json:"reply_msg_id"`
	ReplyResult        string `json:"reply_result"`
}

type BeginUnlocking struct {
	Height          int64  `json:"height"`
	Step            string `json:"step"`
	ContractAddress string `json:"contract_address"`
	PendingMsg      string `json:"pending-msg"`
}

type ToBeMintedBondsAtHeight struct {
	Height               int64                  `json:"height"`
	DepositorDetailsBond []DepositorDetailsBond `json:"depositor_details_bond"`
}
type DirRange []int64

func (a DirRange) Len() int           { return len(a) }
func (a DirRange) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DirRange) Less(i, j int) bool { return a[i] < a[j] }
