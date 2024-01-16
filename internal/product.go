package internal

type Product struct {
	ID          string   `json:"id"`
	Status      string   `json:"status"` // 'available' or 'unavailable'
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Categories  []string `json:"categories"`
	Description string   `json:"description"`
}
