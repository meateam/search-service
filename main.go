package main

import (
	"github.com/meateam/search-service/server"
)

func main() {
	server.NewServer(nil).Serve(nil)
}
