package model

import (
	"errors"
	"fmt"
	"github.com/uptrace/bun"
	"strings"
	"time"
)

type Game struct {
	bun.BaseModel `bun:"table:games,alias:games"`

	ID        int        `bun:"id,pk,autoincrement"`
	Name      string     `bun:"name"`
	Url       string     `bun:"url"`
	CreatedAt time.Time  `bun:"created_at,nullzero"`
	DeletedAt *time.Time `bun:"deleted_at"`

	Prices     []Price      `bun:"rel:has-many,join:id=game_id"`
	Users      []User       `bun:"m2m:users_games,join:Game=User"`
	UsersGames []UsersGames `bun:"rel:has-many,join:id=game_id"`
}

func (g *Game) Validate() error {
	if g.Name == "" {
		return errors.New("name is empty")
	}

	if g.Url == "" {
		return errors.New("url is empty")
	}

	if len(g.Prices) == 0 {
		return errors.New("prices are empty")
	}

	return nil
}

func (g *Game) SetGameIDToPrices() {
	for i := range g.Prices {
		g.Prices[i].GameID = g.ID
	}
}

func (g *Game) GetMinPrice() float64 {
	min := -1.0

	for i := range g.Prices {
		if min < 0 {
			min = g.Prices[i].Value

			continue
		}

		if min > g.Prices[i].Value {
			min = g.Prices[i].Value
		}
	}

	return min
}

func (g *Game) GetPriceText() string {
	var text string

	for i := range g.Prices {
		text += fmt.Sprintf("- %s\n", g.Prices[i].GetPriceText())
	}

	text = strings.TrimSuffix(text, "\n")

	return text
}

func (g *Game) PriceChangedText() string {
	return fmt.Sprintf("Price changed!\nGame: %s\nNew price:\n%s\n\n%s",
		g.Name,
		g.GetPriceText(),
		g.Url,
	)
}
