package viterbi

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
)

func init() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
}

// Function that takes in a string of 1s and 0s and randomly flips n unique bits.
func corruptBits(bits string, n int) string {
	corrupted := []rune(bits)
	indices := make(map[int]bool)
	for i := 0; i < n; i++ {
		// Pick a random index
		index := rand.Intn(len(bits))
		// If we've already flipped this bit, pick another index
		for indices[index] {
			index = rand.Intn(len(bits))
		}
		// Flip the bit
		if corrupted[index] == '0' {
			corrupted[index] = '1'
		} else {
			corrupted[index] = '0'
		}
		// Mark this index as flipped
		indices[index] = true
	}
	return string(corrupted)
}

func TestEncodeDecode(t *testing.T) {
	// Vars
	constraint := 7
	polynomials := []int{91, 109, 121}
	inputLen := 20
	numOfCorruptBits := 1
	reversePolynomials := false

	// Generate random input bytes
	inputBytes := make([]byte, inputLen)
	for i := 0; i < 6; i++ {
		inputBytes[i] = byte(rand.Intn(256))
	}

	codec, err := Init(Input{Constraint: constraint, Polynomials: polynomials, ReversePolynomials: reversePolynomials})
	if err != nil {
		t.Fatal(err)
	}

	// Encode
	encoded := codec.Encode(BytesToBits(inputBytes))

	// Corrupt bits
	encoded = corruptBits(encoded, numOfCorruptBits)

	// Decode
	decoded := codec.Decode(encoded)

	// Check if decoded is equal to input
	if !bytes.Equal(BitsToBytes(decoded), inputBytes) {
		t.Errorf("Expected %s, got %s", BytesToBits(inputBytes), decoded)
	}
}

// Test EncodeDecode 100 times
func TestEncodeDecode100(t *testing.T) {
	for i := 0; i < 1000; i++ {
		TestEncodeDecode(t)
	}
}

// Test encode/decode using EncodeBytes and DecodeBytes
func TestEncodeDecodeBytes(t *testing.T) {
	// Vars
	constraint := 7
	polynomials := []int{91, 109, 121}
	inputLen := 20
	numOfCorruptBits := 1
	reversePolynomials := false

	// Generate random input bytes
	inputBytes := make([]byte, inputLen)
	for i := 0; i < 6; i++ {
		inputBytes[i] = byte(rand.Intn(256))
	}

	codec, err := Init(Input{Constraint: constraint, Polynomials: polynomials, ReversePolynomials: reversePolynomials})
	if err != nil {
		t.Fatal(err)
	}

	// Encode
	encoded := codec.EncodeBytes(inputBytes)

	// Corrupt bits
	encoded = corruptBits(encoded, numOfCorruptBits)

	// Decode
	decoded := codec.DecodeBytes(encoded)

	// Check if decoded is equal to input
	if !bytes.Equal(decoded, inputBytes) {
		t.Errorf("Expected %s, got %s", BytesToBits(inputBytes), decoded)
	}
}

func TestSeparateEncodeDecoders(t *testing.T) {
	// Vars
	constraint := 7
	polynomials := []int{91, 109, 121}
	inputLen := 20
	numOfCorruptBits := 1
	reversePolynomials := false // Somehow when this is true and using different encode/decode codecs, the test fails

	// Generate random input bytes
	inputBytes := make([]byte, inputLen)
	for i := 0; i < 6; i++ {
		inputBytes[i] = byte(rand.Intn(256))
	}

	codec, err := Init(Input{Constraint: constraint, Polynomials: polynomials, ReversePolynomials: reversePolynomials})
	if err != nil {
		t.Fatal(err)
	}

	// Encode
	encoded := codec.Encode(BytesToBits(inputBytes))

	// Corrupt bits
	encoded = corruptBits(encoded, numOfCorruptBits)

	decodeCodec, err := Init(Input{Constraint: constraint, Polynomials: polynomials, ReversePolynomials: reversePolynomials})
	if err != nil {
		t.Fatal(err)
	}

	// Decode
	decoded := decodeCodec.Decode(encoded)

	// Check if decoded is equal to input
	if !bytes.Equal(BitsToBytes(decoded), inputBytes) {
		t.Errorf("Expected %s, got %s", BytesToBits(inputBytes), decoded)
	}
}
