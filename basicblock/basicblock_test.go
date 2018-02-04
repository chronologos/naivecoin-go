package basicblock

import (
	"testing"
	"time"
)

var TestBlock1 = GenesisBlock.FindBlock([]byte{})
var TestBlock2 = TestBlock1.FindBlock([]byte{})

func TestGetConseqZeroes(t *testing.T) {
	if getConseqZeroes(byte(0)) != 8 {
		t.Fail()
	}
	if getConseqZeroes(byte(240)) != 0 {
		t.Fail()
	}
	if getConseqZeroes(byte(15)) != 4 {
		t.Fail()
	}
	if getConseqZeroes(byte(1)) != 7 {
		t.Fail()
	}
}

func TestHashesMatchDifficulties(t *testing.T) {
	hash1 := []byte{79, 0}
	hash2 := []byte{0, 1}
	if hashMatchesDifficulty(4, hash1) {
		t.Fail()
	}
	if !hashMatchesDifficulty(0, hash1) {
		t.Fail()
	}
	if hashMatchesDifficulty(16, hash2) {
		t.Fail()
	}
	if !hashMatchesDifficulty(15, hash2) {
		t.Fail()
	}
}

func TestEmptyBlockchain(t *testing.T) {
	blockChain := []BasicBlock{}
	if IsValidBasicBlockchain(blockChain) {
		t.Fail()
	}
}

func TestInvalidExtraBlock(t *testing.T) {
	blockChain := []BasicBlock{GenesisBlock}
	for i := 0; i < 5; i++ {
		blockChain = append(blockChain, blockChain[len(blockChain)-1].FindBlock([]byte{}))
	}
	blockChain = append(blockChain, BasicBlock{})
	if IsValidBasicBlockchain(blockChain) {
		t.Fail()
	}
}

func TestInvalidGenesisBlock(t *testing.T) {
	blockChain := []BasicBlock{GenesisBlock}
	for i := 0; i < 5; i++ {
		blockChain = append(blockChain, blockChain[len(blockChain)-1].FindBlock([]byte{}))
	}
	blockChain[0].Data = []byte("DEADBEEF")
	if IsValidBasicBlockchain(blockChain) {
		t.Fail()
	}
}

func TestBlockValidation(t *testing.T) {
	TestBlock2HashWrong := TestBlock2
	TestBlock2HashWrong.Hash = [32]byte{}
	TestBlock1HashWrong := TestBlock1
	TestBlock1HashWrong.Hash = [32]byte{}
	TestBlock2MutatedData := TestBlock2
	TestBlock2MutatedData.Data = []byte("DEADBEEF")
	TestBlock2TimestampTooEarly := TestBlock2
	TestBlock2TimestampTooEarly.Timestamp = TestBlock2TimestampTooEarly.Timestamp.Add(-62 * time.Second)
	TestBlock2TimestampOk := TestBlock2
	TestBlock2TimestampOk.Timestamp = TestBlock2TimestampOk.Timestamp.Add(-5 * time.Second)
	// log.Printf("invalid timestamps: %s and %s \n", TestBlock2TimestampOk.Timestamp.String(), TestBlock1.Timestamp.String())

	if !TestBlock2.IsValidBasicBlock(&TestBlock1) {
		t.Fail()
	}
	if TestBlock2HashWrong.IsValidBasicBlock(&TestBlock1) {
		t.Fail()
	}
	if TestBlock2.IsValidBasicBlock(&TestBlock1HashWrong) {
		t.Fail()
	}
	if TestBlock2MutatedData.IsValidBasicBlock(&TestBlock1) {
		t.Fail()
	}
	if TestBlock2TimestampTooEarly.isValidTimestamp(&TestBlock1) {
		t.Fail()
	}
	if !TestBlock2TimestampOk.isValidTimestamp(&TestBlock1) {
		t.Fail()
	}
}

func TestTimestampAttack(t *testing.T) {
	blockChain := []BasicBlock{GenesisBlock}
	for i := 0; i < 5; i++ {
		blockChain = append(blockChain, blockChain[len(blockChain)-1].FindBlock([]byte{}))
	}
	blockChain[len(blockChain)-1].Timestamp = blockChain[len(blockChain)-1].Timestamp.Add(-61 * time.Second)
	if IsValidBasicBlockchain(blockChain) {
		t.Fail()
	}
}

func TestValidBlockchain(t *testing.T) {
	blockChain := []BasicBlock{GenesisBlock}
	for i := 0; i < 5; i++ {
		blockChain = append(blockChain, blockChain[len(blockChain)-1].FindBlock([]byte{}))
	}
	if !IsValidBasicBlockchain(blockChain) {
		t.Fail()
	}
}

// func TestBlockchainReplaced(t *testing.T) {
// 	blockChainShort := []BasicBlock{GenesisBlock}
// 	blockChainLong := []BasicBlock{GenesisBlock}
//
// 	for i := 0; i < 1; i++ {
// 		blockChainLong = append(blockChainLong, blockChainLong[len(blockChainLong)-1].GenerateNextBasicBlock([]byte{}))
// 	}
//
// 	res := PossiblyReplace(blockChainShort, blockChainLong)
// 	if !deepEqual(res, blockChainShort) {
// 		t.Fail()
// 	}
// }
