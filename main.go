package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	ch "github.com/jamiealquiza/polymur/consistenthash"
	"github.com/jamiealquiza/tachymeter"
)

func main() {
	var nodes, vnodes int

	flag.IntVar(&nodes, "nodes", 10, "Number of nodes")
	flag.IntVar(&vnodes, "vnodes", 10, "Number of vnodes per node")

	flag.Parse()

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

	t := tachymeter.New(&tachymeter.Config{Size: len(keys)})

	ring := ch.HashRing{Vnodes: vnodes}

	for i := 0; i < nodes; i++ {
		node := fmt.Sprintf("this-is-node-%d", i)
		start := time.Now()
		ring.AddNode(node, node)
		t.AddTime(time.Since(start))
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
