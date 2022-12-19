# 📡 Convolutional Encoder and Viterbi Decoder

This package implements a convolutional encoder and a Viterbi decoder.
It can be used as a library or as a command line tool.

## Using as a library
```go
package main

import (
	"fmt"
	"github.com/8ff/viterbi"
)

func main() {
	// Initialize a codec.
	codec, err := viterbi.Init(viterbi.Input{Constraint: 7, Polynomials: []int{91, 109, 121}, ReversePolynomials: false})
	if err != nil {
		panic(err)
	}

	/**** Encode Bytes ****/
	input := []byte("Hello, world!")

	// Encode a message.
	encoded := codec.EncodeBytes(input)

	// Decode the message.
	decoded := codec.DecodeBytes(encoded)

	// Print input and decoded message.
	fmt.Printf("InputBytes:   %s\n", input)
	fmt.Printf("DecodedBytes: %s\n", decoded)

	/**** Encode string of bits ****/
	inputBits := "101010"

	// Encode a message.
	encodedBits := codec.Encode(inputBits)

	// Decode the message.
	decodedBits := codec.Decode(encodedBits)

	// Print input and decoded message.
	fmt.Printf("InputBits:   %s\n", inputBits)
	fmt.Printf("DecodedBits: %s\n", decodedBits)
}

```

## Using as command line tool
```bash
cd cmd/viterbiCli

# Encode
echo 101010 | go run viterbiCli.go encode 3 5 7

# Decode
echo 10000010 | go run viterbiCli.go decode 3 5 7

# Encode/Decode
echo 101010 | go run viterbiCli.go encode 3 5 7 | go run viterbiCli.go decode 3 5 7
```

## Error Handling
Inputs are validated, and proper error messages will be displayed.

* The constraint should be greater than 0.
* A generator polynomial should be greater than 0 and less than 1 << constraint.
* The received bit sequence should be of length N * <num-of-polynomials> where N is an integer. Otherwise some bits must be missing during transmission. We will fill in appropriate number of trailing zeros.

## Dependencies
This code contains no external dependencies.