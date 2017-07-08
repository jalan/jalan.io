package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	sequence := "GGTATTTTAATTTATAGT"
	dna := read()
	counts := count(dna, len(sequence))
	seqCount := 0
	p, ok := counts[sequence]
	if ok {
		seqCount = *p
	}
	fmt.Printf("%v\t%v\n", seqCount, sequence)
}

func read() []byte {
	var buf bytes.Buffer
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		buf.WriteString(scanner.Text())
	}
	return buf.Bytes()
}

func count(dna []byte, length int) map[string]*int {
	counts := make(map[string]*int)
	for i := 0; i < len(dna)-length+1; i++ {
		key := string(dna[i : i+length])
		p, ok := counts[key]
		if !ok {
			p = new(int)
			counts[key] = p
		}
		*p++
	}
	return counts
}
