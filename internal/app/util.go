package app

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tchisom17/internal/app/model"
	"io"

	//"github.com/tchisom17/internal/app/model"
	"github.com/tchisom17/internal/app/repository"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
)

func SubmitEIP4337Operation(userOp *model.UserOperation, signature []byte) (string, error) {
	//userOp.InitCode = common.FromHex("0x")
	userOp.CallData = common.FromHex("0x")
	userOp.PaymasterAndData = common.FromHex("0x")
	userOp.Signature = signature
	payload, err := userOp.MarshalJSON()
	//payload, err := json.Marshal(map[string]interface{}{
	//	"jsonrpc": "2.0",
	//	"id":      1,
	//	"method":  "eth_sendUserOperation",
	//	"params": []interface{}{
	//		userOp,
	//		"0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789",
	//	},
	//})
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", os.Getenv("BUNDLER_API_ENDPOINT"), bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(body), nil
}

func GenerateAccount() (*bind.TransactOpts, *ecdsa.PrivateKey, *big.Int, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, nil, nil, err
	}

	chainIDconv, err := strconv.Atoi(os.Getenv("CHAIN_ID"))
	if err != nil {
		log.Fatalf("Error parsing Chain ID: %v", err)
	}

	chainID := big.NewInt(int64(chainIDconv)) // Replace 1 with the appropriate chain ID
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, nil, nil, err
	}

	return auth, privateKey, chainID, nil
}

//func SignUserOperation(op model.UserOperation, privateKey *ecdsa.PrivateKey) ([]byte, error) {
//	// Serialize the user operation struct
//	data, err := json.Marshal(op)
//	if err != nil {
//		return nil, err
//	}
//
//	// Hash the serialized data
//	hash := crypto.Keccak256Hash(data)
//
//	// Sign the hash with the account's private key
//	signature, err := crypto.Sign(hash.Bytes(), privateKey)
//	if err != nil {
//		return nil, err
//	}
//
//	return signature, nil
//}

func SignUserOperation(op *model.UserOperation, privateKey *ecdsa.PrivateKey, entryPointAddress common.Address, chainID *big.Int) ([]byte, error) {
	// Hash the user operation
	hash := op.GetUserOpHash(entryPointAddress, chainID)

	// Sign the hash
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func GenerateRandomSalt() (*big.Int, error) {
	saltBytes := make([]byte, 32)
	_, err := rand.Read(saltBytes)
	if err != nil {
		return nil, err
	}

	salt := new(big.Int).SetBytes(saltBytes)
	return salt, nil
}

func GetCounterfactualAddress(address string, salt *big.Int) (string, error) {
	client, err := ethclient.Dial(os.Getenv("TESTNET_NODE_URL"))
	if err != nil {
		return "", err
	}
	defer client.Close()

	contractAddress := common.HexToAddress(os.Getenv("SIMPLE_ACCOUNT_FACTORY_ADDRESS"))
	instance := repository.NewSimpleAccountFactory(contractAddress, client)

	counterfactualAddress, err := instance.GetCounterfactualAddress(nil, common.HexToAddress(address), salt)
	if err != nil {
		return "", err
	}

	return counterfactualAddress.Hex(), nil
}
