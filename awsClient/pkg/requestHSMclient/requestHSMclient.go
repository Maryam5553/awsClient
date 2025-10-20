package requestHSMclient

import (
	"fmt"
	"net"
)

const (
	GET_KEY_REQUEST_CODE byte = 0  // request code for the client HSM (get key)
	GET_KEY_SUCCESS_CODE byte = 0  // code returned by the HSM client if the key request was successfull
	ANSWER_BUF_SIZE      int  = 17 // size of the expected answer from the HSM client (success byte + 16 key bytes)
)

// structure to represent a key on a given HSM.
// at the moment we'll use HSM key17 and key22 (hsm_number: 17 or 22)
// and there are 32 indexes on each HSM.
type KeyHSM struct {
	Hsm_number int
	Key_index  int
}

// connect to HSM client via a TCP socket.
// you have to give the HSM client address in parameter
// ex: ConnectHSMClient("localhost:8080")
func ConnectHSMClient(hsm_client_addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", hsm_client_addr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to HSM client %s: %w", hsm_client_addr, err)
	}
	return conn, nil
}

// creates a message to send to HSM client
// to request the key on a given index on a given HSM.
// (we can change this function if we want to adapt request format)
func MakeKeyRequestMessage(keyHSM KeyHSM) []byte {
	request := []byte{GET_KEY_REQUEST_CODE, byte(keyHSM.Hsm_number), byte(keyHSM.Key_index)}
	return request
}

// send a request on the open connexion with the HSM client
// to retrieve key at the index given in parameter.
// returns the key bytes and an error.
func SendKeyRequest(conn net.Conn, keyHSM KeyHSM) ([]byte, error) {
	// create request message
	request := MakeKeyRequestMessage(keyHSM)

	// send request to HSM client
	_, err := conn.Write(request)
	if err != nil {
		return []byte{}, fmt.Errorf("error sending key request to HSM %d at index %d: %w", keyHSM.Hsm_number, keyHSM.Key_index, err)
	}

	// wait for HSM client answer
	buf := make([]byte, ANSWER_BUF_SIZE)
	_, err = conn.Read(buf)
	if err != nil {
		return []byte{}, fmt.Errorf("error reading key request (index %d) answer from HSM %d: %w", keyHSM.Key_index, keyHSM.Hsm_number, err)
	}

	// analyse HSM client answer code
	if buf[0] != GET_KEY_SUCCESS_CODE {
		return []byte{}, fmt.Errorf("the HSM client returned that the key request at HSM %d index %d failed", keyHSM.Hsm_number, keyHSM.Key_index)
	} else {
		// if the first byte indicates a success, we return the 16 following bytes containing the key
		return buf[1:17], nil
	}
}

// result type for the function GetKetFromHSM.
// we're using a structure to better handle the result from the goroutine
// (as we'll exectute the function in multiple goroutines)
type resGetKey struct {
	key []byte
	err error
}

// retrieve a key from a given HSM.
// returns the key in string format and an error.
func GetKeyFromHSM(hsm_client_addr string, keyHSM KeyHSM) resGetKey {
	// opens a connexion to the HSM client, that will interact with the HSM
	conn, err := ConnectHSMClient(hsm_client_addr)
	if err != nil {
		return resGetKey{key: []byte{}, err: fmt.Errorf("error sending request to HSM client: %v", err)}
	}
	defer conn.Close()

	// send a key request to HSM client to retrieve key at a given index on the given HSM
	key, err := SendKeyRequest(conn, keyHSM)
	if err != nil {
		return resGetKey{key: []byte{}, err: fmt.Errorf("error sending request to HSM client: %v", err)}
	}
	return resGetKey{key: key, err: nil}
}

// sends 2 parallel requests to 2 HSM to retrieve a key at a given index.
// parameters : HSM client address, 2 keys (reference by their HSM and index)
// returns the key or an empty byte slice if the request failed for both goroutines
// (in this case, the error will be printed)
func GetKey(hsm_client_addr string, keyHSM_1 KeyHSM, keyHSM_2 KeyHSM) []byte {
	// channel to retrieve HSM request results (key + eventual error)
	return_values := make(chan resGetKey, 2)

	// make 2 parallel requests
	go func() {
		return_values <- GetKeyFromHSM(hsm_client_addr, keyHSM_1)
	}()
	go func() {
		return_values <- GetKeyFromHSM(hsm_client_addr, keyHSM_2)
	}()

	key := []byte{}
	// return values
	for range 2 {
		res := <-return_values
		// if the first goroutine to finish returns an error,
		// print the error and continue
		if res.err != nil {
			// don't print this error as if one request succeed we can't ignore the other error
			// we'll only print it if both requests fail
			// log.Println(res.err)
			continue
		} else {
			// else : a key was returned
			key = res.key
			break
		}
	}
	return key
}
