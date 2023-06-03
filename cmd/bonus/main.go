package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hasura/go-graphql-client"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type myProduct struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	BonusData GQLQuery
}

type slackMessageBlock struct {
	Type string `json:"type"`
	Text struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"text"`
	Accessory struct {
		Type     string `json:"type"`
		ImageUrl string `json:"image_url"`
		AltText  string `json:"alt_text"`
	} `json:"accessory"`
}

type slackMessage struct {
	Channel string              `json:"channel"`
	Blocks  []slackMessageBlock `json:"blocks"`
}

type GQLQuery struct {
	Product struct {
		Id         int
		Title      string
		WebPath    string
		SmartLabel string
		Price      struct {
			Now struct {
				Amount float64
			}
			Was struct {
				Amount float64
			}
			UnitInfo struct {
				Price struct {
					Amount float64
				}
				Description string
			}
			Discount struct {
				SegmentId   int
				Description string
			}
		}
		Images []struct {
			Width    int
			Height   int
			Url      string
			TypeName string `graphql:"__typename"`
		}
	} `graphql:"product(id: $id, date: $date)"`
}

type Query struct {
	Query GQLQuery `graphql:"product($id: Int!, $date: String)"`
}

func main() {
	fmt.Println("Go Bonus")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	myProducts, err := parseMyProducts()
	if err != nil {
		log.Fatal(err)
	}

	err = processMyProducts(myProducts)
	if err != nil {
		log.Fatal(err)
	}
}

func parseMyProducts() ([]myProduct, error) {
	contents, err := os.ReadFile(os.Getenv("PRODUCTS_JSON_FILE"))
	if err != nil {
		return nil, err
	}

	products := new([]myProduct)
	err = json.Unmarshal(contents, products)
	if err != nil {
		return nil, err
	}

	return *products, nil
}

func processMyProducts(myProducts []myProduct) error {

	var bonusProducts []myProduct

	for _, product := range myProducts {
		response, err := checkProduct(product)
		if err != nil {
			return err
		}

		hasDiscount := response.Product.Price.Discount.SegmentId != 0
		if hasDiscount {
			bonusProducts = append(bonusProducts, myProduct{
				Name:      product.Name,
				BonusData: response,
			})
		}
	}

	err := postBonusProducts(bonusProducts)
	if err != nil {
		return err
	}

	return nil
}

func checkProduct(product myProduct) (GQLQuery, error) {
	gql := new(Query)
	client := graphql.NewClient(os.Getenv("APPIE_URL"), nil).WithRequestModifier(func(request *http.Request) {
		request.Header.Add("client-name", os.Getenv("CLIENT_NAME"))
		request.Header.Add("client-version", os.Getenv("CLIENT_VERSION"))
	})

	vars := map[string]interface{}{
		"id":   product.ID,
		"date": time.Now().Format("2006-01-02"),
	}

	err := client.Query(context.Background(), &gql.Query, vars)
	if err != nil {
		log.Fatal(err)
	}

	return gql.Query, nil
}

func postBonusProducts(products []myProduct) error {
	for _, product := range products {
		err := postSlackMessage(product)
		if err != nil {
			return err
		}
	}

	return nil
}

func postSlackMessage(product myProduct) error {

	var message = new(slackMessage)
	message.Channel = os.Getenv("SLACK_CHANNEL")
	message.Blocks = append(message.Blocks, slackMessageBlock{
		Type: "section",
		Text: struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}{
			"mrkdwn",
			fmt.Sprintf("<%s%s|%s> \n *%s*",
				os.Getenv("APPIE_SITE"),
				product.BonusData.Product.WebPath,
				product.BonusData.Product.Title,
				product.BonusData.Product.Price.Discount.Description,
			),
		},
		Accessory: struct {
			Type     string `json:"type"`
			ImageUrl string `json:"image_url"`
			AltText  string `json:"alt_text"`
		}{
			"image",
			product.BonusData.Product.Images[0].Url,
			product.BonusData.Product.Title,
		},
	})
	// Marshall the message
	m, _ := json.Marshal(message) // to bytes

	// make the request
	url := fmt.Sprintf("%s/chat.postMessage", os.Getenv("SLACK_API"))
	request, err := http.NewRequest("POST", url, bytes.NewReader(m))
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("SLACK_TOKEN")))
	request.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	_, err = io.ReadAll(response.Body)

	// meh
	defer func() {
		err = response.Body.Close()
	}()

	return nil
}

func foobar() {
	fmt.Println("foobar")
}
