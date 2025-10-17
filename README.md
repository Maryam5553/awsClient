# AWS Client

Implémentation d'un client AWS qui propose à l'utilisateur de mettre et récupérer des fichiers sur S3.

Il utilise le S3 encryption client pour chiffrer les fichiers côté client. Celui-ci fait des requêtes à un client HSM pour récupérer les clés voulues sur le HSM.

Pour tester le client AWS, il faut que le client HSM soit en train de tourner. On peut lancer le programme mockHSMclient.go pour simuler un client HSM qui traite les requêtes selon le format demandé.