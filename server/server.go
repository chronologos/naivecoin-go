package main

import (
	"fmt"
	"log"
	"net/http"

	bb "github.com/chronologos/naivecoin/basicblock"
)

var blockChain []bb.BasicBlock

func main() {
	blockChain = []bb.BasicBlock{bb.GenesisBlock}
	for i := 0; i < 5; i++ {
		blockChain = append(blockChain, blockChain[len(blockChain)-1].GenerateNextBasicBlock([]byte{}))
	}
	http.HandleFunc("/", displayBlockchain)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func displayBlockchain(w http.ResponseWriter, r *http.Request) {
	for _, blk := range blockChain {
		fmt.Fprint(w, blk.String()+"\n")
	}
}
