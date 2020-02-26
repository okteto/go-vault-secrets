package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/vault/api"
)

func getHardcodedSecrets() (string, string, error) {
	return "hard-coded-user", "hard-coded-password", nil
}

func getVaultSecrets() (string, string, error) {
	client, err := api.NewClient(&api.Config{
		Address: "http://vault",
	})

	if err != nil {
		return "", "", err
	}

	client.SetToken(os.Getenv("VAULT_TOKEN"))

	secretName := "secret/data/helloworld"
	secret, err := client.Logical().Read(secretName)
	if err != nil {
		return "", "", err
	}

	if secret == nil {
		return "", "", fmt.Errorf("%s not found", secretName)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return "", "", fmt.Errorf("couldn't parse the response")
	}

	u := data["username"].(string)
	p := data["password"].(string)

	return u, p, nil
}

func main() {
	fmt.Println("Starting hello-world server...")
	http.HandleFunc("/", helloServer)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func helloServer(w http.ResponseWriter, r *http.Request) {
	u, p, err := getHardcodedSecrets()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		w.WriteHeader(500)
	}

	fmt.Fprint(w, "Hello world!\n")
	fmt.Fprintf(w, "Super secret username: %s\n", u)
	fmt.Fprintf(w, "Super secret password is: %s\n", p)
}
