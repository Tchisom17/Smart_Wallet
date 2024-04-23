package repository

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strings"
)

// SimpleAccountFactoryABI is the ABI of the SimpleAccountFactory contract
const SimpleAccountFactoryABI = `
[
	{
		"constant": true,
		"inputs": [
			{
				"name": "owner",
				"type": "address"
			},
			{
				"name": "salt",
				"type": "uint256"
			}
		],
		"name": "getAddress",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	}
]
`

// SimpleAccountFactory is the Go binding for the SimpleAccountFactory contract
type SimpleAccountFactory struct {
	contract *bind.BoundContract
}

// NewSimpleAccountFactory creates a new instance of the SimpleAccountFactory contract
func NewSimpleAccountFactory(address common.Address, client *ethclient.Client) *SimpleAccountFactory {
	contractAbi, err := abi.JSON(strings.NewReader(SimpleAccountFactoryABI))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	contract := bind.NewBoundContract(address, contractAbi, client, client, client)
	return &SimpleAccountFactory{contract: contract}
}

// GetCounterfactualAddress gets the counterfactual address from the SimpleAccountFactory contract
func (s *SimpleAccountFactory) GetCounterfactualAddress(opts *bind.CallOpts, userAddress common.Address, salt *big.Int) (common.Address, error) {
	var result []interface{}
	err := s.contract.Call(opts, &result, "getAddress", userAddress, salt)
	if err != nil {
		return common.Address{}, err
	}

	// Convert the first element of the result slice to a common.Address
	address, ok := result[0].(common.Address)
	if !ok {
		return common.Address{}, errors.New("unexpected type conversion")
	}

	return address, nil
}
