package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	bb "github.com/chronologos/naivecoin/basicblock"
	"github.com/gorilla/websocket"
)

var ip = flag.String("ip", "80", "ip address for this server")
var mines = flag.Bool("mines", false, "True if this servdr actually mines blocks.")
var wsconns []*websocket.Conn
var blockChain bb.BlockChain
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var ticker *time.Ticker
var inCh chan bb.BlockChain
var outCh chan bb.BlockChain
var dialer = &websocket.Dialer{
	Proxy: http.ProxyFromEnvironment,
}

func main() {
	flag.Parse()
	blockChain = []bb.BasicBlock{bb.GenesisBlock}

	outCh = make(chan bb.BlockChain)
	inCh = make(chan bb.BlockChain)
	go wsWriter(outCh)

	http.HandleFunc("/", displayIndex)
	http.HandleFunc("/blocks", displayBlockchain)
	http.HandleFunc("/p", parsePost)
	http.HandleFunc("/ws", websocketHandler)

	ticker = time.NewTicker(5 * time.Second) // TODO(chronologos) remove eventually, when we have real mining.
	defer ticker.Stop()

	var s string
	if *mines {
		s = "mining node"
		go mine(outCh)

	} else {
		s = "non-mining node"
	}

	go updateBlockchain(inCh)

	fmt.Printf("🖥 Server initialized, listening on port %s, %s. \n", *ip, s)
	log.Fatal(http.ListenAndServe("localhost:"+*ip, nil))
}

func adjustDifficulty(bc bb.BlockChain) {
	var err error
	bb.Difficulty, err = bb.GetDifficulty(bc)
	if err != nil {
		log.Fatalln("Blockchain length is 0")
	}
}

func mine(ch chan<- bb.BlockChain) {
	for {
		<-ticker.C
		newBlock := blockChain[len(blockChain)-1].FindBlock([]byte{})
		newBlockChain := append(blockChain, newBlock)
		if !newBlockChain.IsValid() { // TODO necessary?
			log.Fatal("Mined an invalid blockchain somehow.")
		}
		blockChain = bb.PossiblyReplace(blockChain, newBlockChain)
		adjustDifficulty(blockChain)
		ch <- blockChain
	}
}

func updateBlockchain(ch <-chan bb.BlockChain) {
	for bc := range ch {
		if !bc.IsValid() {
			log.Println("Received invalid blockchain.")
			continue
		}
		log.Println("Blockchain updated!")
		blockChain = bb.PossiblyReplace(blockChain, bc)
		adjustDifficulty(blockChain)
	}
}

func wsReader(wsconn *websocket.Conn, ch chan<- bb.BlockChain) {
	defer wsconn.Close()
	for {

		_, p, err := wsconn.ReadMessage()
		if err != nil {
			log.Fatalf("ReadMessage failed in wsReader: %v\n", err)
		}

		var buff bytes.Buffer
		_, err = buff.Write(p)
		if err != nil {
			log.Fatalf("Read failed in wsReader: %v\n", err)
		}

		var newBlockChain bb.BlockChain
		decoder := gob.NewDecoder(&buff)
		err = decoder.Decode(&newBlockChain)
		if err != nil {
			log.Fatal("decode error 1:", err)
		}

		fmt.Printf("Received blockchain: %s\n", newBlockChain.String())
		ch <- newBlockChain
	}
}

func wsWriter(ch <-chan bb.BlockChain) {
	for newBlockChain := range ch {
		log.Printf("wsWriter got the blockchain: %s\n", newBlockChain)
		var buf bytes.Buffer
		encoder := gob.NewEncoder(&buf)
		err := encoder.Encode(newBlockChain)
		if err != nil {
			log.Fatal("encode error:", err)
		}

		for _, ws := range wsconns {
			ws.WriteMessage(1, buf.Bytes())
			log.Printf("wrote to %s", ws.RemoteAddr().String())
		}
	}
}

func displayIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the Naivecoin http server")
}

func displayBlockchain(w http.ResponseWriter, r *http.Request) {
	for _, blk := range blockChain {
		fmt.Fprint(w, blk.String()+"\n")
	}
}

func parsePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprint(w, "please make a POST request.")
	}
	r.ParseForm()
	for k, v := range r.PostForm {
		fmt.Fprintf(w, "key is %s, val is %s \n", k, v)

		if k == "data" {
			blockChain = append(blockChain, blockChain[len(blockChain)-1].FindBlock([]byte(v[0])))
		}

		if k == "addpeer" {
			// We use this to manually add websocket peers, as there is no peer discovery mechanism.
			var err error
			var wsconn *websocket.Conn
			u := url.URL{Scheme: "ws", Host: v[0], Path: "/ws"}
			log.Printf("⧉ connecting to %s", u.String())
			wsconn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Fatal("dial:", err)
			}
			wsconns = append(wsconns, wsconn)
			go wsReader(wsconn, inCh)
		}
	}
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var wsconn *websocket.Conn
	wsconn, err = upgrader.Upgrade(w, r, nil)
	wsconns = append(wsconns, wsconn)
	if err != nil {
		log.Fatal(err)
	}
	go wsReader(wsconn, inCh)
	log.Printf("⧉ connection from %s", wsconn.RemoteAddr().String())

}
