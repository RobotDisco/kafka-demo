package main

import (
	"encoding/json"
	"fmt"

	"tulip.io/hackday/kafka"
)

func main() {
	for {
		msg := kafka.Generate()
		j, err := json.Marshal(msg)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(j))
	}
}
