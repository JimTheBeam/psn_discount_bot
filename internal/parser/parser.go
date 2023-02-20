package parser

import (
	"fmt"
	"net/http"
	"psn_discount_bot/internal/model"
)

type HTTPClient interface {
	DoRequest(method, url string, body []byte) ([]byte, error)
}

type Parser struct {
	httpClient HTTPClient
}

func NewParser(client HTTPClient) *Parser {
	return &Parser{
		httpClient: client,
	}
}

func (p *Parser) ParseGame(url string) (model.Game, error) {
	body, err := p.httpClient.DoRequest(http.MethodGet, url, nil)
	if err != nil || len(body) == 0 {
		return model.Game{}, fmt.Errorf("do request: %w", err)
	}

	game, err := parseGamePage(body)
	if err != nil {
		return model.Game{}, fmt.Errorf("parse page: %w", err)
	}

	game.Url = url

	if err := game.Validate(); err != nil {
		return model.Game{}, fmt.Errorf("validate game after parsing: %w", err)
	}

	return game, nil
}
