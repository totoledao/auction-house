package main

import (
	"fmt"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	migrations := "./internal/store/pgstore/migrations"

	cmd := exec.Command("tern", "migrate", "--migrations", migrations, "--config", "./internal/store/pgstore/migrations/tern.conf")
	outputHandler(cmd, "tern migrate")
}

// sqlc generate -f ./internal/store/pgstore/sqlc.yaml

func outputHandler(cmd *exec.Cmd, name string) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s command error. Output: ", name)
		fmt.Println(string(output))
		panic(err)
	}

	fmt.Printf("%s command success. Output: ", name)
	fmt.Println(string(output))
}
