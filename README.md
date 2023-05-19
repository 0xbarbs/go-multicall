# go-multicall

A golang utility library to help with all things multi-call.

The current implementation is built to interact with the [MakerDAO Multicall2 contracts](https://github.com/makerdao/multicall#multicall2-contract-addresses) but future versions plan to support custom integrations.

## Features:

- [x] Support multiple parameters/return types
- [x] Cast to provided return types
- [x] Create key value map for multi-calls
- [x] Use MakerDAO tryAggregate() function to allow calls to fail
- [x] Support for tuple[] parameters
- [ ] Custom contract implementations via go interfaces
- [ ] Alternative ViewCallResult structures

---

### Union User Manager example:

```go
client, _ := rpc.DialHTTP("https://mainnet.infura.io/v3/{API_KEY}")
address := common.HexToAddress("0xb8150a1b6945e75d05769d685b127b41e6335bbc")

contract := multicall.NewContract(client, multicall.Config{
    ContractAddress:   "0x5ba1e12693dc8f9c48aad8770482f4739beed696", // mainnet multicall
    FunctionSignature: "0xbce38bd7",
})

calls := make(multicall.ViewCalls, 0)
calls["is_member"] = multicall.NewViewCall(
    "0x49c910Ba694789B58F53BFF80633f90B8631c195", // Union mainnet user manager
    "checkIsMember(address)(bool)",
    address,
)
calls["credit_limit"] = multicall.NewViewCall(
    "0x49c910Ba694789B58F53BFF80633f90B8631c195", // Union mainnet user manager
    "getCreditLimit(address)(int256)",
    address,
)
calls["staker_addresses"] = multi.NewViewCall(
    "0x49c910Ba694789B58F53BFF80633f90B8631c195", // Union mainnet user manager
    "getStakerAddresses(address)(address[])",
    address,
)

result, _ := contract.Call(calls)
```

The above multi-call will return:

```
{
    "credit_limit": {
        "success": true,
        "data": 162734905059436317281
    },
    "is_member": {
        "success": true,
        "data": true
    },
    "staker_addresses": {
        "success": true,
        "data": [
            "0x497c20fed24d61c7506ef2500065e4fd662f3779",
            "0x258bad0751299ce659e443db1d166fd881cba281",
            "0x7a0c61edd8b5c0c5c1437aeb571d7ddbf8022be4",
            "0x230d31eec85f4063a405b0f95bde509c0d0a8b5d"
        ]
    }
}
```