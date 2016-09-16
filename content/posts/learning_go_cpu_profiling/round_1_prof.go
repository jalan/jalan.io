package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"runtime/pprof"
)

func main() {
	f, _ := os.Create("round_1.prof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	sequence := "GGTATTTTAATTTATAGT"
	dna := read()
	counts := count(dna, len(sequence))
	fmt.Printf("%v\t%v\n", counts[sequence], sequence)
}

func read() []byte {
	var buf bytes.Buffer
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		buf.WriteString(scanner.Text())
	}
	return buf.Bytes()
}

func count(dna []byte, length int) map[string]int {
	counts := make(map[string]int)
	for i := 0; i < len(dna)-length+1; i++ {
		counts[string(dna[i:i+length])]++
	}
	return counts
}
