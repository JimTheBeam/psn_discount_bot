package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"time"
)

const (
	GeneralType = "usual"
	PSPlusType  = "ps_plus"
	EAType      = "ea"
)

type Price struct {
	bun.BaseModel `bun:"table:prices,alias:prices"`

	ID        int       `bun:"id,pk,autoincrement"`
	GameID    int       `bun:"game_id"`
	Value     float64   `bun:"value"`
	IsFree    bool      `bun:"is_free"`
	Type      string    `bun:"type"`
	Currency  string    `bun:"currency"`
	CreatedAt time.Time `bun:"created_at,nullzero"`
}

func (p *Price) GetPriceText() string {
	var text string

	if p.IsFree {
		text = "FREE"
	} else {
		text = fmt.Sprintf("%.2f %s", p.Value, p.Currency)
	}

	switch p.Type {
	case GeneralType:
	case PSPlusType:
		text += " with PlayStation Plus"
	case EAType:
		text += " with EA subscription"
	}

	return text
}
