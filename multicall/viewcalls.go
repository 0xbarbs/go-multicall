package multicall

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"sort"
)

type ViewCalls map[string]ViewCall

func (calls ViewCalls) parseCallData(requireSuccess bool) ([]byte, error) {
	boolean, _ := abi.NewType("bool", "", nil)
	data, _ := abi.NewType("tuple[]", "", []abi.ArgumentMarshaling{
		{Type: "address", Name: "Target"},
		{Type: "bytes", Name: "CallData"},
	})

	args := abi.Arguments{
		{Type: boolean, Name: "requireSuccess"},
		{Type: data, Name: "calls"},
	}

	var test = make([]MulticallItem, 0)
	for _, call := range calls.getSortedCalls() {
		mc, err := call.parseCallData()
		if err != nil {
			return nil, err
		}

		test = append(test, mc)
	}

	return args.Pack(requireSuccess, test)
}

func (calls ViewCalls) getSortedKeys() []string {
	keys := make([]string, 0)
	for key := range calls {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

func (calls ViewCalls) getSortedCalls() []ViewCall {
	keys := calls.getSortedKeys()
	sorted := make([]ViewCall, 0)

	for _, k := range keys {
		sorted = append(sorted, calls[k])
	}

	return sorted
}
