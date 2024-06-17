package main

import (
	// "crypto/md5"
	"encoding/binary"
	"fmt"
	// "unicode/utf8"
)

func md5(plainText string) string {

	// Step 1 Append padding bits
	// get length of string in
	// append a single "1" bit to message
	// append "0" bits until message length in bits ≡ 448 (mod 512)
	// append bits even if message length is already congruent to 448 mod 512
	// get total number of bits in the message and store in variable
	// Using maths find out how many more bits need to be added to make the total number of bits congruent to 448 mod 512
	// for loop to add the bits to the message until the total number of bits is congruent to 448 mod 512

	plainTextByteSize := len(plainText)

	remainingBits := (plainTextByteSize * 8) % 512

	var bytesToAppend int

	if remainingBits < 448 {
		bytesToAppend = (448 - remainingBits) / 8
	} else if remainingBits == 448 {
		bytesToAppend = 448 / 8
	} else {
		bytesToAppend = (512 - remainingBits + 448) / 8
	}

	fmt.Println(bytesToAppend)

	// Create a list of bytes out of string
	plainTextBytes := []byte(plainText)
	plainTextBytes = append(plainTextBytes, 0x80) // append 1 to the end of the message
	bytesToAppend--                               // decrement the number of bytes to append by 1

	for i := 0; i < bytesToAppend; i++ {
		plainTextBytes = append(plainTextBytes, 0x00)
	}

	fmt.Println(plainTextBytes)
	fmt.Println(len(plainTextBytes))

	// Step 2 Append length in binary
	// append length of message in bits as 64-bit little-endian integer to message
	// the length is the number of bits in the original message
	// add the binary representation of the binary length of the original message as a 64 bit integer
	// if it exceeds 64 bits, take the least significant 64 bits, the rightmost 64 bits

	var plainTextBitLength uint64 = uint64(plainTextByteSize) * 8
	lengthBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(lengthBytes, plainTextBitLength)

	fmt.Println("Final padded: ", plainTextBytes)
	fmt.Println(len(plainTextBytes))
	fmt.Println("Length bytes: ", lengthBytes)

	paddedPlainTextBytes := append(plainTextBytes, lengthBytes...)
	fmt.Println("Padded plain text bytes: ", paddedPlainTextBytes)
	fmt.Println("Number of bytes in padded text: ", len(paddedPlainTextBytes))

	// Step 3 Initialize MD Buffer
	// Initialize 4 word 32 bit each, used for compute message digest
	// word A: 01 23 45 67
	// word B: 89 ab cd ef
	// word C: fe dc ba 98
	// word D: 76 54 32 10
	// Initialize fixed rotate constants map for each operation
	// Initialize table of round constants K

	// Initialize shift constants for each operations in each round
	s := make([]int, 64)
	rotatePatterns := [][]int{
		{7, 12, 17, 22},
		{5, 9, 14, 20},
		{4, 11, 16, 23},
		{6, 10, 15, 21},
	}

	for i, pattern := range rotatePatterns {
		for j := 0; j < 4; j++ {
			copy(s[i*16+j*4:i*16+(j+1)*4], pattern)
		}
	}

	fmt.Println(s)

	// var A int = 0x67452301
	// var B int = 0xefcdab89
	// var C int = 0x98badcfe
	// var D int = 0x10325476

	// Step 4 Process Message in 16-Word Blocks
	// for each 512-bit block of message
	// break the block into sixteen 32-bit words
	// initialize hash value for this chunk, in map M
	// initialize 4 working variables, A, B, C, D
	// main loop

	// Iterate over each 512 bit block
	for i := 0; i < len(paddedPlainTextBytes); i += 64 {
		// break the block into sixteen 32-bit words
		// initialize hash value for this chunk, in map M
		// initialize 4 working variables, A, B, C, D
		// main loop

		// Initialize hash values for this chunk
		A := 0x67452301
		B := 0xefcdab89
		C := 0x98badcfe
		D := 0x10325476

		// Initialize working variables
		var F int
		var g int

		// Initialize M
		M := make([]int, 16)

		for j := 0; j < 16; j++ {
			M[j] = int(paddedPlainTextBytes[i+j])
		}

		// Main loop
		for j := 0; j < 64; j++ {
			if j < 16 {
				F = (B & C) | ((^B) & D)
				g = j
			} else if j < 32 {
				F = (D & B) | ((^D) & C)
				g = (5*j + 1) % 16
			} else if j < 48 {
				F = B ^ C ^ D
				g = (3*j + 5) % 16
			} else {
				F = C ^ (B | (^D))
				g = (7 * j) % 16
			}

			dTemp := D
			D = C
			C = B
			B = B + (A + F + M[g] + 0x5a827999)
			B = (B << s[j]) | (B >> (32 - s[j]))
			B = B + dTemp
			A = D
		}

		// Step 4 Output
		// add this chunk's hash to result so far
		// result = result + hash of this chunk

		// Step 5 Output
		// return result

		fmt.Println(A, B, C, D)
	}

	// Step 5 Output

	return plainText
}

func main() {
	fmt.Println(md5("They are deterministic"))
	// fmt.Println(len("abcde"))
	// data := []byte("hello")
	// fmt.Printf("%x", md5.Sum(data))
}
