package parser

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"math"
	"psn_discount_bot/internal/model"
	"strconv"
	"strings"
)

const (
	freePrice      = "Free"
	includedPrice  = "Included"
	gameTrialPrice = "Game Trial"

	//todo;
	defaultCurrency = "TL"
)

const (
	eaPlay      = "EA Play"
	playstation = "PlayStation"
)

func parseGamePage(body []byte) (model.Game, error) {
	content, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return model.Game{}, fmt.Errorf("new document from reader: %w", err)
	}

	game := parsePsnHTML(content)

	return game, nil
}

// parseHabrHTML - parse articles from habr HTML.
func parsePsnHTML(doc *goquery.Document) model.Game {
	gameCard := doc.FindMatcher(goquery.Single(".psw-c-bg-card-1"))

	gameName := gameCard.Find(".psw-m-b-5").Text()

	prices := parseGamePrices(gameCard)

	game := model.Game{
		Name:   gameName,
		Prices: prices,
	}

	return game
}

func parseGamePrices(gameCard *goquery.Selection) []model.Price {
	prices := make([]model.Price, 0, 1)

	gameCard.Find(".psw-l-anchor.psw-l-stack-left.psw-fill-x").Each(func(i int, priceTitle *goquery.Selection) {
		priceStr := priceTitle.Find(".psw-t-title-m").Text()

		if priceStr == "" {
			return
		}

		psPlusString := priceTitle.Find(".psw-c-t-ps-plus").Text()
		eaString := priceTitle.Find(".psw-m-r-3").Text()

		isPSPlus := strings.Contains(psPlusString, playstation)
		isEA := strings.Contains(eaString, eaPlay)

		var subscriptionType string
		switch {
		case isPSPlus:
			subscriptionType = model.PSPlusType
		case isEA:
			subscriptionType = model.EAType
		default:
			subscriptionType = model.GeneralType
		}

		price := model.Price{
			Type: subscriptionType,
		}

		switch priceStr {
		case freePrice, includedPrice:
			price.IsFree = true
		case gameTrialPrice:
			return
		}

		if !price.IsFree {
			gamePrice, gameCurrency := parsePriceCurrency(priceStr)

			price.Value = gamePrice
			price.Currency = gameCurrency
		}

		prices = append(prices, price)
	})

	return prices
}

func parsePriceCurrency(priceStr string) (float64, string) {
	parts := strings.Split(priceStr, " ")

	if len(parts) != 2 {
		return 0, ""
	}

	price := strings.Replace(parts[0], ".", "", 1)

	price = strings.Replace(price, ",", ".", 1)

	priceFloat, err := strconv.ParseFloat(price, 32)
	if err != nil {
		return 0, ""
	}

	priceFloat = math.Round(priceFloat*100) / 100

	return priceFloat, parts[1]
}
