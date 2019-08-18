package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/linkai-io/secrets"
)

var env = os.Getenv("APP_ENV")
var region = os.Getenv("APP_REGION")

func main() {
	if env == "" {
		log.Fatalf("error reading APP_ENV did you forget to set the environment variable?")
		return
	}

	s := secrets.NewSecretsCache(env, region)
	password, err := s.Password()
	if err != nil {
		log.Fatalf("error reading password: %v\n", err)
	}

	http.HandleFunc("/", secretServer(password))
	http.ListenAndServe(":8000", nil)
}

type secretResp struct {
	Secret string `json:"secret,omitempty"`
}

func secretServer(password string) func(w http.ResponseWriter, r *http.Request) {
	secret := secretResp{Secret: password}
	data, _ := json.Marshal(secret)

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(data))
	}
}
