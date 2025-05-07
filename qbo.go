package main

import (
	"encoding/json"
	"fmt"
	"os"

	quickbooks "github.com/rwestlund/quickbooks-go"
)

func FillEstimate() quickbooks.Estimate {
	dat, err := os.ReadFile("./estimate.json")
	if err != nil {
		panic(err)
	}

	var estimate quickbooks.Estimate
	if err := json.Unmarshal(dat, &estimate); err != nil {
		panic(err)
	}
	return estimate

	// l := quickbooks.Line{
	// 	Amount:     "150.0",
	// 	DetailType: "foobaz",
	// }

	// estimate := quickbooks.Estimate{
	// 	TxnStatus: "foobar",
	// 	TotalAmt:  150.0,
	// 	Line:      []quickbooks.Line{l},
	// }
}

func LoadClient(token *quickbooks.BearerToken) (c *quickbooks.Client, err error) {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("SECRET")
	realmId := os.Getenv("REALM_ID")
	return quickbooks.NewClient(clientId, clientSecret, realmId, false, "", token)
}

func MakeRequest() string {
	// FIXME: load from DB:
	token := quickbooks.BearerToken{
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		AccessToken:  os.Getenv("ACCESS_TOKEN"),
	}

	client, err := LoadClient(&token)
	if err != nil {
		panic(err)
	}

	// To do first when you receive the authorization code from quickbooks callback
	// authorizationCode := "XAB11746551225hXNdSW2iGUcTdTLImx5gzNIF59QnhMmM40tX"
	// redirectURI := "https://developer.intuit.com/v2/OAuth2Playground/RedirectUrl"
	// bearerToken, err := client.RetrieveBearerToken(authorizationCode, redirectURI)
	// if err != nil {
	// 	panic(err)
	// }

	// When the token expire, you can use the following function
	_, err = client.RefreshToken(token.RefreshToken)
	if err != nil {
		panic(err)
	}

	// Make a request!
	info, err := client.FindCompanyInfo()
	if err != nil {
		panic(err)
	}
	fmt.Println(info)

	estimate := FillEstimate()
	estresp, err := client.CreateEstimate(&estimate)
	if err != nil {
		panic(err)
	}

	jsonBytes, err := json.MarshalIndent(estresp, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}
