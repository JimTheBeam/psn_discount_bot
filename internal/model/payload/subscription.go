package payload

type (
	Subscribe struct {
		UserID int
		GameID int
		Price  float64
	}

	Subscriptions struct {
		UserID int
		Limit  int
		Offset int
	}
)
