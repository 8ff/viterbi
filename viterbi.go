package viterbi

import (
	"fmt"
	"os"
	"strconv"
)

type Input struct {
	Constraint         int
	Polynomials        []int
	ReversePolynomials bool
}

var FLAGS_reverse_polynomials bool
var FLAGS_encode bool

type Trellis [][]int

const MaxInt = 1<<(32-1) - 1

type ViterbiCodec struct {
	constraint_  int
	polynomials_ []int
	outputs_     []string
}

func (v *ViterbiCodec) num_parity_bits() int {
	return len(v.polynomials_)
}

func (v *ViterbiCodec) NextState(current_state int, input int) int {
	return (current_state >> 1) | (input << (v.constraint_ - 2))
}

func (v *ViterbiCodec) Output(current_state int, input int) string {
	ind := int(current_state | input<<(v.constraint_-1))
	return v.outputs_[ind]
}

func (v *ViterbiCodec) InitializeOutputs() {
	new_len := (1 << v.constraint_)
	for i := 0; i < new_len; i++ {
		v.outputs_ = append(v.outputs_, "")

		for j := 0; j < v.num_parity_bits(); j++ {
			// Reverse polynomial bits to make the convolution code simpler.
			polynomial := ReverseBits(v.constraint_, v.polynomials_[j])
			input := i
			output := 0
			for k := 0; k < v.constraint_; k++ {
				output ^= (input & 1) & (polynomial & 1)
				polynomial >>= 1
				input >>= 1
			}
			if output == 1 {
				v.outputs_[i] += "1"
			} else {
				v.outputs_[i] += "0"
			}

		}
	}

}

// Encode the given message bits.
func (v *ViterbiCodec) Encode(bits string) string {
	state := 0
	var encoded string

	// Encode the message bits.
	for i := 0; i < len(bits); i++ {
		c := bits[i]
		input := int(c) - int('0') //c - '0'

		encoded += v.Output(state, input)

		state = v.NextState(state, input)
	}

	for i := 0; i < v.constraint_-1; i++ {
		encoded += v.Output(state, 0)
		state = v.NextState(state, 0)
	}

	return encoded
}

// Return the branch metric for the given source and target states.
func (v *ViterbiCodec) BranchMetric(bits string, source_state int, target_state int) int {
	var output string = v.Output(source_state, target_state>>(v.constraint_-2))

	return HammingDistance(bits, output)
}

// Return the path metric and the source state for the given target state.
func (v *ViterbiCodec) PathMetric(bits string, prev_path_metrics []int, state int) [2]int {
	ret := [2]int{0, 0}
	s := (state & ((1 << (v.constraint_ - 2)) - 1)) << 1
	source_state1 := s
	source_state2 := s | 1

	pm1 := prev_path_metrics[source_state1]
	if pm1 < MaxInt {
		pm1 += v.BranchMetric(bits, source_state1, state)
	}

	pm2 := prev_path_metrics[source_state2]
	if pm2 < MaxInt {
		pm2 += v.BranchMetric(bits, source_state2, state)
	}

	if pm1 <= pm2 {
		ret[0] = pm1
		ret[1] = source_state1
	} else {
		ret[0] = pm2
		ret[1] = source_state2
	}
	return (ret)
}

// Update the path metrics and trellis for the next bit.
func (v *ViterbiCodec) UpdatePathMetrics(bits string, path_metrics []int, trellis Trellis) ([]int, []int) {
	var new_path_metrics []int
	var new_trellis_column []int

	for w := 0; w < len(path_metrics); w++ {
		p := v.PathMetric(bits, path_metrics, w)
		new_path_metrics = append(new_path_metrics, p[0])
		new_trellis_column = append(new_trellis_column, p[1])
	}

	return new_path_metrics, new_trellis_column
}

// Decode a string of bits using the Viterbi algorithm.
func (v *ViterbiCodec) Decode(bits string) string {
	var trellis Trellis
	var auxTrellis []int
	var path_metrics []int
	var current_bits string

	for i := 0; i < (1 << (v.constraint_ - 1)); i++ {
		path_metrics = append(path_metrics, MaxInt)
	}

	path_metrics[0] = 0
	for i := 0; i < len(bits); i += v.num_parity_bits() {

		if i+v.num_parity_bits() >= len(bits) {
			current_bits = bits[i:]
		} else {
			current_bits = bits[i : i+v.num_parity_bits()]
		}

		// If some bits are missing, fill with trailing zeros.
		// This is not ideal but it is the best we can do.
		for len(current_bits) < v.num_parity_bits() {
			current_bits += "0"
		}

		path_metrics, auxTrellis = v.UpdatePathMetrics(current_bits, path_metrics, trellis)
		trellis = append(trellis, auxTrellis)
	}

	// Traceback.
	var decoded string
	var state int = findMin(path_metrics)

	for i := len(trellis) - 1; i >= 0; i-- {

		if state>>(v.constraint_-2) == 1 {
			decoded += "1"
		} else {
			decoded += "0"
		}
		state = trellis[i][state]

	}

	reverse := reverseString(decoded)
	return reverse[0 : len(reverse)-v.constraint_+1]

}

