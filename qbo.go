package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/brtholomy/templui-quickstart/ui/pages"
	quickbooks "github.com/rwestlund/quickbooks-go"
)

var estimatestr string = `{
  "TotalAmt": 31.5,
  "Line": [
    {
      "Description": "Pest Control Services",
      "DetailType": "SalesItemLineDetail",
      "SalesItemLineDetail": {
        "TaxCodeRef": {
          "value": "NON"
        },
        "Qty": 1,
        "UnitPrice": 35,
        "ItemRef": {
          "name": "Pest Control",
          "value": "10"
        }
      },
      "LineNum": 1,
      "Amount": 35.0,
      "Id": "1"
    },
    {
      "DetailType": "SubTotalLineDetail",
      "Amount": 35.0,
      "SubTotalLineDetail": {}
    },
    {
      "DetailType": "DiscountLineDetail",
      "Amount": 3.5,
      "DiscountLineDetail": {
        "DiscountAccountRef": {
          "name": "Discounts given",
          "value": "86"
        },
        "PercentBased": true,
        "DiscountPercent": 10
      }
    }
  ],
  "CustomerRef": {
    "name": "Cool Cars",
    "value": "3"
  },
  "TxnTaxDetail": {
    "TotalTax": 0
  },
  "ApplyTaxAfterDiscount": false
}`

func NewQboHandler(request func(string) string) QboHandler {
	return QboHandler{Request: request}
}

type QboHandler struct {
	// must match the signature of the QboRequest func below.
	Request func(string) string
}

// implements the HTTP handler interface on the QboHandler type.
func (qh QboHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// NOTE: can pass in http.Request fields and methods.
	pages.Qbo(qh.Request(r.UserAgent())).Render(r.Context(), w)
}

func FillEstimate() quickbooks.Estimate {
	// dat, err := os.ReadFile("./estimate.json")
	// if err != nil {
	// 	panic(err)
	// }

	var estimate quickbooks.Estimate
	if err := json.Unmarshal([]byte(estimatestr), &estimate); err != nil {
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
	fmt.Println(clientId, clientSecret, realmId)
	return quickbooks.NewClient(clientId, clientSecret, realmId, false, "", token)
}

func QboRequest(agent string) string {
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
	return agent + string(jsonBytes)
}
