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

import "github.com/qqiao/webapp/datastore"

// Direction is the sorting direction
//
// Deprecated: please use datastore.Direction instead.
type Direction = datastore.Direction

// Ordering directions
//
// Deprecated: please use their datastore.DirectionASC and
// datastore.DirectionDESC instead.
const (
	DirectionASC  = datastore.DirectionASC
	DirectionDESC = datastore.DirectionDESC
)

// Query represents the abstraction of any datastore query.
//
// Deprecated: please use datastore.Query instead.
type Query = datastore.Query

// Order represents ordering criterion.
//
// Deprecated: please use datastore.Order instead.
type Order = datastore.Order

// Filter represents filtering criterion.
//
// Deprecated: please use datastore.Filter instead.
type Filter = datastore.Filter
