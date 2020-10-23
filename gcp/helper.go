package gcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
)

type apiErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Check HTTP response code, return error if HTTP request is not successed.
func apiResponseChecker(statusCode int, response []byte, funcName string) error {

	if statusCode >= 300 || statusCode < 200 {
		log.Printf("%s request failed", funcName)
		var errorResponse apiErrorResponse
		responseContent := bytes.NewBuffer(response).String()
		if err := json.Unmarshal(response, &errorResponse); err != nil {
			log.Printf("Failed to unmarshall error response from %s", funcName)
			return fmt.Errorf(responseContent)
		}
		return fmt.Errorf("code: %d, message: %s", errorResponse.Code, errorResponse.Message)
	}

	return nil

}

func nextRandomInt(min int, max int) int {
	return rand.Intn(max-min) + min
}
