package basicblock

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"log"
	"testing"
)

func checkFatal(err error) {
	if err != nil {
		log.Fatalln("Fatal error.")
	}
}
func TestValidateCoinbaseTx(t *testing.T) {
	// privateKeyFrom, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	// checkFatal(err)
	// publicKeyFrom := privateKeyFrom.PublicKey

	privateKeyTo, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	checkFatal(err)
	publicKeyTo := privateKeyTo.PublicKey
	txOut := TxOut{publicKeyTo, CoinbaseAmount}
	tx := Transaction{
		txIns: []TxIn{
			TxIn{txOutIndex: int32(12)},
		},
		txOuts: []TxOut{txOut},
	}
	tx.id = tx.getID()
	if !validateCoinbaseTx(tx, 12) {
		fmt.Printf("validateCoinbaseTx failed\n")
		t.Fail()
	}
}
