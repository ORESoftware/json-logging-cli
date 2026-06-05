package main

import (
	"encoding/json"
	"fmt"
	"os"

	logging "github.com/oresoftware/json-logging/jlog/helper"
)

func main() {
	var v any

	if err := json.NewDecoder(os.Stdin).Decode(&v); err != nil {
		panic(err)
	}

	fmt.Println(logging.GetPrettyString(v, 0))
}
