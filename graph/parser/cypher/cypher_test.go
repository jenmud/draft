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
			Query: `MATCH (n:Person)`,
			Expected: QueryPlan{
				ReadingClause: []MatchQueryPlan{
					MatchQueryPlan{
						[]NodeQueryPlan{
							NodeQueryPlan{
								Variable: "n",
								Labels:   []string{"Person"},
							},
						},
					},
				},
			},
		},
		TestCase{
			Query: `MATCH (node:Person)`,
			Expected: QueryPlan{
				ReadingClause: []MatchQueryPlan{
					MatchQueryPlan{
						[]NodeQueryPlan{
							NodeQueryPlan{
								Variable: "node",
								Labels:   []string{"Person"},
							},
						},
					},
				},
			},
		},
		TestCase{
			Query: `MATCH (n:Person:Animal)`,
			Expected: QueryPlan{
				ReadingClause: []MatchQueryPlan{
					MatchQueryPlan{
						[]NodeQueryPlan{
							NodeQueryPlan{
								Variable: "n",
								Labels:   []string{"Person", "Animal"},
							},
						},
					},
				},
			},
		},
		TestCase{
			Query: `MATCH (n:Person:Animal {})`,
			Expected: QueryPlan{
				ReadingClause: []MatchQueryPlan{
					MatchQueryPlan{
						[]NodeQueryPlan{
							NodeQueryPlan{
								Variable:   "n",
								Labels:     []string{"Person", "Animal"},
								Properties: map[string][]byte{},
							},
						},
					},
				},
			},
		},
		TestCase{
			Query: `MATCH (n:Person:Animal {name: "Foo", surname: 'Foo-Bar', age: 21, sex: "Rather not say", male: True, female: false})`,
			Expected: QueryPlan{
				ReadingClause: []MatchQueryPlan{
					MatchQueryPlan{
						[]NodeQueryPlan{
							NodeQueryPlan{
								Variable: "n",
								Labels:   []string{"Person", "Animal"},
								Properties: map[string][]byte{
									"name":    []byte("Foo"),
									"surname": []byte("Foo-Bar"),
									"age":     []byte("21"),
									"sex":     []byte("Rather not say"),
									"male":    []byte("true"),
									"female":  []byte("false"),
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
