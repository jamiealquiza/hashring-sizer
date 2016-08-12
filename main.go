package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"time"

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

	ring := ch.HashRing{Vnodes: 100000}

	for i := 0; i < 10; i++ {
		node := fmt.Sprintf("this-is-node-%d", i)
		ring.AddNode(node, node)
	}

	nodeCount := make(map[string]int)
	lat := make([]time.Duration, len(keys))

	for i, k := range keys {
		s := time.Now()
		n, _ := ring.GetNode(k)
		lat[i] = time.Since(s)
		nodeCount[n]++
	}

	// Print node counts.
	for k, v := range nodeCount {
		fmt.Printf("%s: %d\n", k, v)
	}

	// Get min/max latency.
	var highest time.Duration
	lowest := 10 * time.Second
	for _, v := range lat {
		if v > highest {
			highest = v
		}
		if v < lowest {
			lowest = v
		}
	}

	fmt.Printf("Max latency: %s\n", highest)
	fmt.Printf("Min latency: %s\n", lowest)

	// Get sorted latencies.
	latInt := make([]int, len(lat))
	for n := range lat {
		latInt[n] = int(lat[n])
	}

	sort.Ints(latInt)

	for n := range lat {
		lat[n] = time.Duration(latInt[n])
	}

	// Get other stuff.
	var latTotal time.Duration
	for _, v := range lat {
		latTotal += v
	}
	fmt.Printf("Avg. lantency: %s\n",
		time.Duration(int(latTotal)/len(lat)))

	if len(lat)%2 == 0 {
		med := (lat[len(lat)/2] + lat[(len(lat)/2)+1]) / 2
		fmt.Printf("Median latency: %s\n", med)
	} else {
		fmt.Printf("Median latency: %s\n",
			lat[int(math.Floor(float64(len(lat)/2)))])
	}
}
