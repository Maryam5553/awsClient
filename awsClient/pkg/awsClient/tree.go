package awsClient

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/amazon-s3-encryption-client-go/v3/client"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//En fait le problème venait du fait que même si S3 afficher les buckets sous forme d'arborescence, elle les stocke en format liste
//Il fallait donc faire la conversion

type Node struct {
	Name     string
	Children map[string]*Node
	IsFile   bool
}

// Cette fonction permet de remplir un arbre en y ajoutant un noeud (qui peut être un dossier ou un fichier)
func addNode(root *Node, key string) {
	parts := strings.Split(key, "/")
	currentNode := root

	for i, part := range parts {
		if part == "" {
			continue
		}
		if currentNode.Children[part] == nil {
			currentNode.Children[part] = &Node{Name: part, Children: make(map[string]*Node)}
		}
		if i == len(parts)-1 {
			currentNode.Children[part].IsFile = true
		}
		currentNode = currentNode.Children[part]
	}
}

// Affiche un Tree sous forme arborescente
func printTree(node *Node, indent string) {
	if node.IsFile {
		fmt.Println(indent + node.Name)
	} else {
		fmt.Println(indent + node.Name + "/")
	}
	for _, child := range node.Children {
		printTree(child, indent+"  ")
	}
}

// Cette fonction permet d'afficher tous les sous dossiers d'un bucket
func auxArbo(client *client.S3EncryptionClientV3, bucket string) *Node {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}
	root := &Node{Name: bucket, Children: make(map[string]*Node)}
	paginator := s3.NewListObjectsV2Paginator(client, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to list objects, %v", err)
		}
		for _, object := range output.Contents {
			addNode(root, *object.Key)
		}
	}
	return root
}

// On affiche chaque bucket via les deux fonctions précédentes
func AfficherArborescence(client *client.S3EncryptionClientV3) {
	listOut, err := client.ListBuckets(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	listBuckets := listOut.Buckets
	//fmt.Printf("%d\n", len(listBuckets))
	for _, bucket := range listBuckets {
		//fmt.Println(*bucket.Name)
		printTree(auxArbo(client, *bucket.Name), "")
	}
}

// Recherche récursive d'un chemin dans l'arbre
func inTree(path string, root *Node) (bool, *Node) {
	//fmt.Printf("path = %s\n", path)
	parts := strings.Split(path, "/")
	//fmt.Printf("%s\n", root.Name)
	if len(parts) == 0 {
		return false, nil
	} else {
		if parts[0] == root.Name {
			//fmt.Printf("cas 1\n")
			if len(parts) == 1 {
				//fmt.Println("cas 2")
				return true, root
			} else {
				//fmt.Println("cas 3")
				newPath := strings.Join(parts[1:], "/")
				//fmt.Printf("newPath = %s\n", newPath)
				for _, value := range root.Children {
					tmp, newRoot := inTree(newPath, value)
					if tmp {
						return true, newRoot
					}
				}
			}
		}
	}
	//fmt.Println("cas 4")
	return false, nil

}