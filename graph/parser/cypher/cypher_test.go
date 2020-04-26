package cypher

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Query    string
	Expected QueryPlan
}

func TestMatchQueries(t *testing.T) {
	tests := []TestCase{
		TestCase{
			Query: `MATCH (:Person)`,
			Expected: QueryPlan{
				ReadingClause: []MatchQueryPlan{
					MatchQueryPlan{
						[]NodeQueryPlan{NodeQueryPlan{Labels: []string{"Person"}}},
					},
				},
			},
		},
		TestCase{
			Query: `MATCH (:Person:Animal)`,
			Expected: QueryPlan{
				ReadingClause: []MatchQueryPlan{
					MatchQueryPlan{
						[]NodeQueryPlan{NodeQueryPlan{Labels: []string{"Person", "Animal"}}},
					},
				},
			},
		},
		TestCase{
			Query: `MATCH (:Person:Animal {})`,
			Expected: QueryPlan{
				ReadingClause: []MatchQueryPlan{
					MatchQueryPlan{
						[]NodeQueryPlan{
							NodeQueryPlan{
								Labels:     []string{"Person", "Animal"},
								Properties: map[string][]byte{},
							},
						},
					},
				},
			},
		},
		TestCase{
			Query: `MATCH (:Person:Animal {name: "Foo", surname: 'Bar', age: 21})`,
			Expected: QueryPlan{
				ReadingClause: []MatchQueryPlan{
					MatchQueryPlan{
						[]NodeQueryPlan{
							NodeQueryPlan{
								Labels: []string{"Person", "Animal"},
								Properties: map[string][]byte{
									"name":    []byte("Foo"),
									"surname": []byte("Bar"),
									"age":     []byte("21"),
								},
							},
						},
					},
				},
			},
		},
		TestCase{
			Query: `MATCH (:Person:Animal {name: "Foo", surname: 'The Bar'})`,
			Expected: QueryPlan{
				ReadingClause: []MatchQueryPlan{
					MatchQueryPlan{
						[]NodeQueryPlan{
							NodeQueryPlan{
								Labels: []string{"Person", "Animal"},
								Properties: map[string][]byte{
									"name":    []byte("Foo"),
									"surname": []byte("The Bar"),
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		got, err := Parse("", []byte(test.Query))
		if err != nil {
			log.Printf("Error: %s", err)
		}
		assert.Nil(t, err)
		actual := got.(QueryPlan)
		assert.Equal(t, test.Expected, actual)

	}
}
