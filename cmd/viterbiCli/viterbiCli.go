package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/8ff/viterbi"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: viterbiCli encode/decode <constraint> <polynomial> <polynomial> ... <message>")
		os.Exit(1)
	}

	constraint, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid constraint: %s\n", args[2])
		os.Exit(2)
	}

	// Read polynomials as remaining arguments except the last one
	poly := make([]int, len(args)-4)
	for i, p := range args[3 : len(args)-1] {
		poly[i], err = strconv.Atoi(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid polynomial: %s\n", p)
			os.Exit(2)
		}
	}

	// Read message from stdin
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	message := scanner.Text()

	// fmt.Fprintf(os.Stderr, "constraint: %d, poly: %v, message: %s\n", constraint, poly, message)

	// Init viterbi codec
	codec, err := viterbi.Init(viterbi.Input{Constraint: constraint, Polynomials: poly, ReversePolynomials: false})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}

	// Switch for encode/decode
	switch args[1] {
	case "encode":
		// Encode message
		encoded := codec.Encode(message)

		// Print encoded message
		fmt.Println(encoded)
	case "decode":
		// Decode message
		decoded := codec.Decode(message)

		// Print decoded message
		fmt.Println(decoded)
	default:
		fmt.Fprintf(os.Stderr, "Invalid command: %s\n", args[1])
		os.Exit(1)
	}
}
