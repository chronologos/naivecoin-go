# naivecoin
Simple cryptocurrency implementation in Go. Loosely based on javascript [naivecoin](https://lhartikk.github.io/jekyll/update/2017/07/14/chapter1.html) tutorial.

## Sample CURL commands
```
curl -X POST -d "addpeer=localhost:9000" -H 'Content-Type: application/x-www-form-urlencoded' 'localhost:8000/p'
```

Add a peer


```
curl -X POST -d "data=bob" -H 'Content-Type: application/x-www-form-urlencoded' 'localhost:8000/p'
```

Mine a block with data = "bob"

## TODO
1. Find out are coinbase tx inserted into the blockchain and when `validateCoinbaseTx` is called.
