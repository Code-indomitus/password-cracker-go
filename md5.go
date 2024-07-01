package main

import (
	// "crypto/md5"

	"encoding/binary"
	"fmt"
	"math"
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
	var A uint32 = 0x01234567
	var C uint32 = 0x89abcdef
	var B uint32 = 0xfedcba98
	var D uint32 = 0x76543210

	s := make([]uint32, 64)
	rotatePatterns := [][]uint32{
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

	T := make([]uint32, 64)
	for i := 0; i < 64; i++ {
		T[i] = uint32(math.Pow(2, 32) * math.Sin(float64(i+1)))
	}

	fmt.Println(s)

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

		// Initialize M
		X := make([]uint32, 16)

		for j := 0; j < 16; j++ {
			start := i + (j * 4)
			X[j] = uint32(bytesToInt(paddedPlainTextBytes[start : start+4]))
		}

		var AA uint32 = A
		var BB uint32 = B
		var CC uint32 = C
		var DD uint32 = D

		var k uint32
		var f auxFunc

		// Main loop
		for j := uint32(0); j < 64; j++ {

			if j < 16 {
				k = j
				f = F
			} else if j < 32 {
				k = (5*j + 1) % 16
				f = G
			} else if j < 48 {
				k = (3*j + 5) % 16
				f = H
			} else {
				k = (7 * j) % 16
				f = I
			}

			if j%4 == 0 {
				roundOperation(&A, &B, &C, &D, f, X[k], T[j], s[j])
			} else if j%4 == 1 {
				roundOperation(&D, &A, &B, &C, f, X[k], T[j], s[j])
			} else if j%4 == 2 {
				roundOperation(&C, &D, &A, &B, f, X[k], T[j], s[j])
			} else if j%4 == 3 {
				roundOperation(&B, &C, &D, &A, f, X[k], T[j], s[j])
			}
		}

		// Step 4 Output
		// add this chunk's hash to result so far
		// result = result + hash of this chunk

		A = A + AA
		B = B + BB
		C = C + CC
		D = D + DD

	}

	// Step 5 Output
	// return result
	fmt.Printf("%s %s %s %s", fmt.Sprintf("%x", A), fmt.Sprintf("%x", B), fmt.Sprintf("%x", C), fmt.Sprintf("%x", D))

	return plainText
}

/* Let [abcd k s i] denote the operation
   a = b + ((a + F(b,c,d) + X[k] + T[i]) <<< s). */
/* Do the following 16 operations for each round. */
func roundOperation(a, b, c, d *uint32, f auxFunc, k, i, s uint32) {
	*a = *b + rotateLeft((*a+f(*b, *c, *d)+k+i), s)
}

/*
 Auxiliar functions
*/

type auxFunc func(uint32, uint32, uint32) uint32

func F(x, y, z uint32) uint32 {
	return (x & y) | ((^x) & z)
}

func G(x, y, z uint32) uint32 {
	return (x & z) | (y & (^z))
}

func H(x, y, z uint32) uint32 {
	return x ^ y ^ z
}

func I(x, y, z uint32) uint32 {
	return y ^ (x | (^z))
}

/*
Helper functions
*/
func rotateLeft(x uint32, n uint32) uint32 {
	return ((x) << (n)) | ((x) >> (32 - (n)))
}

func bytesToInt(binary []byte) uint32 {
	return uint32(binary[0])<<24 | uint32(binary[1])<<16 | uint32(binary[2])<<8 | uint32(binary[3])
}

func main() {
	fmt.Println(md5("They are deterministic"))
	// fmt.Println(len("abcde"))
	// data := []byte("hello")
	// fmt.Printf("%x", md5.Sum(data))
}
