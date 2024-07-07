package main

import (
	"fmt"
	"password-cracker-go/utils"
	"time"
)

func Md5Hash() {
	fmt.Println(utils.Md5("1"))
}

func FourLetterPasswordCracker(hashes []string) {
	// Loop that generates all permutations of alphanmeric characters of lenght 4
	// Each 4 letter word compare with the hashes in the list
	// If it is a hit save result in a map, else continue
	fmt.Println("Brute forcing passwords that are four letters long...")

	charSet := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	length := 4
	resultMap := make(map[string]string)

	startTime := time.Now()
	Md5PasswordCracker(charSet, []rune{}, length, hashes, resultMap)
	elapsedTime := time.Since(startTime)
	minutes := int(elapsedTime.Minutes())
	seconds := int(elapsedTime.Seconds()) % 60

	fmt.Printf("Brute Force completed, took %d minutes and %d seconds.\n", minutes, seconds)
	fmt.Printf("\nResults:\n\n")

	// Get results and print them out
	for _, hash := range hashes {
		if val, ok := resultMap[hash]; ok {
			fmt.Printf("Hash: %s ; Password: %s\n", hash, val)
		} else {
			fmt.Printf("Hash: %s, Not Found\n", hash)
		}
	}

}

// Helper function to convert combinations to permutations
func Md5PasswordCracker(chars []rune, prefix []rune, length int, hashes []string, resultMap map[string]string) {
	// Generate permutations
	if len(prefix) == length {
		password := string(prefix)
		for _, hash := range hashes {
			if utils.Md5(password) == hash {
				resultMap[hash] = password
			}
		}
		// *result = append(*result, string(prefix))
		return
	}
	for _, char := range chars {
		Md5PasswordCracker(chars, append(prefix, char), length, hashes, resultMap)
	}
}

func main() {
	hashes := []string{"7a95bf926a0333f57705aeac07a362a2", "08054846bbc9933fd0395f8be516a9f9"}
	FourLetterPasswordCracker(hashes)
}
