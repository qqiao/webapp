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

package datastore

// Direction represents the sorting direction.
type Direction string

// Possible ordering directions.
const (
	DirectionASC  Direction = "ASC"
	DirectionDESC Direction = "DESC"
)

// Query represents the abstraction of any datastore query.
//
// For a query with multiple Filters, they are treated as a set of criterion
// joined with AND condition behind the scenes.
//
// For queries needing to use the OR condition, it is more efficient to split
// the query into multiple separate ones, run them separately in concurrently
// and combine the results afterwards.
type Query struct {
	Limit   int      `json:"limit"`
	Orders  []Order  `json:"orders"`
	Filters []Filter `json:"filters"`
}

// Order represents an ordering criteria.
type Order struct {
	Path      string    `json:"path"`
	Direction Direction `json:"direction"`
}

// Filter represents a filtering criteria.
type Filter struct {
	Path     string      `json:"path"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}
