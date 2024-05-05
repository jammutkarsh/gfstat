package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// AnyToJSON converts any type to JSON string.
// Any type is used instead of a specific struct type to allow for flexibility.
// Primarily used for printing, debugging, testing purposes and possibly in future sending the data to a client.
func AnyToJSON(v any) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	_, err = os.Stdout.WriteString(string(jsonData))
	if err != nil {
		log.Println(err)
		return
	}
}
