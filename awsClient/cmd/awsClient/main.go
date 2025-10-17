package main

import (
	"context"
	"fmt"
	"log"

	"awsClient/pkg/awsClient"
	MyMaterials "awsClient/pkg/awsEncryptionMaterials"
	hsmClient "awsClient/pkg/requestHSMclient"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/amazon-s3-encryption-client-go/v3/client"
)

/*
	This program implements a AWS client that can put and get files from AWS S3
	with a client side encryption with key retrieved by a HSM client of
	address given in the program.
*/

const HSM_CLIENT_ADDRESS = "localhost:8080"

// retourne un client s3 Localstack (j'ai harcodé les infos localstack mais
// il faudrait que je crée un client qui accepte de vrais id AWS proprement)
func CreateS3Client() (*s3.Client, error) {
	// infos localstack
	awsRegion := "us-east-1"
	awsEndpoint := "http://localhost:4566"

	// créer la configuration avec la région
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot load the AWS configs: %w", err)
	}

	// créer le client S3 avec la configuration précédente et l'endpoint
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(awsEndpoint)
	})
	return client, nil
}

// retourne un S3 encryption client
func CreateS3EncryptionClient(hsm_client_address string, keyHSM_1 hsmClient.KeyHSM, keyHSM_2 hsmClient.KeyHSM) (*client.S3EncryptionClientV3, error) {
	s3Client, err := CreateS3Client()
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
	s3EncryptionClient, err := CreateS3EncryptionClient(HSM_CLIENT_ADDRESS, keyHSM_1, keyHSM_2)
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
