package naivecoin

import (
	"fmt"

	nc "github.com/chronologos/naivecoin"
)

// BlockChain!

func main() {

	blockChain := []nc.BasicBlock{nc.GenesisBlock}
	for i := 0; i < 5; i++ {
		blockChain = append(blockChain, blockChain[len(blockChain)-1].GenerateNextBasicBlock([]byte{}))
	}

	fmt.Print(blockChain)
	fmt.Print(nc.ValidBasicBlockchain(blockChain))
}
