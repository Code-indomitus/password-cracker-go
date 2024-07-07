package utils

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"math/bits"
)

func Md5(plainText string) string {
	// Prepare to append to input
	plainTextByteSize := len(plainText)
	remainingBits := (plainTextByteSize * 8) % 512
	var bytesToAppend int

	// Find out how many bytes to append
	if remainingBits < 448 {
		bytesToAppend = (448 - remainingBits) / 8
	} else if remainingBits == 448 {
		bytesToAppend = 448 / 8
	} else {
		bytesToAppend = (512 - remainingBits + 448) / 8
	}

	// Create a list of bytes out of input string
	plainTextBytes := []byte(plainText)
	plainTextBytes = append(plainTextBytes, 0x80) // append "1" bit to the end of the message
	bytesToAppend--                               // decrement the number of bytes to append by 1

	// Append remaing required "0" bits till congruent to 448 mod 512
	for i := 0; i < bytesToAppend; i++ {
		plainTextBytes = append(plainTextBytes, 0x00)
	}

	// append length of message in bits as 64-bit little-endian integer to message buffer
	var plainTextBitLength uint64 = uint64(plainTextByteSize) * 8 // the length is the number of bits in the original message
	lengthBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(lengthBytes, plainTextBitLength)
	paddedPlainTextBytes := append(plainTextBytes, lengthBytes...)

	// Initialize MD buffer values in little endian form
	var A uint32 = 0x67452301
	var B uint32 = 0xEFCDAB89
	var C uint32 = 0x98BADCFE
	var D uint32 = 0x10325476

	T := make([]uint32, 64)
	for i := 0; i < 64; i++ {
		T[i] = uint32(math.Pow(2, 32) * math.Abs(math.Sin(float64(i+1))))
	}

	// Iterate over each 512 bit block
	for i := 0; i <= len(paddedPlainTextBytes)-64; i += 64 {
		// break the block into sixteen 32-bit words
		// Initialize hash value for this chunk, in map M
		X := make([]uint32, 16)
		for j := 0; j < 16; j++ {
			start := i + (j * 4)
			X[j] = binary.LittleEndian.Uint32(paddedPlainTextBytes[start : start+4])
		}

		// Save state of A, B, C, D
		var AA uint32 = A
		var BB uint32 = B
		var CC uint32 = C
		var DD uint32 = D

		// Round 1
		A = B + bits.RotateLeft32(F(B, C, D)+A+X[0]+T[0], 7) // [ABCD  0  7  1]
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
	}

	// Concatenate little endian of A, B, C, D to get final 128 bit hash value
	mdBytes := make([]byte, 16)
	binary.LittleEndian.PutUint32(mdBytes[0:], A)
	binary.LittleEndian.PutUint32(mdBytes[4:], B)
	binary.LittleEndian.PutUint32(mdBytes[8:], C)
	binary.LittleEndian.PutUint32(mdBytes[12:], D)

	return hex.EncodeToString(mdBytes)
}

/*
 Auxiliar functions
*/

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

func BytesToInt(binary []byte) uint32 {
	return uint32(binary[0])<<24 | uint32(binary[1])<<16 | uint32(binary[2])<<8 | uint32(binary[3])
}

func PrintByteAsBits(b byte) {
	for i := 7; i >= 0; i-- {
		if b&(1<<i) != 0 {
			fmt.Print("1")
		} else {
			fmt.Print("0")
		}
	}
	fmt.Println()
}

func PrintIntAsBits(num uint32) {
	for i := 31; i >= 0; i-- {
		bit := (num >> i) & 1
		fmt.Print(bit)
		if i%8 == 0 {
			fmt.Print(" ") // Insert space for better readability every 8 bits
		}
	}
	fmt.Println()
}
