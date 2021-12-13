package db

import (
	"stockx-backend/db/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func PutListOfRegisteredUsersInTheTable(item models.ListOfRegisteredUsers) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	tableName := "ListOfRegisteredUsers"

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

func GetListOfRegisteredUsers() (models.ListOfRegisteredUsers, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	var result *dynamodb.GetItemOutput

	var err error

	tableName := "ListOfRegisteredUsers"

	result, err = svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"region": {
				S: aws.String("Europe/Copenhagen"),
			},
		},
	})
	if err != nil {
		return models.ListOfRegisteredUsers{}, err
	}

	if result.Item == nil {
		return models.ListOfRegisteredUsers{}, nil
	}

	item := models.ListOfRegisteredUsers{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return models.ListOfRegisteredUsers{}, err
	}

	return item, nil
}
