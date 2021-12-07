package db

import (
	"log"
	"net/mail"
	"stockx-backend/db/models"
	"stockx-backend/reserr"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// This file contains functions for any DB interactions

func PutItemInTable(item interface{}, tableName string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

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

func UpdateUsersPassword(useremail, hashedPass string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// Update item in table Movies

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c": {
				S: aws.String(hashedPass),
			},
		},
		TableName: aws.String("User"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(useremail),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set password = :r"),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		log.Fatalf("Got error calling UpdateItem: %s", err)
	}

	return nil
}

func GetUserFromTable(email string) (models.User, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	var result *dynamodb.GetItemOutput

	var err error

	tableName := "User"

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
			return models.User{}, err
		}
	} else {
		return models.User{}, reserr.NotFound("invalid email", err, "email did not pass regex")
	}

	if result.Item == nil {
		return models.User{}, nil
	}

	item := models.User{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return models.User{}, err
	}

	return item, nil
}
