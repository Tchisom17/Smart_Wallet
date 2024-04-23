package app

import (
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	//"github.com/tchisom17/internal/app/model"
	"github.com/tchisom17/internal/app/repository"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func SubmitEIP4337Operation() error {
	payload := strings.NewReader("{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"eth_chainId\",\"params\":[]}")

	req, err := http.NewRequest("POST", os.Getenv("BUNDLER_API_ENDPOINT"), payload)
	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bundler endpoint returned non-200 status code: %d", res.StatusCode)
	}
	//body, _ := io.ReadAll(res.Body)
	//fmt.Println(string(body))

	return nil
}

func GenerateAccount() (*bind.TransactOpts, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	chainIDconv, err := strconv.Atoi(os.Getenv("CHAIN_ID"))
	if err != nil {
		log.Fatalf("Error parsing Chain ID: %v", err)
	}

	chainID := big.NewInt(int64(chainIDconv)) // Replace 1 with the appropriate chain ID
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, err
	}

	return auth, nil
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
