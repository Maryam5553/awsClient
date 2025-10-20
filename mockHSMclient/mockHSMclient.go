package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

/*
	Code for a mock HSM client, used to test the AWS client requests.
	This mock HSM client will listen for incoming requests and check if
	the first byte is the code for a "get key" request. If so, it will read
	the next to bytes (keystore, key index on the keystore), and return a hardcoded key
	after a random amount of milliseconds (<500 ms) to simulate the real behaviour of the HSM client.
*/

// mock HSM client port
const PORT = "6123"

const (
	GET_KEY_REQUEST_CODE byte = 0 // request code for the client HSM (get key)
	GET_KEY_SUCCESS_CODE byte = 0 // code returned by the HSM client if the key request was successful
)

func handleConnection(conn net.Conn) {
	fmt.Printf("Client %s connected.\n", conn.RemoteAddr().String())
	defer conn.Close()

	// read client requests
	buf := make([]byte, 128)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("read error : %v", err)
		return
	}

	// read first byte
	if buf[0] != GET_KEY_REQUEST_CODE {
		fmt.Printf("client %s asked for an unknown request (code %d)\n", conn.RemoteAddr().String(), buf[0])
		return
	}

	// read request data (for the mock HSM client)
	hsm_number := buf[1]
	key_index := buf[2]
	fmt.Printf("client %s asked to get key at HSM %d index %d\n", conn.RemoteAddr().String(), hsm_number, key_index)

	// sleeps a random amount of milliseconds to simulate the real HSM client behaviour
	n := rand.Intn(500) // 500 ms
	time.Sleep(time.Duration(n) * time.Millisecond)

	// sends the key (here, it's just a 16 bytes hardcoded key)
	_, err = conn.Write([]byte{
		GET_KEY_SUCCESS_CODE,
		0x01, 0x02, 0x03, 0x04,
		0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0b, 0x0c,
		0x0d, 0x0e, 0x0f, 0x10})
	if err != nil {
		fmt.Printf("error sending key to client %s: %v\n", conn.RemoteAddr().String(), err)
	}
	fmt.Printf("key (HSM %d, index %d) sent to client %s.\n", hsm_number, key_index, conn.RemoteAddr().String())
}

func RunMockHSMclient(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	fmt.Printf("Mock HSM client listening at address localhost:%s...\n", port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("client failed to connect")
		}
		go handleConnection(conn)
	}
}

func main() {
	RunMockHSMclient(PORT)
}
