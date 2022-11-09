package multicall

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/rpc"
	"reflect"
	"strings"
)

type Contract struct {
	client *rpc.Client
	config *Config
}

func NewContract(client *rpc.Client, config Config) Contract {
	return Contract{
		client: client,
		config: &config,
	}
}

func (c *Contract) Call(calls ViewCalls) (ViewCallResults, error) {
	callData, err := calls.parseCallData()
	if err != nil {
		return nil, err
	}

	result, err := c.process(callData)
	if err != nil {
		return nil, err
	}

	return c.decodeResult(result, calls)
}

func (c *Contract) process(callData []byte) (result string, err error) {
	req := RPCRequest{
		To:   c.config.ContractAddress,
		Data: c.config.FunctionSignature + hex.EncodeToString(callData),
	}

	err = c.client.Call(&result, "eth_call", req, "latest")
	if err != nil {
		return
	}

	return
}

func (c *Contract) decodeResult(result string, calls ViewCalls) (ViewCallResults, error) {
	results := make(ViewCallResults, 0)
	sortedKeys := calls.getSortedKeys()
	sortedCalls := calls.getSortedCalls()
	rawBytes, err := hex.DecodeString(strings.Replace(result, "0x", "", -1))
	if err != nil {
		return nil, err
	}

	resultType, _ := abi.NewType("tuple[]", "", []abi.ArgumentMarshaling{
		{Name: "Success", Type: "bool"},
		{Name: "Data", Type: "bytes"},
	})

	resultArgs := abi.Arguments{
		{
			Name: "ResultItem",
			Type: resultType,
		},
	}

	data, err := resultArgs.Unpack(rawBytes)
	if err != nil {
		return nil, err
	}

	reflected := reflect.ValueOf(data[0])
	for i := 0; i < reflected.Len(); i++ {
		item := reflected.Index(i)
		success := item.FieldByName("Success").Bool()
		bytes := item.FieldByName("Data").Bytes()
		call := sortedCalls[i]
		value, err := c.convertToReturnType(bytes, call)

		if err != nil {
			return nil, err
		}

		results[sortedKeys[i]] = ViewCallResult{
			Success: success,
			Data:    value,
		}
	}

	return results, nil
}

func (c *Contract) convertToReturnType(input []byte, call ViewCall) (interface{}, error) {
	args, err := call.getReturnAbiArguments()
	if err != nil {
		return nil, err
	}

	value, err := args.Unpack(input)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack args when converting return type: %s", err)
	}

	if len(value) == 1 {
		return value[0], nil
	}

	return value, nil
}
