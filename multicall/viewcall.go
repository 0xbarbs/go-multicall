package multicall

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"regexp"
	"strings"
)

type ViewCall struct {
	target string
	method string
	params []interface{}
}

func NewViewCall(target string, method string, params ...interface{}) ViewCall {
	return ViewCall{
		target: target,
		method: method,
		params: params,
	}
}

func (vc *ViewCall) parseCallData() (call MulticallItem, err error) {
	types, err := vc.getParameterTypes()
	if err != nil {
		return
	}
	args, err := vc.getAbiArguments(types)
	if err != nil {
		return
	}

	b, err := args.Pack(vc.params...)
	if err != nil {
		return call, fmt.Errorf("failed to pack args when parsing call data: %s", err)
	}

	data := append(vc.getMethodSignature(), b...)
	call = MulticallItem{
		Target:   common.HexToAddress(vc.target),
		CallData: data,
	}

	return
}

func (vc *ViewCall) getAbiArguments(types []string) (args abi.Arguments, err error) {
	for _, t := range types {
		abiType, err := abi.NewType(t, "", nil)
		if err != nil {
			return args, errors.New(fmt.Sprintf("%s is not a valid parameter type", t))
		}

		args = append(args, abi.Argument{
			Type: abiType,
		})
	}

	return
}

func (vc *ViewCall) getParameterTypes() ([]string, error) {
	r := regexp.MustCompile("\\((.*?)\\)")
	matches := r.FindAllStringSubmatch(vc.method, -1)

	if len(matches) < 2 {
		return nil, errors.New(fmt.Sprintf("%s is not a valid method signature", vc.method))
	}

	return strings.Split(matches[0][1], ","), nil
}

func (vc *ViewCall) getParameterAbiArguments() (abi.Arguments, error) {
	types, err := vc.getParameterTypes()
	if err != nil {
		return nil, err
	}

	return vc.getAbiArguments(types)
}

func (vc *ViewCall) getReturnTypes() ([]string, error) {
	r := regexp.MustCompile("\\((.*?)\\)")
	matches := r.FindAllStringSubmatch(vc.method, -1)

	if len(matches) < 2 {
		return nil, errors.New(fmt.Sprintf("%s is not a valid method signature", vc.method))
	}

	return strings.Split(matches[1][1], ","), nil
}

func (vc *ViewCall) getReturnAbiArguments() (abi.Arguments, error) {
	types, err := vc.getReturnTypes()
	if err != nil {
		return nil, err
	}

	return vc.getAbiArguments(types)
}

func (vc *ViewCall) getMethodSignature() []byte {
	r := regexp.MustCompile("^(.*?)\\)")
	sig := r.FindString(vc.method) // todo: validate?
	hash := crypto.Keccak256([]byte(sig))
	return hash[0:4]
}
