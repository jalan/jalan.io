+++
date = "2016-09-16"
draft = true
title = "Learning Go: CPU Profiling"
+++

This is the first post in what I hope becomes a series about learning the Go
programming language (golang) and tools. This time, let's create a simple
program, profile it, and try to make it go fast.

## How about a game?

[The Computer Language Benchmarks Game] is a fun source of simple program
requirements. Their [k-nucleotide benchmark] appears to be one on which Go
[does not perform very well], losing to 15 other languages at the time I begin
this post. To make this exercise fun, we'll use an adaptation of that
benchmark--not exactly that benchmark, but a fair representation of it. Here
are our requirements:

[The Computer Language Benchmarks Game]: http://benchmarksgame.alioth.debian.org/
[k-nucleotide benchmark]: http://benchmarksgame.alioth.debian.org/u64q/knucleotide-description.html#knucleotide
[does not perform very well]: http://benchmarksgame.alioth.debian.org/u64q/performance.php?test=knucleotide

1. Read a DNA sequence from standard input. It will look something like this:

        CCCATAACTACAATAGTCGGCAATCTTTTATTACCCAGAACTAACGTTTTTATTTCCCGG
        TACGTATCACATTAATCTTAATTTAATGCGTGAGAGTAACGATGAACGAAAGTTATTTAT
        GTTTAAGCCGCTTCTTGAGAATACAGATTACTGTTAGAATGAAGGCATCATAACTAGAAC
        ACCAACGCGCACCTCGCACATTACTCTAATAGTAGCTTTATTCAGTTTAATATAGACAGT
        .
        .
        .

2. Count up all the subsequences of a given length.

3. Print the count for a specific subsequence, just to make sure everything
   works.

The full input will be about 120 MB. We'll count up all subsequences of length
18 and then print the count for `GGTATTTTAATTTATAGT`.

## Round 1

Here is a straightforward implementation, with one function to read the input,
another function to count all subsequences of a given length, and a main
function.

```go
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
```

It "takes" about 24 seconds to run on my machine and finds 893 instances of the
specified subsequence. In order to see where the program spends these 24
seconds, we can add a bit of code to enable the CPU profiler. With
`runtime/pprof` added to the imports and the following added to the beginning
of `main`, building and running the program again saves a profile to
`round_1.prof`.

```go
	f, _ := os.Create("round_1.prof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
```

The [pprof] tool can be used to analyze the resulting profile.

[pprof]: https://github.com/google/pprof



```
$ go tool pprof round_1 round_1.prof
Entering interactive mode (type "help" for commands)
(pprof) top 8
20950ms of 26430ms total (79.27%)
Dropped 67 nodes (cum <= 132.15ms)
Showing top 8 nodes out of 40 (cum >= 4710ms)
      flat  flat%   sum%        cum   cum%
    9130ms 34.54% 34.54%    11960ms 45.25%  runtime.mapaccess1_faststr
    2420ms  9.16% 43.70%     3810ms 14.42%  runtime.mallocgc
    2220ms  8.40% 52.10%     5080ms 19.22%  runtime.mapassign1
    2110ms  7.98% 60.08%     2110ms  7.98%  runtime.memeqbody
    1460ms  5.52% 65.61%    24500ms 92.70%  main.count
    1410ms  5.33% 70.94%     1410ms  5.33%  runtime.aeshashbody
    1290ms  4.88% 75.82%     1290ms  4.88%  runtime.memmove
     910ms  3.44% 79.27%     4710ms 17.82%  runtime.rawstring
```

## Round 2: Pointers

Let's first focus on the third most expensive function in the pprof output:
`runtime.mapassign1`. In Go, map elements cannot be modified, so every

```go
		counts[string(dna[i:i+length])]++
```

is actually a map access, then an increment, and then a map assign. On each
iteration through the count loop, we pay for two map operations.


We can avoid almost all the map assigns by storing the counts themselves
outside the map and storing pointers to the counts in the map.




This can be
verified with another `pprof` command, `disasm`, which shows the assembly for a
given function.

```
(pprof) disasm count
.
.
.
      40ms        12s     40163f: CALL runtime.mapaccess1_faststr(SB)
         .          .     401644: MOVQ 0x20(SP), BX
      80ms       80ms     401649: MOVQ 0(BX), BX
     780ms      780ms     40164c: INCQ BX
      20ms       20ms     40164f: MOVQ BX, 0x38(SP)
      10ms       10ms     401654: LEAQ 0xc42a5(IP), BX
         .          .     40165b: MOVQ BX, 0(SP)
      10ms       10ms     40165f: MOVQ 0x40(SP), BX
         .          .     401664: MOVQ BX, 0x8(SP)
      30ms       30ms     401669: LEAQ 0x48(SP), BX
         .          .     40166e: MOVQ BX, 0x10(SP)
      30ms       30ms     401673: LEAQ 0x38(SP), BX
         .          .     401678: MOVQ BX, 0x18(SP)
      20ms      5.10s     40167d: CALL runtime.mapassign1(SB)
.
.
.
```






Instead of storing integers in a map, we can store pointers


[go/issues/3117]: https://github.com/golang/go/issues/3117

## Round 3: Key Encoding




## Round 4: Faster Maps



## Results



## Caveat: Don't Do Any of This

The usual advice about performance optimization applies

Professional driver on a closed course.

In this demo I created some pretty nasty code that I would not want to put in production.

- "Should I use pointers for counting something with a map?"
- "Do I need to reencode my data so each key fits in 64 bits?"
- "Are Go maps so slow that I need to implement my own?"

Probably not!
