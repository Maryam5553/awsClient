# Mock HSM Client
A mock HSM client, used to test the AWS client requests.

This mock HSM client will listen for incoming requests and check if the first byte is the code for a "get key" request. If so, it will read the next to bytes (keystore, key index on the keystore), and return a hardcoded key after a random amount of milliseconds (<500 ms) to simulate the real behaviour of the HSM client.

## Prerequisites

Golang

## Run

By default the server will listen on port 6123.

```go run mockHSMclient.go```

You can then run the AWS Client to make requests.