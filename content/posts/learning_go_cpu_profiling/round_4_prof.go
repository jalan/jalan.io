package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"runtime/pprof"
	"strings"
)

var toNum = strings.NewReplacer(
	"A", string(0),
	"C", string(1),
	"G", string(2),
	"T", string(3),
)

const hashSize = 1 << 18

func hash(key uint64) int {
	return int(key) % hashSize
}

type entry struct {
	key   uint64
	value int
	next  *entry
}

type hashCounter struct {
	entries [hashSize]*entry
}

func (hc *hashCounter) get(key uint64) int {
	for e := hc.entries[hash(key)]; e != nil; e = e.next {
		if e.key == key {
			return e.value
		}
	}
	return 0
}

func (hc *hashCounter) inc(k uint64) {
	h := hash(k)
	p := &hc.entries[h]
	for e := *p; e != nil; e = e.next {
		if e.key == k {
			e.value++
			return
		}
	}
	e := &entry{k, 1, nil}
	if *p == nil {
		*p = e
	} else {
		e.next = *p
		*p = e
	}
}

func main() {
	f, _ := os.Create("round_4.prof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	sequence := "GGTATTTTAATTTATAGT"
	dna := read()
	counts := count(dna, len(sequence))
	seqCount := counts.get(encode([]byte(toNum.Replace(sequence))))
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

func count(dna []byte, length int) (counts hashCounter) {
	for i := 0; i < len(dna)-length+1; i++ {
		key := encode(dna[i : i+length])
		counts.inc(key)
	}
	return
}

func encode(sequence []byte) uint64 {
	var num uint64
	for _, char := range sequence {
		num = (num << 2) | uint64(char)
	}
	return num
}
