package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// AnyToJSON converts any type to JSON string.
// The any type is used instead of a specific struct type to allow for flexibility.
// Primiarly used for printing, debugging, testing purposes and possibly in future sending the data to a client.
func AnyToJSON(v any) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	os.Stdout.WriteString(string(jsonData))
}
