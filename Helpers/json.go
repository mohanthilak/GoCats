package helpers

import (
	"encoding/json"
	"log"
)

func ConvertToJSON(v any) []byte {
	result, err := json.Marshal(v)
	if err != nil {
		log.Println("error while Marshalling to JSON", err)
		return nil
	}
	return result
}

func ConvertFromJSON(data []byte) any {
	var V any
	if err := json.Unmarshal(data, V); err != nil {
		log.Println("error while unmarshalling JSON. Data: ", data, err)
		return nil
	}
	return V
}
