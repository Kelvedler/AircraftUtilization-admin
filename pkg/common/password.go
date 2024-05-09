package common

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func PasswordFromCli() (password []byte) {
	for {
		fmt.Println("Enter password: ")
		password1, err := term.ReadPassword(0)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		passwordLen := len(string(password1))
		if passwordLen >= 20 {
			fmt.Println("Password too long (>= 20)")
			os.Exit(1)
		} else if passwordLen < 6 {
			fmt.Println("Password too short (< 6)")
			os.Exit(1)
		}
		fmt.Println("Repeat password: ")
		password2, err := term.ReadPassword(0)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if string(password1) != string(password2) {
			fmt.Println("Passwords don't match")
		} else {
			return password1
		}
	}
}
