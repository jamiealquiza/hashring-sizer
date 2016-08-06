package main

import (
	"bufio"
	"fmt"
	"os"

	ch "github.com/jamiealquiza/polymur/consistenthash"
)


func main() {
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

	ring := ch.HashRing{Vnodes: 100}

	for i := 0; i < 10; i++ {
		node := fmt.Sprintf("this-is-node-%d", i)
		ring.AddNode(node, node)
	}

	nodeCount := make(map[string]int)

	for _, k := range keys {
		n, _ := ring.GetNode(k)
		nodeCount[n]++	
	}

	for k, v := range nodeCount {
		fmt.Printf("%s: %d\n", k, v)
	}
}
