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
	"context"
	"reflect"

	"cloud.google.com/go/firestore"
	"github.com/qqiao/pipeline/v2"
	"github.com/qqiao/webapp/v2/datastore"
	"google.golang.org/api/iterator"
)

// ApplyQuery takes collection reference and a custom query and applies the
// query to the collection reference.
func ApplyQuery(col *firestore.CollectionRef,
	query datastore.Query) firestore.Query {
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

// Or takes a set of datastore queries, and run them in the same transaction,
// with OR condition connecting the queries.
//
// Due to the fact that firestore does not support OR query, we have to
// simulate an OR query with multiple simultaneous queries in parallel.
//
// There a few limitations users should be aware of with the current
// implementation:
//     1. Ordering will not function. Since queries are run in parallel,
//     to prevent any potential data race conditions, the results of the
//     queries are simply streamed, with no consideration of ordering.
//     2. Limit on the overall result set will not function. Given that queries
//     are run in parallel and ordering will not function, there is currently
//     no way for the function ot support limiting the results.
//
// As a workaround, users should take the total results, apply any sorting,
// further filtering and limiting of the results in their own code.
func Or[O any](ctx context.Context, parallelQueries int, bufferSize int,
	t *firestore.Transaction, col *firestore.CollectionRef,
	queries ...datastore.Query) (<-chan O, <-chan error) {
	out := make(chan O)
	err := make(chan error)

	go func() {
		defer close(out)
		defer close(err)

		// make sure that we feed the workers
		in := make(chan firestore.Query, len(queries))
		go func() {
			defer close(in)
			for _, query := range queries {
				in <- ApplyQuery(col, query)
			}
		}()

		sw := func(ctx context.Context, producer pipeline.Producer[firestore.
			Query]) (<-chan O, <-chan error) {
			out := make(chan O)
			err := make(chan error)

			go func() {
				defer close(out)
				defer close(err)

				for query := range producer {
					func() {
						iter := t.Documents(query)
						defer iter.Stop()

						for {
							ref, e := iter.Next()

							if e == iterator.Done {
								return
							}

							if e != nil {
								err <- e
								return
							}

							var object O
							var toLoad any

							// We need to deal with the case where O is a
							// pointer type. In such case,
							// we have to use reflection to instantiate a real
							// instance of the object, and use its address
							rv := reflect.ValueOf(object)
							if rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
								t := reflect.TypeOf(object).Elem()
								toLoad = reflect.Indirect(reflect.New(t)).
									Addr().Interface()
								object = toLoad.(O)
							} else {
								toLoad = &object
							}

							if e = ref.DataTo(toLoad); e != nil {
								err <- e
								return
							}
							out <- object
						}
					}()
				}
			}()
			return out, err
		}

		stage, e := pipeline.NewStageStreamWorker(parallelQueries,
			bufferSize, in, sw)
		if e != nil {
			err <- e
			return
		}

		so := stage.Produces()
		ec := stage.Start(ctx)

		for {
			select {
			case o, ok := <-so:
				if !ok {
					return
				}
				out <- o
			case e, ok := <-ec:
				if ok && e != nil {
					err <- e
					return
				}
			}
		}
	}()

	return out, err
}
