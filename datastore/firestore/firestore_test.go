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

package firestore_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/qqiao/webapp/v2/datastore"

	fs "github.com/qqiao/webapp/v2/datastore/firestore"
)

var client *firestore.Client

const collectionNameOrTest = "test_or_collection"

func setUp() {
	c, err := firestore.NewClient(context.Background(), "test-project")
	if err != nil {
		log.Fatalf("Unable to initialize firebase client. Error: %v", err)
	}
	client = c

	setUpOrTest()
}

func setUpOrTest() {
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		if _, err := client.Collection(collectionNameOrTest).Doc(fmt.Sprintf(
			"Or-%d",
			i)).Set(
			ctx, map[string]string{
				"name": fmt.Sprintf("Or-%d", i),
			}); err != nil {
			log.Fatalf("Unable to initialize data for Or test: %v", err)
		}
	}
}

func TestMain(m *testing.M) {
	setUp()
	os.Exit(m.Run())
}

func TestOr(t *testing.T) {
	var expected sync.Map
	expected.Store("Or-1", false)
	expected.Store("Or-4", false)
	expected.Store("Or-7", false)

	col := client.Collection(collectionNameOrTest)
	queries := make([]datastore.Query, 0)
	expected.Range(func(key, _ any) bool {
		queries = append(queries, datastore.Query{
			Filters: []datastore.Filter{{
				Path:     "name",
				Operator: "==",
				Value:    key.(string),
			}},
		})
		return true
	})

	if err := client.RunTransaction(context.Background(), func(ctx context.Context,
		transaction *firestore.Transaction) error {
		o, e := fs.Or[map[string]string](ctx, 3, 0, transaction, col,
			queries...)
		for {
			select {
			case found, ok := <-o:
				if !ok {
					return nil
				}
				expected.Store(found["name"], true)
			case er := <-e:
				if er != nil {
					return er
				}
			}
		}
	}); err != nil {
		t.Errorf("Error testing Or: %v", err)
	}

	expected.Range(func(key, value any) bool {
		v := value.(bool)
		if !v {
			t.Errorf("%s should have been found by query", key.(string))
		}
		return v
	})
}
