package accounthand

import (
	"encoding/json"
	"fmt"
	"github.com/tchisom17/internal/app"
	"github.com/tchisom17/internal/app/model"
	"math/big"
	"net/http"
)

func HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	account, privateKey, chainID, err := app.GenerateAccount()
	if err != nil {
		http.Error(w, "Failed to generate Ethereum account", http.StatusInternalServerError)
		return
	}

	salt, err := app.GenerateRandomSalt()
	if err != nil {
		http.Error(w, "Failed to generate random salt", http.StatusInternalServerError)
		return
	}

	counterfactualAddress, err := app.GetCounterfactualAddress(account.From.Hex(), salt)
	if err != nil {
		http.Error(w, "Failed to get counterfactual address", http.StatusInternalServerError)
		return
	}

	userOp := model.UserOperation{
		Sender:               account.From,
		Nonce:                big.NewInt(0),
		CallGasLimit:         big.NewInt(100000),
		VerificationGasLimit: big.NewInt(100000),
		PreVerificationGas:   big.NewInt(100000),
		MaxFeePerGas:         big.NewInt(50000),      // in wei
		MaxPriorityFeePerGas: big.NewInt(1000000000), // in wei
	}

	signature, err := app.SignUserOperation(&userOp, privateKey, account.From, chainID)
	if err != nil {
		http.Error(w, "Failed to sign user operation", http.StatusInternalServerError)
		return
	}
	//userOp.Signature = signature
	result, err := app.SubmitEIP4337Operation(&userOp, signature)
	if err != nil {
		http.Error(w, "Failed to submit EIP-4337 operation", http.StatusInternalServerError)
		return
	}

	fmt.Println(result)
	response := map[string]string{"counterfactualAddress": counterfactualAddress}
	json.NewEncoder(w).Encode(response)
}
