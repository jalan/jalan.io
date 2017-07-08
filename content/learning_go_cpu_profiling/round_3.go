package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

var toNum = strings.NewReplacer(
	"A", string(0),
	"C", string(1),
	"G", string(2),
	"T", string(3),
)

func main() {
	sequence := "GGTATTTTAATTTATAGT"
	dna := read()
	counts := count(dna, len(sequence))
	seqCount := 0
	p, ok := counts[encode([]byte(toNum.Replace(sequence)))]
	if ok {
		seqCount = *p
	}
	fmt.Printf("%v\t%v\n", seqCount, sequence)
}

func read() []byte {
	var buf bytes.Buffer
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		buf.WriteString(toNum.Replace(scanner.Text()))
	}
	return buf.Bytes()
}

func count(dna []byte, length int) map[uint64]*int {
	counts := make(map[uint64]*int)
	for i := 0; i < len(dna)-length+1; i++ {
		key := encode(dna[i : i+length])
		p, ok := counts[key]
		if !ok {
			p = new(int)
			counts[key] = p
		}
		*p++
	}
	return counts
}

func encode(sequence []byte) uint64 {
	var num uint64
	for _, char := range sequence {
		num = (num << 2) | uint64(char)
	}
	return num
}
