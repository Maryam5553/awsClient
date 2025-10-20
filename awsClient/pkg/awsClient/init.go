package awsClient

import (
	"bufio"
	"fmt"
	"os"

	"github.com/aws/amazon-s3-encryption-client-go/v3/client"
)



// Le menu d'interface principal avec le client, lui proposant des actions
// (put, get, afficher l'arborescence...)
func InteractionConsole(client *client.S3EncryptionClientV3) int {
	fmt.Println()
	fmt.Println("Entrer une lettre pour effectuer une action:\nP = mettre un fichier sur Amazon S3\nG = récupérer un fichier\nL = lister les buckets présents sur S3\nA = afficher l'arborescence de fichiers sur S3\nD = supprimer tous les fichiers de amazon S3\nX = arrêter le programme")
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
	}
	//TODO : gérer le resultats des putObject et GetObject : en fait on en a pas besoin donc
	// la fonction pourrait juste retourner "err"
	if char == 'P' {
		_, err = traiterPut(client, reader) // cf put.go
	} else if char == 'G' {
		_, err = traiterGet(client, reader) // cf get.go
	} else if char == 'X' {
		return 0
	} else if char == 'L' {
		traiterList(client) // cf tools.go
	} else if char == 'A' {
		AfficherArborescence(client) // cf tree.go
	} else if char == 'D' {
		CleanS3(client) // cf clean.go
	} else {
		fmt.Println("Entrée non reconnue, veuillez recommencer")
	}
	if err != nil {
		fmt.Printf("erreur lors du traitement de la demande %v\n", err)
		fmt.Println("Une erreur a été engendrée et la requête n'a pas été effectuée")
	}
	return 1
}
