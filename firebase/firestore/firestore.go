// Copyright 2022 Qian Qiao
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
