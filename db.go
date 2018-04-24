package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("eu-west-1"))
/*
type AddressInfo struct {
    Street string`json:"street"`
    City string`json:"city"`
    State string`json:"state"`
    Country string`json:"country"`
}

type UserInfo struct {
    Id int`json:"id"`
    Name string`json:"name"`
    Address AddressInfo`json:"address"`
}

type Product struct {
    Name string`json:"name"`
    MinTemp int`json:"mintemp"`
    MaxTemp int`json:"maxtemp"`
    User UserInfo`json:"user"`
}
*/

func getItem(name string) (*product, error) {
    input := &dynamodb.GetItemInput{
        TableName: aws.String("TrustForTags_Product"),
        Key: map[string]*dynamodb.AttributeValue{
            "Name": {
                S: aws.String(name),
            },
        },
    }

    result, err := db.GetItem(input)
    if err != nil {
        return nil, err
    }
    if result.Item == nil {
        return nil, nil
    }

    p := new(product)
    err = dynamodbattribute.UnmarshalMap(result.Item, p)
    if err != nil {
        return nil, err
    }

    return p, nil
}

// Add a product record to DynamoDB.
func putItem(p *product) error {

    input := &dynamodb.PutItemInput{
        TableName: aws.String("TrustForTags_Product"),
        Item: map[string]*dynamodb.AttributeValue{
            "Name": {
                S: aws.String(p.name),
            },
            "MinTemp": {
                N: aws.String(p.minTemp),
            },
            "MaxTemp": {
                N: aws.String(p.maxTemp),
            },
            "User": {
                S: aws.String(p.user),
            },
        },
    }

    _, err := db.PutItem(input)
    return err
}