func reverseString(str string) string {
	byte_str := []rune(str)
	for i, j := 0, len(byte_str)-1; i < j; i, j = i+1, j-1 {
		byte_str[i], byte_str[j] = byte_str[j], byte_str[i]
	}
	return string(byte_str)
}

func findMin(v []int) int {
	min := v[0]
	ind := 0

	for i, value := range v {
		if value < min {
			min = value
			ind = i
		}

	}
	return ind
}

func HammingDistance(x string, y string) int {
	distance := 0
	for i := 0; i < len(x); i++ {
		if x[i] != y[i] {
			distance += 1
		}
	}
	return distance
}

func ReverseBits(num_bits int, input int) int {
	var output int
	output = 0
	for i := num_bits - 1; i >= 0; i-- {
		output = (output << 1) + (input & 1)
		input >>= 1
	}
	return output
}

func ParseInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("Expected a number, found %s\n", s)
		os.Exit(0)
	}
	return val
}

// Function that takes in []byte and converts it to a string of bits.
func BytesToBits(bytes []byte) string {
	var bits string
	for _, b := range bytes {
		bits += fmt.Sprintf("%08b", b)
	}
	return bits
}

// Function that takes in a string of bits and converts it to a []byte.
func BitsToBytes(bits string) []byte {
	// Check if the number of bits is a multiple of 8.
	if len(bits)%8 != 0 {
		return nil
	}

	var bytes []byte
	for i := 0; i < len(bits); i += 8 {
		b, _ := strconv.ParseUint(bits[i:i+8], 2, 8)
		bytes = append(bytes, byte(b))
	}
	return bytes
}

func Init(input Input) (*ViterbiCodec, error) {
	// Go over polynomials and check if they are valid.
	for i := 0; i < len(input.Polynomials); i++ {
		if input.Polynomials[i] <= 0 {
			return nil, fmt.Errorf("polynomial should be greater than 0, found %d", input.Polynomials[i])
		}

		if input.Polynomials[i] >= (1 << input.Constraint) {
			return nil, fmt.Errorf("polynomial should be less than %d found %d", (1 << input.Constraint), input.Polynomials[i])
		}
	}

	// Go over polynomials and reverse them if needed.
	if input.ReversePolynomials {
		for i := 0; i < len(input.Polynomials); i++ {
			input.Polynomials[i] = ReverseBits(input.Constraint, input.Polynomials[i])
		}
	}

	var codec = ViterbiCodec{constraint_: input.Constraint, polynomials_: input.Polynomials}
	codec.InitializeOutputs()
	return &codec, nil
}

// Encode bytes using BytesToBits and run codec.Encode.
func (v *ViterbiCodec) EncodeBytes(bytes []byte) string {
	return v.Encode(BytesToBits(bytes))
}

// Decode bytes using BitsToBytes and run codec.Decode.
func (v *ViterbiCodec) DecodeBytes(bits string) []byte {
	return BitsToBytes(v.Decode(bits))
}

// Encode bytes using BitsToBytes and run codec.Encode then return bytes using BytesToBits or error
func (v *ViterbiCodec) EncodeBytesToBytes(bytes []byte) ([]byte, error) {
	bits := v.Encode(BytesToBits(bytes))
	bytesToReturn := BitsToBytes(bits)
	if bytesToReturn == nil {
		return nil, fmt.Errorf("could not convert bits to bytes")
	}
	return bytesToReturn, nil
}

// Decode bytes using BytesToBits and run codec.Decode then return bytes using BitsToBytes or error
func (v *ViterbiCodec) DecodeBytesToBytes(bytes []byte) ([]byte, error) {
	bits := BytesToBits(bytes)
	bytesToReturn := BitsToBytes(v.Decode(bits))
	if bytesToReturn == nil {
		return nil, fmt.Errorf("could not convert bits to bytes")
	}
	return bytesToReturn, nil
}
