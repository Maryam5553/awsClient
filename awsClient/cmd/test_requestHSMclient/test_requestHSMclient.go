package main

import (
	hsmClient "awsClient/pkg/requestHSMclient"
	"fmt"
)

// address of the HSM client
const HSM_CLIENT_ADDRESS = "localhost:8080"

func main() {
	// keys we want to retrieve on 2 different HSM
	keyHSM_1 := hsmClient.KeyHSM{
		Hsm_number: 17, // keystore key17
		Key_index:  1,  // key at index 1
	}
	keyHSM_2 := hsmClient.KeyHSM{
		Hsm_number: 22, // keystore key22
		Key_index:  1,  // key at index 1
	}

	// make parallel key requests
	key := hsmClient.GetKey(HSM_CLIENT_ADDRESS, keyHSM_1, keyHSM_2)

	// print result
	if len(key) == 0 {
		fmt.Println("Get key request failed.")
	} else {
		fmt.Println("key :")
		fmt.Println(key)
	}
}
