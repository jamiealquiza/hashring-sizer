package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	ch "github.com/jamiealquiza/polymur/consistenthash"
	"github.com/jamiealquiza/tachymeter"
)

const nodes, vnodes = 10, 100

func main() {
	t := tachymeter.New(&tachymeter.Config{Size: 1000})

	f, err := os.Open("keys.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	var keys []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		keys = append(keys, scanner.Text())
	}

	ring := ch.HashRing{Vnodes: vnodes}

	for i := 0; i < nodes; i++ {
		node := fmt.Sprintf("this-is-node-%d", i)
		start := time.Now()
		ring.AddNode(node, node)
		t.AddTime(time.Since(start))
		t.AddCount(1)
	}

	fmt.Printf("\n> Inserted %d nodes with %d vnodes:\n", nodes, vnodes)
	t.Calc().Dump()
	t.Reset()
	fmt.Println()

	nodeCount := make(map[string]int)

	fmt.Printf("\n> Node balance:\n")
	for _, k := range keys {
		start := time.Now()
		n, _ := ring.GetNode(k)
		t.AddTime(time.Since(start))
		t.AddCount(1)
		nodeCount[n]++
	}

	var count int
	// Print node counts.
	for k, v := range nodeCount {
		fmt.Printf("%s: %d\n", k, v)
		count += v
	}

	fmt.Printf("\n> Performed %d lookups\n", count)
	t.Calc().Dump()
}
