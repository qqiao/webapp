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

// Direction is the sorting direction
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
