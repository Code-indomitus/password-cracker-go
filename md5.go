package main

import (
	// "crypto/md5"

	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"math/bits"
	// "unicode/utf8"
)

func md6(plainText string) string {

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
	binary.LittleEndian.PutUint64(lengthBytes, plainTextBitLength)

	fmt.Println("Final padded: ", plainTextBytes)
	fmt.Println(len(plainTextBytes))
	fmt.Println("Length bytes: ", lengthBytes)

	paddedPlainTextBytes := append(plainTextBytes, lengthBytes...)
	fmt.Println("Padded plain text bytes: ", paddedPlainTextBytes)

	for i := range paddedPlainTextBytes {
		printByteAsBits(paddedPlainTextBytes[i])
	}
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
	// var A uint32 = 0x01234567
	// var C uint32 = 0x89abcdef
	// var B uint32 = 0xfedcba98
	// var D uint32 = 0x76543210
	var A uint32 = 0x67452301
	var B uint32 = 0xEFCDAB89
	var C uint32 = 0x98BADCFE
	var D uint32 = 0x10325476

	fmt.Println("A: ", A)

	fmt.Println("First here let us go yay: ")
	fmt.Printf("%s %s %s %s", fmt.Sprintf("%x", A), fmt.Sprintf("%x", B), fmt.Sprintf("%x", C), fmt.Sprintf("%x", D))

	// s := make([]uint32, 64)
	// rotatePatterns := [][]uint32{
	// 	{7, 12, 17, 22},
	// 	{5, 9, 14, 20},
	// 	{4, 11, 16, 23},
	// 	{6, 10, 15, 21},
	// }

	// for i, pattern := range rotatePatterns {
	// 	for j := 0; j < 4; j++ {
	// 		copy(s[i*16+j*4:i*16+(j+1)*4], pattern)
	// 	}
	// }

	// fmt.Println("Rotations: ", s)

	T := make([]uint32, 64)
	// Z := make([]int, 64)

	for i := 0; i < 64; i++ {
		T[i] = uint32(math.Pow(2, 32) * math.Abs(math.Sin(float64(i+1))))
		// Z[i] = int(math.Pow(2, 32) * math.Abs(math.Sin(float64(i+1))))

		fmt.Println("Table Values: ")

		fmt.Printf(fmt.Sprintf("T[%d] : %x\n", i, T[uint32(i)]))

	}

	// fmt.Println("T here hello: ", T[uint32(13)])
	// fmt.Println(Z[0])

	// fmt.Println(s)

	// Step 4 Process Message in 16-Word Blocks
	// for each 512-bit block of message
	// break the block into sixteen 32-bit words
	// initialize hash value for this chunk, in map M
	// initialize 4 working variables, A, B, C, D
	// main loop
	count := 0

	// Iterate over each 512 bit block
	for i := 0; i <= len(paddedPlainTextBytes)-64; i += 64 {
		count += 1
		fmt.Println("Count: ", count)
		// break the block into sixteen 32-bit words
		// initialize hash value for this chunk, in map M
		// initialize 4 working variables, A, B, C, D
		// main loop

		// Initialize M
		X := make([]uint32, 16)
		fmt.Println("Start")
		for j := 0; j < 16; j++ {
			start := i + (j * 4)
			X[j] = binary.LittleEndian.Uint32(paddedPlainTextBytes[start : start+4])
			// bytesToInt(paddedPlainTextBytes[start : start+4])
			fmt.Println("Values:", j)
			printIntAsBits(X[j])
		}
		fmt.Println("End")

		var AA uint32 = A
		var BB uint32 = B
		var CC uint32 = C
		var DD uint32 = D

		// Before first round:
		fmt.Println("Before first round:")
		fmt.Printf("%s %s %s %s\n", fmt.Sprintf("%x", A), fmt.Sprintf("%x", B), fmt.Sprintf("%x", C), fmt.Sprintf("%x", D))

		// Round 1
		A = B + bits.RotateLeft32(F(B, C, D)+A+X[0]+T[0], 7) // [ABCD  0  7  1]

		// After first round:
		fmt.Println("After first round:")
		fmt.Printf("%s %s %s %s\n", fmt.Sprintf("%x", A), fmt.Sprintf("%x", B), fmt.Sprintf("%x", C), fmt.Sprintf("%x", D))

		D = A + bits.RotateLeft32(F(A, B, C)+D+X[1]+T[1], 12)
		C = D + bits.RotateLeft32(F(D, A, B)+C+X[2]+T[2], 17)
		B = C + bits.RotateLeft32(F(C, D, A)+B+X[3]+T[3], 22)
		A = B + bits.RotateLeft32(F(B, C, D)+A+X[4]+T[4], 7)
		D = A + bits.RotateLeft32(F(A, B, C)+D+X[5]+T[5], 12)
		C = D + bits.RotateLeft32(F(D, A, B)+C+X[6]+T[6], 17)
		B = C + bits.RotateLeft32(F(C, D, A)+B+X[7]+T[7], 22)
		A = B + bits.RotateLeft32(F(B, C, D)+A+X[8]+T[8], 7)
		D = A + bits.RotateLeft32(F(A, B, C)+D+X[9]+T[9], 12)
		C = D + bits.RotateLeft32(F(D, A, B)+C+X[10]+T[10], 17)
		B = C + bits.RotateLeft32(F(C, D, A)+B+X[11]+T[11], 22)
		A = B + bits.RotateLeft32(F(B, C, D)+A+X[12]+T[12], 7)
		D = A + bits.RotateLeft32(F(A, B, C)+D+X[13]+T[13], 12)
		C = D + bits.RotateLeft32(F(D, A, B)+C+X[14]+T[14], 17)
		B = C + bits.RotateLeft32(F(C, D, A)+B+X[15]+T[15], 22)

		// Round 2
		A = B + bits.RotateLeft32(G(B, C, D)+A+X[1]+T[16], 5)
		D = A + bits.RotateLeft32(G(A, B, C)+D+X[6]+T[17], 9)
		C = D + bits.RotateLeft32(G(D, A, B)+C+X[11]+T[18], 14)
		B = C + bits.RotateLeft32(G(C, D, A)+B+X[0]+T[19], 20)
		A = B + bits.RotateLeft32(G(B, C, D)+A+X[5]+T[20], 5)
		D = A + bits.RotateLeft32(G(A, B, C)+D+X[10]+T[21], 9)
		C = D + bits.RotateLeft32(G(D, A, B)+C+X[15]+T[22], 14)
		B = C + bits.RotateLeft32(G(C, D, A)+B+X[4]+T[23], 20)
		A = B + bits.RotateLeft32(G(B, C, D)+A+X[9]+T[24], 5)
		D = A + bits.RotateLeft32(G(A, B, C)+D+X[14]+T[25], 9)
		C = D + bits.RotateLeft32(G(D, A, B)+C+X[3]+T[26], 14)
		B = C + bits.RotateLeft32(G(C, D, A)+B+X[8]+T[27], 20)
		A = B + bits.RotateLeft32(G(B, C, D)+A+X[13]+T[28], 5)
		D = A + bits.RotateLeft32(G(A, B, C)+D+X[2]+T[29], 9)
		C = D + bits.RotateLeft32(G(D, A, B)+C+X[7]+T[30], 14)
		B = C + bits.RotateLeft32(G(C, D, A)+B+X[12]+T[31], 20)

		// Round 3
		A = B + bits.RotateLeft32(H(B, C, D)+A+X[5]+T[32], 4)
		D = A + bits.RotateLeft32(H(A, B, C)+D+X[8]+T[33], 11)
		C = D + bits.RotateLeft32(H(D, A, B)+C+X[11]+T[34], 16)
		B = C + bits.RotateLeft32(H(C, D, A)+B+X[14]+T[35], 23)
		A = B + bits.RotateLeft32(H(B, C, D)+A+X[1]+T[36], 4)
		D = A + bits.RotateLeft32(H(A, B, C)+D+X[4]+T[37], 11)
		C = D + bits.RotateLeft32(H(D, A, B)+C+X[7]+T[38], 16)
		B = C + bits.RotateLeft32(H(C, D, A)+B+X[10]+T[39], 23)
		A = B + bits.RotateLeft32(H(B, C, D)+A+X[13]+T[40], 4)
		D = A + bits.RotateLeft32(H(A, B, C)+D+X[0]+T[41], 11)
		C = D + bits.RotateLeft32(H(D, A, B)+C+X[3]+T[42], 16)
		B = C + bits.RotateLeft32(H(C, D, A)+B+X[6]+T[43], 23)
		A = B + bits.RotateLeft32(H(B, C, D)+A+X[9]+T[44], 4)
		D = A + bits.RotateLeft32(H(A, B, C)+D+X[12]+T[45], 11)
		C = D + bits.RotateLeft32(H(D, A, B)+C+X[15]+T[46], 16)
		B = C + bits.RotateLeft32(H(C, D, A)+B+X[2]+T[47], 23)

		// Round 4
		A = B + bits.RotateLeft32(I(B, C, D)+A+X[0]+T[48], 6)
		D = A + bits.RotateLeft32(I(A, B, C)+D+X[7]+T[49], 10)
		C = D + bits.RotateLeft32(I(D, A, B)+C+X[14]+T[50], 15)
		B = C + bits.RotateLeft32(I(C, D, A)+B+X[5]+T[51], 21)
		A = B + bits.RotateLeft32(I(B, C, D)+A+X[12]+T[52], 6)
		D = A + bits.RotateLeft32(I(A, B, C)+D+X[3]+T[53], 10)
		C = D + bits.RotateLeft32(I(D, A, B)+C+X[10]+T[54], 15)
		B = C + bits.RotateLeft32(I(C, D, A)+B+X[1]+T[55], 21)
		A = B + bits.RotateLeft32(I(B, C, D)+A+X[8]+T[56], 6)
		D = A + bits.RotateLeft32(I(A, B, C)+D+X[15]+T[57], 10)
		C = D + bits.RotateLeft32(I(D, A, B)+C+X[6]+T[58], 15)
		B = C + bits.RotateLeft32(I(C, D, A)+B+X[13]+T[59], 21)
		A = B + bits.RotateLeft32(I(B, C, D)+A+X[4]+T[60], 6)
		D = A + bits.RotateLeft32(I(A, B, C)+D+X[11]+T[61], 10)
		C = D + bits.RotateLeft32(I(D, A, B)+C+X[2]+T[62], 15)
		B = C + bits.RotateLeft32(I(C, D, A)+B+X[9]+T[63], 21)

		A = A + AA
		B = B + BB
		C = C + CC
		D = D + DD

		// fmt.Printf("%s %s %s %s", fmt.Sprintf("%x", A), fmt.Sprintf("%x", B), fmt.Sprintf("%x", C), fmt.Sprintf("%x", D))

		// Step 4 Output
		// add this chunk's hash to result so far
		// result = result + hash of this chunk
		fmt.Println("Hare Krishna")
	}

	fmt.Println("Bits:")
	printIntAsBits(A)
	printIntAsBits(B)
	printIntAsBits(C)
	printIntAsBits(D)

	// Step 5 Output
	// return result
	fmt.Printf("%s %s %s %s\n", fmt.Sprintf("%x", A), fmt.Sprintf("%x", B), fmt.Sprintf("%x", C), fmt.Sprintf("%x", D))
	// It looks like the code you provided is a comment in the Go programming language. Comments in Go
	// start with "//" for single-line comments and "/*" and "*/" for multi-line comments. In this case,
	// the comment is simply "bs".

	mdBytes := make([]byte, 16)
	binary.LittleEndian.PutUint32(mdBytes[0:], A)
	binary.LittleEndian.PutUint32(mdBytes[4:], B)

	binary.LittleEndian.PutUint32(mdBytes[8:], C)
	binary.LittleEndian.PutUint32(mdBytes[12:], D)

	return hex.EncodeToString(mdBytes)
}

