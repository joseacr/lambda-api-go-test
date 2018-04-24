package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "regexp"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

var isbnRegexp = regexp.MustCompile(`[0-9]{3}\-[0-9]{10}`)
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

type product struct {
    Name   string `json:"name"`
    MinTemp  string `json:"mintemp"`
    MaxTemp string `json:"maxtemp"`
    User string`json:"user"`
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    switch req.HTTPMethod {
    case "GET":
        return show(req)
    case "POST":
        return create(req)
    default:
        return clientError(http.StatusMethodNotAllowed)
    }
}

func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    name := req.QueryStringParameters["name"]
    if !isbnRegexp.MatchString(name) {
        return clientError(http.StatusBadRequest)
    }

    p, err := getItem(name)
    if err != nil {
        return serverError(err)
    }
    if p == nil {
        return clientError(http.StatusNotFound)
    }

    js, err := json.Marshal(p)
    if err != nil {
        return serverError(err)
    }

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       string(js),
    }, nil
}

func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    if req.Headers["Content-Type"] != "application/json" {
        return clientError(http.StatusNotAcceptable)
    }

    p := new(product)
    err := json.Unmarshal([]byte(req.Body), p)
    if err != nil {
        return clientError(http.StatusUnprocessableEntity)
    }

    if !isbnRegexp.MatchString(p.Name) {
        return clientError(http.StatusBadRequest)
    }
    if p.MinTemp == "" || p.MaxTemp == "" {
        return clientError(http.StatusBadRequest)
    }

    err = putItem(p)
    if err != nil {
        return serverError(err)
    }

    return events.APIGatewayProxyResponse{
        StatusCode: 201,
        Headers:    map[string]string{"Location": fmt.Sprintf("/product?name=%s", p.Name)},
    }, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
    errorLogger.Println(err.Error())

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusInternalServerError,
        Body:       http.StatusText(http.StatusInternalServerError),
    }, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse{
        StatusCode: status,
        Body:       http.StatusText(status),
    }, nil
}

func main() {
    lambda.Start(router)
}