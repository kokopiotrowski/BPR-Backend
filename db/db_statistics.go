package db

import (
	"net/mail"
	"stockx-backend/db/models"
	"stockx-backend/reserr"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func PutStatisticsInTheTable(email string, item models.Statistics) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	tableName := "Statistics"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func GetStatisticsFromTableForUser(email string) (models.Statistics, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	var result *dynamodb.GetItemOutput

	var err error

	tableName := "Statistics"

	if _, err = mail.ParseAddress(email); err == nil {
		result, err = svc.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"email": {
					S: aws.String(email),
				},
			},
		})
		if err != nil {
			return models.Statistics{}, err
		}
	} else {
		return models.Statistics{}, reserr.NotFound("invalid email", err, "email did not pass regex")
	}

	if result.Item == nil {
		return models.Statistics{}, nil
	}

	item := models.Statistics{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return models.Statistics{}, err
	}

	return item, nil
}
