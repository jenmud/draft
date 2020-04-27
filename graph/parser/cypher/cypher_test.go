package cypher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Query       string
	Expected    QueryPlan
	ShouldError bool
}

func TestMatchQueries(t *testing.T) {
	tests := []TestCase{
		TestCase{
			Query: `MATCH (n:Person) RETURN n`,
			Expected: QueryPlan{
				ReadingClause: []ReadingClause{
					ReadingClause{
						Returns: []string{"n"},
						Match: MatchQuery{
							Nodes: []NodeQuery{
								NodeQuery{
									Variable: "n",
									Labels:   []string{"Person"},
								},
							},
						},
					},
				},
			},
		},
		TestCase{
			Query: `MATCH (n:Person:Animal) RETURN n`,
			Expected: QueryPlan{
				ReadingClause: []ReadingClause{
					ReadingClause{
						Returns: []string{"n"},
						Match: MatchQuery{
							Nodes: []NodeQuery{
								NodeQuery{
									Variable: "n",
									Labels:   []string{"Person", "Animal"},
								},
							},
						},
					},
				},
			},
		},
		TestCase{
			Query: `MATCH (n:Person:Animal {name: "Foo", surname: 'Bar', age: 21, active: true, address: "My address is private"}) RETURN n`,
			Expected: QueryPlan{
				ReadingClause: []ReadingClause{
					ReadingClause{
						Returns: []string{"n"},
						Match: MatchQuery{
							Nodes: []NodeQuery{
								NodeQuery{
									Variable: "n",
									Labels:   []string{"Person", "Animal"},
									Properties: map[string][]byte{
										"name":    []byte("Foo"),
										"surname": []byte("Bar"),
										"age":     []byte("21"),
										"active":  []byte("true"),
										"address": []byte("My address is private"),
									},
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
		if !test.ShouldError {
			assert.Nil(t, err)
			actual := got.(QueryPlan)
			assert.Equal(t, test.Expected, actual)
		} else {
			assert.NotNil(t, err, "Expected query %s to fail", test.Query)
		}
	}
}
