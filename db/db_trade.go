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

func PutTradesInTheTable(email string, item models.Trades) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	tableName := "Trades"

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

func GetTradesFromTableForUser(email string) (models.Trades, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	var result *dynamodb.GetItemOutput

	var err error

	tableName := "Trades"

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
			return models.Trades{}, err
		}
	} else {
		return models.Trades{}, reserr.NotFound("invalid email", err, "email did not pass regex")
	}

	if result.Item == nil {
		return models.Trades{}, nil
	}

	item := models.Trades{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return models.Trades{}, err
	}

	if len(item.BoughtStocks) == 0 {
		item.BoughtStocks = []models.BoughtStock{}
	}

	if len(item.SoldStocks) == 0 {
		item.SoldStocks = []models.SoldStock{}
	}

	if len(item.ShortStocks) == 0 {
		item.ShortStocks = []models.ShortStock{}
	}

	if len(item.BoughtToCover) == 0 {
		item.BoughtToCover = []models.BoughtToCover{}
	}

	if len(item.HoldLong) == 0 {
		item.HoldLong = []models.HoldLong{}
	}

	if len(item.HoldShort) == 0 {
		item.HoldShort = []models.HoldShort{}
	}

	return item, nil
}

func DeleteTrades(email string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String("Trades"),
	}

	_, err := svc.DeleteItem(input)

	return err
}
