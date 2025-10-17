package awsClient

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/amazon-s3-encryption-client-go/v3/client"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//Ce fichier permet de supprimer récursivement tous les objects de S3

// On commence par parcourir tous les buckets puis pour chaque bucket on supprime tous ses objets
func CleanS3(client *client.S3EncryptionClientV3) {
	listOut, err := client.ListBuckets(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	listBuckets := listOut.Buckets
	for _, bucket := range listBuckets {
		CleanS3Bucket(client, *aws.String(*bucket.Name))
	}
}

// On parcourt tous les objects d'un buckets et on les supprime
func CleanS3Bucket(client *client.S3EncryptionClientV3, bucketName string) error {
	//On liste tous les objets du buckets
	listInput := &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	}
	paginator := s3.NewListObjectsV2Paginator(client, listInput)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return fmt.Errorf("échec de la pagination : %w", err)
		}

		// Supprimer chaque objet
		for _, obj := range page.Contents {
			err := CleanS3Object(client, bucketName, *obj.Key)
			if err != nil {
				return fmt.Errorf("échec de la suppression de l'objet : %w", err)
			}
		}

	}
	_, err := client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: &bucketName,
	})
	if err != nil {
		return fmt.Errorf("échec de la suppression du bucket %s: %w", bucketName, err)
	}

	fmt.Printf("Le bucket %s a été supprimé avec succès\n", bucketName)
	return nil
}

func CleanS3Object(client *client.S3EncryptionClientV3, bucketName, objectKey string) error {
	// Créer une requête pour supprimer un objet
	input := &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	}

	// Supprimer l'objet
	_, err := client.DeleteObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("échec de la suppression de l'objet %s dans le bucket %s: %w", objectKey, bucketName, err)
	}

	fmt.Printf("L'objet %s a été supprimé avec succès du bucket %s\n", objectKey, bucketName)
	return nil
}