/* Let [abcd k s i] denote the operation
   a = b + ((a + F(b,c,d) + X[k] + T[i]) <<< s). */
/* Do the following 16 operations for each round. */
func roundOperation(a, b, c, d *uint32, f auxFunc, k, i, s uint32) uint32 {
	return rotateLeft((*a + f(*b, *c, *d) + k + i), s)
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

func printByteAsBits(b byte) {
	for i := 7; i >= 0; i-- {
		if b&(1<<i) != 0 {
			fmt.Print("1")
		} else {
			fmt.Print("0")
		}
	}
	fmt.Println()
}

func printIntAsBits(num uint32) {
	for i := 31; i >= 0; i-- {
		bit := (num >> i) & 1
		fmt.Print(bit)
		if i%8 == 0 {
			fmt.Print(" ") // Insert space for better readability every 8 bits
		}
	}
	fmt.Println()
}

func outputMd5(text string) string {
	// Create a new MD5 hash object
	hash := md5.New()

	// Write the data to the hash object
	hash.Write([]byte(text))

	// Compute the MD5 checksum
	hashInBytes := hash.Sum(nil)

	// Convert the bytes to a hexadecimal string
	hashInHex := hex.EncodeToString(hashInBytes)

	// Print the result
	fmt.Println("MD5 hash:", hashInHex)
	fmt.Println()
	return hashInHex
}

func main() {
	outputMd5("1")
	// fmt.Println()
	fmt.Println(md6("1"))

	// value := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	// fmt.Println(outputMd5(value) == md6(value))
	// fmt.Println(rotateLeft(1, 31))
	// var a uint32 = 1
	// var b uint32 = a

	// fmt.Println(a, b)
	// a = 2
	// fmt.Println(a, b)

	// var b uint32 = 0x89abcdef
	// var c uint32 = 0xfedcba98
	// var d uint32 = 0x76543210
	// fmt.Printf(fmt.Sprintf("%x\n", F(b, c, d)))

	// var one uint32 = 1
	// var zero uint32 = 1

	// // fmt.Println(^one)
	// printIntAsBits(zero ^ one)

	// var e uint32 = 0x2bd309f0
	// fmt.Printf(fmt.Sprintf("%x\n", rotateLeft(e, 7)))

	// fmt.Println(len("abcde"))
	// data := []byte("hello")
	// fmt.Printf("%x", md5.Sum(data))
}
