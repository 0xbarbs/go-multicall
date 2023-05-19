package multicall_test

import (
	"github.com/0xbarbs/go-multicall/multicall"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"testing"
)

const (
	MULTICALL_CONTRACT         = "0x5ba1e12693dc8f9c48aad8770482f4739beed696"
	TRY_AGGREGATE_SIG          = "0xbce38bd7"
	UNION_USERMANAGER_CONTRACT = "0x49c910Ba694789B58F53BFF80633f90B8631c195"
)

func TestContract_Call(t *testing.T) {
	client, err := rpc.DialHTTP("https://mainnet.infura.io/v3/")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	address := common.HexToAddress("0xb8150a1b6945e75d05769d685b127b41e6335bbc")
	contract := multicall.NewContract(client, multicall.Config{
		ContractAddress:   MULTICALL_CONTRACT,
		FunctionSignature: TRY_AGGREGATE_SIG,
	})

	calls := make(multicall.ViewCalls, 0)
	calls["is_member"] = multicall.NewViewCall(
		UNION_USERMANAGER_CONTRACT,
		"checkIsMember(address)(bool)",
		address,
	)

	result, err := contract.Call(calls, true)
	if err != nil {
		log.Fatal(err)
	}

	if result["is_member"].Data != true {
		t.Errorf("is_member: expected %t, got %t", true, result["is_member"].Data)
	}
}
