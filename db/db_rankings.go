package db

import (
	"stockx-backend/db/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func PutRankingsInTheTable(item models.Rankings) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	tableName := "Rankings"

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

func GetRankingsFromTableForUser(date string) (models.Rankings, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	var result *dynamodb.GetItemOutput

	var err error

	tableName := "Rankings"
	region := "Europe/Copenhagen"

	result, err = svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"region": {
				S: aws.String(region),
			},
			"date": {
				S: aws.String(date),
			},
		},
	})
	if err != nil {
		return models.Rankings{}, err
	}

	if result.Item == nil {
		return models.Rankings{}, nil
	}

	item := models.Rankings{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return models.Rankings{}, err
	}

	return item, nil
}

func DeleteRankings(date string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	region := "Europe/Copenhagen"

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"region": {
				S: aws.String(region),
			},
			"date": {
				S: aws.String(date),
			},
		},
		TableName: aws.String("Rankings"),
	}

	_, err := svc.DeleteItem(input)

	return err
}
