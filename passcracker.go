package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"password-cracker-go/utils"
	"strings"
	"time"
)

func Md5Hash() {
	fmt.Println(utils.Md5("1"))
}

func FourLetterPasswordCracker(hashes []string) {
	// Loop that generates all permutations of alphanmeric characters of lenght 4
	// Each 4 letter word compare with the hashes in the list
	// If it is a hit save result in a map, else continue
	fmt.Println("Brute forcing passwords that are four characters long...")

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
			fmt.Printf("Hash: %s -> Password: %s\n", hash, val)
		} else {
			fmt.Printf("Hash: %s, Not Found\n", hash)
		}
	}

}

func Md5PasswordCracker(chars []rune, prefix []rune, length int, hashes []string, resultMap map[string]string) {
	// Generate permutations
	if len(prefix) == length {
		password := string(prefix)
		for _, hash := range hashes {
			if utils.Md5(password) == hash {
				resultMap[hash] = password
			}
		}
		return
	}
	for _, char := range chars {
		Md5PasswordCracker(chars, append(prefix, char), length, hashes, resultMap)
	}
}

func Md5WordListCracker(hash string, wordListFilePath string) {
	// Open the file that has been passed as an argument
	file, err := os.Open(wordListFilePath)
	if err != nil {
		fmt.Println("passcracker:", wordListFilePath+":", "No such file or directory")
		return
	}

	// Close the file later
	defer file.Close()

	// Create a scanner to read the file
	scanner := bufio.NewScanner(file)

	fmt.Println("Checking word list for password hash", hash, "...")

	startTime := time.Now()

	lineCount := 0

	// Read each line and strip newline characters
	for scanner.Scan() {
		lineCount++
		password := strings.TrimSpace(scanner.Text())
		if utils.Md5(password) == hash {
			fmt.Println("Hash found in word list:")
			fmt.Printf("Hash: %s -> Password: %s\n", hash, password)
			return
		}
	}
	fmt.Println("Hash NOT found in word list.")
	fmt.Printf("\nNumber of Words in list: %d\n", lineCount)

	elapsedTime := time.Since(startTime)
	minutes := int(elapsedTime.Minutes())
	seconds := int(elapsedTime.Seconds()) % 60
	fmt.Printf("Word list check completed, took %d minutes and %d seconds.\n", minutes, seconds)
}

func main() {
	// Declare the flag options of password cracker program
	var wordListFlag bool

	// Parse command line arguments for flags
	flag.BoolVar(&wordListFlag, "l", false, "use a word list to crack a password hash")
	flag.Parse()

	if len(os.Args) <= 1 {
		fmt.Println("Error: Missing parameters required for password cracking.")
		return
	}

	if wordListFlag {
		// Get the hash to be cracked
		hash := os.Args[2]
		fmt.Println(hash)
		// Get the file name to be read and start cracking
		filePath := os.Args[len(os.Args)-1]
		Md5WordListCracker(hash, filePath)
	} else {
		// Get the hash to be cracked
		hash := os.Args[1]
		fmt.Println(hash)
		FourLetterPasswordCracker([]string{hash})
	}
}

// hashes := []string{"7a95bf926a0333f57705aeac07a362a2", "08054846bbc9933fd0395f8be516a9f9"} PASS and CODE
