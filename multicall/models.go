package multicall

type (
	Config struct {
		FunctionSignature string
		ContractAddress   string
	}

	RPCRequest struct {
		To   string `json:"to"`
		Data string `json:"data"`
	}

	MulticallItem struct {
		Target   [20]byte
		CallData []byte
	}

	ViewCallResult struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
	}

	ViewCallResults map[string]ViewCallResult
)
