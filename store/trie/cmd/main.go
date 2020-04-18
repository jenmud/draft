package main

import (
	"log"

	"bitbucket.dataprocessors.com.au/st5g/memgraph/index/trie"
)

func main() {
	t := trie.New()

	if err := t.Add("bat", "bat", false); err != nil {
		log.Printf("%s", err)
	}

	if err := t.Add("cat", "cat", false); err != nil {
		log.Printf("%s", err)
	}

	if err := t.Add("co", "co", false); err != nil {
		log.Printf("%s", err)
	}

	if err := t.Add("cot", "cot", false); err != nil {
		log.Printf("%s", err)
	}

	if err := t.Add("cought", "cought", false); err != nil {
		log.Printf("%s", err)
	}

	if err := t.Add("count", "count", true); err != nil {
		log.Printf("%s", err)
	}

	for _, k := range t.Keys() {
		log.Printf("%s", k)
	}
}
