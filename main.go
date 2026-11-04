package main

import (
	"fmt"

	"github.com/queueunit/Bloggo/internal/config"
)

func main() {
	x := config.Read()
	fmt.Println(x.Db_url)
	x.SetUser("Queue")

	return
}
