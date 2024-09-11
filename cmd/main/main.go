package main

import (
	"fmt"
	"github.com/victorspinei/urban-disco/internal/api"
)

func main() {
	fmt.Println("Hello, world!")
	api.GetSongsFromAlbum("seventh_son_of_a_seventh_son")
}