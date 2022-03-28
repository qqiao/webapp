package firestore

type Direction string

// Ordering directions
const (
	DirectionASC  Direction = "ASC"
	DirectionDESC Direction = "DESC"
)

// Query represents the abstraction of any datastore query.
type Query struct {
	Limit   int      `json:"limit"`
	Orders  []Order  `json:"orders"`
	Filters []Filter `json:"filters"`
}

// Order represents ordering criterion.
type Order struct {
	Path      string    `json:"path"`
	Direction Direction `json:"direction"`
}

// Filter represents filtering criterion.
type Filter struct {
	Path     string      `json:"path"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}
