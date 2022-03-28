package firestore

import (
	"cloud.google.com/go/firestore"
)

// ApplyQuery takes collection reference and a custom query and applies the
// query.
//
// This function returns a firestore Query object.
func ApplyQuery(col *firestore.CollectionRef, query Query) firestore.Query {
	q := col.Query

	if query.Limit != 0 {
		q = q.Limit(query.Limit)
	}

	for _, order := range query.Orders {
		var dir firestore.Direction
		switch order.Direction {
		case "ASC":
			dir = firestore.Asc
		default:
			dir = firestore.Desc
		}
		q = q.OrderBy(order.Path, dir)
	}

	for _, filter := range query.Filters {
		q = q.Where(filter.Path, filter.Operator, filter.Value)
	}

	return q
}
