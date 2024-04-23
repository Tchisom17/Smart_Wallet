package accounthand

import (
	"encoding/json"
	"github.com/tchisom17/internal/app"
	"net/http"
)

func HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	account, err := app.GenerateAccount()
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

	err = app.SubmitEIP4337Operation()
	if err != nil {
		http.Error(w, "Failed to submit EIP-4337 operation", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"counterfactualAddress": counterfactualAddress}
	json.NewEncoder(w).Encode(response)
}
