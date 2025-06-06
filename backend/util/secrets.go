package util

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../../.env")
}

type Secrets struct {
	JWTSecretKey string `json:"JWT_SECRET_KEY"`
	DBPassword   string `json:"DB_PASSWORD"`
}

func getSecrets() (Secrets, error) {
	var secrets Secrets

	secretsManagerName := os.Getenv("SECRETS_MANAGER_NAME")
	if len(secretsManagerName) == 0 {
		err := errors.New("AWS secrets manager name is not specified")
		return Secrets{}, err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return Secrets{}, err
	}

	svc := secretsmanager.NewFromConfig(cfg)

	result, err := svc.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretsManagerName),
	})
	if err != nil {
		return Secrets{}, err
	}

	err = json.Unmarshal([]byte(*result.SecretString), &secrets)
	if err != nil {
		return Secrets{}, err
	}

	return secrets, nil
}
