package awsClient

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/aws/amazon-s3-encryption-client-go/v3/client"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Fonctions générales pour intéragir avec S3, utilisées par nos autres fichiers.go

// Cette fonction permet de regarder si un fichier est présent déja dans un bucket
func inAWSS3(client *client.S3EncryptionClientV3, bucket, path string) (bool, *Node) {
	root := auxArbo(client, bucket)
	estPresent, newRoot := inTree(bucket+"/"+path, root)
	return estPresent, newRoot
}

// Regarde si un dossier est présent et si ce n'est pas le cas propose de le créer
func verifierDossier(client *client.S3EncryptionClientV3, bucket, sous_rep string) (int, error) {
	estPresent, _ := inAWSS3(client, bucket, sous_rep)
	if estPresent {
		return 2, nil
	} else {
		fmt.Println("le dossier où vous voulez mettre le fichier n'existe pas, voulez-vous en créer un avec le nom que vous avez donné ? (O/N)")
		reader := bufio.NewReader(os.Stdin)
		char, _, err := reader.ReadRune()
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
		}
		//TODO : gérer l'erreur
		if char == 'O' {
			return 1, err
		} else if char == 'N' {
			return 0, nil
		} else {
			fmt.Println("réponse invalide- Par défaut, l'action n'est pas effectuée")
			return 0, err
		}
	}
}

// Regarde si un fichier (associé à un chemin) est présent, si c'est déja le cas, propose de le remplacer
func verifierKey(client *client.S3EncryptionClientV3, bucket, key string) int {
	estPresent, _ := inAWSS3(client, bucket, key)
	if estPresent {
		fmt.Printf("Il y a déja un fichier avec ce nom dans le bucket %s, l'action de Put remplacera entièrement le fichier (ou modifira le dossier)", bucket)
		fmt.Println("Voulez-vous continuer ? (O/N)")
		reader := bufio.NewReader(os.Stdin)
		char, _, err := reader.ReadRune()
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
		}
		if char == 'O' {
			return 1
		} else if char == 'N' {
			return 0
		} else {
			fmt.Println("réponse invalide- Par défaut, l'action n'est pas effectuée")
			return 0
		}
	}
	return 1
}

// Vérifier si un bucket est présent ou pas
func bucketPresent(client *client.S3EncryptionClientV3, bucket string) bool {
	listOut, _ := client.ListBuckets(context.TODO(), nil)
	listBuckets := listOut.Buckets
	var estPresent bool = false
	for _, buck := range listBuckets {
		if *buck.Name == bucket {
			estPresent = true
		}
	}
	return estPresent
}

// Cette fonction vérifie que le bucket est présent, si ce n'est pas le cas, demande à l'utilisateur s'il veut crée le bucket
func verifierBucket(client *client.S3EncryptionClientV3, bucket string) (int, error) {
	estPresent := bucketPresent(client, bucket)
	if estPresent {
		return 2, nil
	} else {
		fmt.Println("le bucket où vous voulez mettre le fichier n'existe pas, voulez-vous en créer un avec le nom que vous avez donné ? (O/N)")
		reader := bufio.NewReader(os.Stdin)
		char, _, err := reader.ReadRune()
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
		}
		//TODO : gérer l'erreur
		if char == 'O' {
			//TODO : proposer la création automatique de bucket -
			//TODO2 : vérifier l'existence d'un bucket
			fmt.Println("Petites précisions quant à la création de Buckets :")
			fmt.Println("1)Votre nom de bucket doit être unique- Autrement dit, aucune autre personne dans le monde doit avoir un bucket ayant le même nom")
			fmt.Println("2)Votre bucket doit faire entre 3 et 63 caractères")
			fmt.Println("Seules les lettres minuscules les chiffres et les tirets \"-\" sont autorisés")
			_, err = client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
				Bucket: aws.String(bucket),
				CreateBucketConfiguration: &types.CreateBucketConfiguration{
					LocationConstraint: "eu-west-3",
				},
			})
			return 1, err
		} else if char == 'N' {
			return 0, nil
		} else {
			fmt.Println("réponse invalide- Par défaut, l'action n'est pas effectuée")
			return 0, err
		}
	}
}

// Fonction pour afficher tous les buckets de l'utilisateur
func traiterList(client *client.S3EncryptionClientV3) {
	listOut, _ := client.ListBuckets(context.TODO(), nil)
	listBuckets := listOut.Buckets
	for _, bucket := range listBuckets {
		fmt.Println(*bucket.Name)
	}
}