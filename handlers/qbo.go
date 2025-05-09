package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/brtholomy/templui-quickstart/ui/pages"
	quickbooks "github.com/rwestlund/quickbooks-go"
)

func NewQboHandler(request func(*http.Request) string) QboHandler {
	return QboHandler{Request: request}
}

type QboHandler struct {
	// must match the signature of the QboRequest func below.
	Request func(r *http.Request) string
}

// implements the HTTP handler interface on the QboHandler type.
func (qh QboHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// NOTE: can pass in http.Request fields and methods.
	resp := qh.Request(r)
	pages.Qbo(resp).Render(r.Context(), w)
}

func FillEstimate() quickbooks.Estimate {
	var estimate quickbooks.Estimate
	if err := json.Unmarshal([]byte(ESTIMATE), &estimate); err != nil {
		panic(err)
	}
	estimate.CustomerRef.Value = "2"
	estimate.Line[0].Description = "BTH RULES"
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

func QboRequest(req *http.Request) string {
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

	// TODO: figure out how often to refresh?
	_, err = client.RefreshToken(token.RefreshToken)
	if err != nil {
		panic(err)
	}

	// make a more interesting request:
	estimate := FillEstimate()
	estresp, err := client.CreateEstimate(&estimate)
	if err != nil {
		panic(err)
	}
	jsonBytes, err := json.MarshalIndent(estresp, "", "  ")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%v\n%v\n%v\n", req.URL.Path, req.URL.RawQuery, string(jsonBytes))
}
