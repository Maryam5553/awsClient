package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"

	"awsClient/pkg/awsClient"
	MyMaterials "awsClient/pkg/awsEncryptionMaterials"
	hsmClient "awsClient/pkg/requestHSMclient"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/amazon-s3-encryption-client-go/v3/client"
)

/*
	This program implements a AWS client that can put and get files from AWS S3
	with a client side encryption with key retrieved by a HSM client of
	address given in the program.
*/

const HSM_CLIENT_DEFAULT_PORT = 6123
const AWS_CREDENTIALS_PATH = ".aws/credentials"
const AWS_CONFIG_PATH = ".aws/config"
const LOCALSTACK_ENDPOINT = "http://localhost:4566"

// retourne un client s3 Localstack
func CreateS3Client_Locastack() (*s3.Client, error) {
	// infos localstack
	awsRegion := "us-east-1"
	awsEndpoint := LOCALSTACK_ENDPOINT

	// créer la configuration avec la région
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion), config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("test", "test", ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot load the AWS configs: %w", err)
	}

	// créer le client S3 avec la configuration précédente et l'endpoint
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(awsEndpoint)
	})

	fmt.Printf("S3 client connected to LocalStack %s\n", awsEndpoint)
	return client, nil
}

// retourne un client s3 relié à un compte AWS
func CreateS3Client_AWS() (*s3.Client, error) {
	// crée la configuration selon les informations des fichiers credentials et config
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedCredentialsFiles(
			[]string{AWS_CREDENTIALS_PATH},
		),
		config.WithSharedConfigFiles(
			[]string{AWS_CONFIG_PATH},
		))
	if err != nil {
		return nil, fmt.Errorf("cannot load the AWS configs: %w", err)
	}

	// créer le client S3 avec la configuration précédente
	client := s3.NewFromConfig(awsCfg)

	fmt.Printf("S3 client connected to AWS account defined in %s\n", AWS_CREDENTIALS_PATH)
	return client, nil
}

// retourne un client s3. si localstack==true, le client est lié à l'endpoint localstack.
// sinon, il est lié à un vrai compte AWS
func CreateS3Client(localstack bool) (*s3.Client, error) {
	if localstack {
		return CreateS3Client_Locastack()
	} else {
		return CreateS3Client_AWS()
	}
}

// retourne un S3 encryption client
func CreateS3EncryptionClient(hsm_client_address string, keyHSM_1 hsmClient.KeyHSM, keyHSM_2 hsmClient.KeyHSM, localstack bool) (*client.S3EncryptionClientV3, error) {
	s3Client, err := CreateS3Client(localstack)
	if err != nil {
		return nil, fmt.Errorf("couldn't create S3 client: %v", err)
	}
	cmm := MyMaterials.NewCustomCryptographicMaterialsManager(hsm_client_address, keyHSM_1, keyHSM_2)
	encryptionClient, err := client.New(s3Client, cmm)
	if err != nil {
		return nil, fmt.Errorf("couldn't create encryption client: %v", err)
	}
	return encryptionClient, nil
}

func main() {
	// command-line arguments
	hsm_client_port_flag := flag.Int("HSMclient", HSM_CLIENT_DEFAULT_PORT, "HSM client port")
	localstack_flag := flag.Bool("localstack", false, "if true, the AWS client will connect to LocalStack. Otherwise (default behaviour), it will connect to a AWS account")

	flag.Parse()
	HSM_CLIENT_ADDRESS := "localhost:" + strconv.Itoa(*hsm_client_port_flag)
	fmt.Printf("HSM client address %s", HSM_CLIENT_ADDRESS)

	// keys we want to retrieve on 2 different HSM
	keyHSM_1 := hsmClient.KeyHSM{
		Hsm_number: 17, // keystore key17
		Key_index:  1,  // key at index 1
	}
	keyHSM_2 := hsmClient.KeyHSM{
		Hsm_number: 22, // keystore key22
		Key_index:  1,  // key at index 1
	}

	// créer le S3 encryption client avec les informations cryptographiques ci-dessus
	s3EncryptionClient, err := CreateS3EncryptionClient(HSM_CLIENT_ADDRESS, keyHSM_1, keyHSM_2, *localstack_flag)
	if err != nil {
		log.Fatal("error creating encryption client")
	}

	// Une fois le mode de chiffrement décidé, on peut demander à l'utilisateur
	// ce qu'il veut faire comme actions. cf fichier init.go
	tmp := 1
	for tmp != 0 {
		tmp = awsClient.InteractionConsole(s3EncryptionClient) // fonction dans init.go
	}
}
