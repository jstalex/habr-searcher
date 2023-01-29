package main

import hp "habr-searcher/internal/habr-parser"

func main() {
	fprser := hp.New("go")
	fprser.GetNewPost()
}
