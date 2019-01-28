package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

// this command will try to display the value at the path passed as the first argument
func main() {
	flag.Parse()
	args := flag.Args()

	client, err := api.NewClient(&api.Config{
		Address: "http://localhost:8200",
	})
	if err != nil {
		panic(err)
	}

	if err := login(client); err != nil {
		panic(err)
	}

	value, err := getSecretValue(client, args[0])
	if err != nil {
		panic(err)
	}
	fmt.Printf("Secret: %+v\n", value)
}

// login returns a Vault token based on the JWT in NETLIFY_SECRET_TOKEN
func login(client *api.Client) error {
	jwt, exists := os.LookupEnv("NETLIFY_SECRET_TOKEN")
	if !exists {
		return errors.New("Login token missing. Please set NETLIFY_SECRET_TOKEN")
	}

	data, err := client.Logical().Write("auth/jwt/login", map[string]interface{}{
		"jwt":  jwt,
		"role": "demo",
	})
	if err != nil {
		return err
	}

	client.SetToken(data.Auth.ClientToken)

	return nil
}

// getSecretValue returns the value of a secret at the given path
func getSecretValue(client *api.Client, path string) (string, error) {
	secret, err := client.Logical().Read(fmt.Sprintf("secret/data/%s", path))
	if err != nil {
		return "", err
	}
	secretData := secret.Data["data"]
	var value string
	if dataMap, ok := secretData.(map[string]interface{}); ok {
		value = dataMap["value"].(string)
	} else {
		return "", errors.New("Invalid data structure")
	}

	return value, nil
}
