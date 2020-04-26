// Stripped down OpenCypher parser for quering the graph.
package cypher

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type KV struct {
	Key   string
	Value []byte
}

type QueryPlan struct {
	ReadingClause []MatchQueryPlan
}

type MatchQueryPlan struct {
	Nodes []NodeQueryPlan
}

type NodeQueryPlan struct {
	Labels     []string
	Properties map[string][]byte
}

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: cypher query 'QUERY'")
	}

	got, err := ParseReader("", strings.NewReader(os.Args[1]))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=", got)
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Statement",
			pos:  position{line: 49, col: 1, offset: 730},
			expr: &actionExpr{
				pos: position{line: 49, col: 14, offset: 743},
				run: (*parser).callonStatement1,
				expr: &seqExpr{
					pos: position{line: 49, col: 14, offset: 743},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 49, col: 14, offset: 743},
							label: "query",
							expr: &ruleRefExpr{
								pos:  position{line: 49, col: 20, offset: 749},
								name: "Query",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 49, col: 26, offset: 755},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Query",
			pos:  position{line: 57, col: 1, offset: 874},
			expr: &actionExpr{
				pos: position{line: 57, col: 10, offset: 883},
				run: (*parser).callonQuery1,
				expr: &labeledExpr{
					pos:   position{line: 57, col: 10, offset: 883},
					label: "regularQuery",
					expr: &ruleRefExpr{
						pos:  position{line: 57, col: 23, offset: 896},
						name: "RegularQuery",
					},
				},
			},
		},
		{
			name: "RegularQuery",
			pos:  position{line: 61, col: 1, offset: 943},
			expr: &actionExpr{
				pos: position{line: 61, col: 18, offset: 960},
				run: (*parser).callonRegularQuery1,
				expr: &labeledExpr{
					pos:   position{line: 61, col: 18, offset: 960},
					label: "singleQuery",
					expr: &ruleRefExpr{
						pos:  position{line: 61, col: 30, offset: 972},
						name: "SingleQuery",
					},
				},
			},
		},
		{
			name: "SingleQuery",
			pos:  position{line: 65, col: 1, offset: 1017},
			expr: &actionExpr{
				pos: position{line: 65, col: 16, offset: 1032},
				run: (*parser).callonSingleQuery1,
				expr: &seqExpr{
					pos: position{line: 65, col: 16, offset: 1032},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 65, col: 16, offset: 1032},
							label: "clause",
							expr: &ruleRefExpr{
								pos:  position{line: 65, col: 23, offset: 1039},
								name: "ReadingClause",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 65, col: 37, offset: 1053},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "ReadingClause",
			pos:  position{line: 69, col: 1, offset: 1083},
			expr: &actionExpr{
				pos: position{line: 69, col: 18, offset: 1100},
				run: (*parser).callonReadingClause1,
				expr: &labeledExpr{
					pos:   position{line: 69, col: 18, offset: 1100},
					label: "match",
					expr: &ruleRefExpr{
						pos:  position{line: 69, col: 24, offset: 1106},
						name: "Match",
					},
				},
			},
		},
		{
			name: "Match",
			pos:  position{line: 73, col: 1, offset: 1139},
			expr: &actionExpr{
				pos: position{line: 73, col: 10, offset: 1148},
				run: (*parser).callonMatch1,
				expr: &seqExpr{
					pos: position{line: 73, col: 10, offset: 1148},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 73, col: 10, offset: 1148},
							name: "M",
						},
						&ruleRefExpr{
							pos:  position{line: 73, col: 12, offset: 1150},
							name: "A",
						},
						&ruleRefExpr{
							pos:  position{line: 73, col: 14, offset: 1152},
							name: "T",
						},
						&ruleRefExpr{
							pos:  position{line: 73, col: 16, offset: 1154},
							name: "C",
						},
						&ruleRefExpr{
							pos:  position{line: 73, col: 18, offset: 1156},
							name: "H",
						},
						&ruleRefExpr{
							pos:  position{line: 73, col: 20, offset: 1158},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 73, col: 22, offset: 1160},
							label: "pattern",
							expr: &ruleRefExpr{
								pos:  position{line: 73, col: 30, offset: 1168},
								name: "Pattern",
							},
						},
					},
				},
			},
		},
		{
			name: "Pattern",
			pos:  position{line: 78, col: 1, offset: 1278},
			expr: &ruleRefExpr{
				pos:  position{line: 78, col: 12, offset: 1289},
				name: "PatternPart",
			},
		},
		{
			name: "PatternPart",
			pos:  position{line: 80, col: 1, offset: 1303},
			expr: &ruleRefExpr{
				pos:  position{line: 80, col: 16, offset: 1318},
				name: "AnonymousPatternPart",
			},
		},
		{
			name: "AnonymousPatternPart",
			pos:  position{line: 82, col: 1, offset: 1340},
			expr: &ruleRefExpr{
				pos:  position{line: 82, col: 25, offset: 1364},
				name: "PatternElement",
			},
		},
		{
			name: "PatternElement",
			pos:  position{line: 84, col: 1, offset: 1380},
			expr: &ruleRefExpr{
				pos:  position{line: 84, col: 19, offset: 1398},
				name: "NodePattern",
			},
		},
		{
			name: "NodePattern",
			pos:  position{line: 86, col: 1, offset: 1411},
			expr: &actionExpr{
				pos: position{line: 86, col: 16, offset: 1426},
				run: (*parser).callonNodePattern1,
				expr: &seqExpr{
					pos: position{line: 86, col: 16, offset: 1426},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 86, col: 16, offset: 1426},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 86, col: 20, offset: 1430},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 86, col: 22, offset: 1432},
							label: "labels",
							expr: &zeroOrOneExpr{
								pos: position{line: 86, col: 29, offset: 1439},
								expr: &ruleRefExpr{
									pos:  position{line: 86, col: 29, offset: 1439},
									name: "NodeLabels",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 86, col: 41, offset: 1451},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 86, col: 43, offset: 1453},
							label: "props",
							expr: &zeroOrOneExpr{
								pos: position{line: 86, col: 49, offset: 1459},
								expr: &ruleRefExpr{
									pos:  position{line: 86, col: 50, offset: 1460},
									name: "Properties",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 86, col: 63, offset: 1473},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 86, col: 65, offset: 1475},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "NodeLabels",
			pos:  position{line: 100, col: 1, offset: 1685},
			expr: &actionExpr{
				pos: position{line: 100, col: 15, offset: 1699},
				run: (*parser).callonNodeLabels1,
				expr: &seqExpr{
					pos: position{line: 100, col: 15, offset: 1699},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 100, col: 15, offset: 1699},
							label: "label",
							expr: &ruleRefExpr{
								pos:  position{line: 100, col: 21, offset: 1705},
								name: "NodeLabel",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 100, col: 31, offset: 1715},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 100, col: 33, offset: 1717},
							label: "labels",
							expr: &zeroOrMoreExpr{
								pos: position{line: 100, col: 40, offset: 1724},
								expr: &ruleRefExpr{
									pos:  position{line: 100, col: 41, offset: 1725},
									name: "NodeLabel",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "NodeLabel",
			pos:  position{line: 111, col: 1, offset: 1933},
			expr: &actionExpr{
				pos: position{line: 111, col: 14, offset: 1946},
				run: (*parser).callonNodeLabel1,
				expr: &seqExpr{
					pos: position{line: 111, col: 14, offset: 1946},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 111, col: 14, offset: 1946},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 111, col: 18, offset: 1950},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 111, col: 20, offset: 1952},
							label: "label",
							expr: &ruleRefExpr{
								pos:  position{line: 111, col: 26, offset: 1958},
								name: "String",
							},
						},
					},
				},
			},
		},
		{
			name: "Properties",
			pos:  position{line: 115, col: 1, offset: 1992},
			expr: &ruleRefExpr{
				pos:  position{line: 115, col: 15, offset: 2006},
				name: "MapLiteral",
			},
		},
		{
			name: "ProperyKV",
			pos:  position{line: 116, col: 1, offset: 2017},
			expr: &actionExpr{
				pos: position{line: 116, col: 14, offset: 2030},
				run: (*parser).callonProperyKV1,
				expr: &seqExpr{
					pos: position{line: 116, col: 14, offset: 2030},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 116, col: 14, offset: 2030},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 116, col: 18, offset: 2034},
								name: "String",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 116, col: 25, offset: 2041},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 116, col: 27, offset: 2043},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 116, col: 31, offset: 2047},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 116, col: 33, offset: 2049},
							label: "value",
							expr: &choiceExpr{
								pos: position{line: 116, col: 40, offset: 2056},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 116, col: 40, offset: 2056},
										name: "StringLiteral",
									},
									&ruleRefExpr{
										pos:  position{line: 116, col: 54, offset: 2070},
										name: "Integer",
									},
									&ruleRefExpr{
										pos:  position{line: 116, col: 62, offset: 2078},
										name: "BoolLiteral",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "MapLiteral",
			pos:  position{line: 129, col: 1, offset: 2478},
			expr: &actionExpr{
				pos: position{line: 129, col: 15, offset: 2492},
				run: (*parser).callonMapLiteral1,
				expr: &seqExpr{
					pos: position{line: 129, col: 15, offset: 2492},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 129, col: 15, offset: 2492},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 129, col: 19, offset: 2496},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 129, col: 21, offset: 2498},
							label: "kv",
							expr: &zeroOrOneExpr{
								pos: position{line: 129, col: 24, offset: 2501},
								expr: &seqExpr{
									pos: position{line: 129, col: 25, offset: 2502},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 129, col: 25, offset: 2502},
											name: "ProperyKV",
										},
										&ruleRefExpr{
											pos:  position{line: 129, col: 35, offset: 2512},
											name: "_",
										},
										&zeroOrMoreExpr{
											pos: position{line: 129, col: 37, offset: 2514},
											expr: &seqExpr{
												pos: position{line: 129, col: 38, offset: 2515},
												exprs: []interface{}{
													&litMatcher{
														pos:        position{line: 129, col: 38, offset: 2515},
														val:        ",",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 129, col: 42, offset: 2519},
														name: "_",
													},
													&ruleRefExpr{
														pos:  position{line: 129, col: 44, offset: 2521},
														name: "ProperyKV",
													},
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 129, col: 59, offset: 2536},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 129, col: 61, offset: 2538},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "StringLiteral",
			pos:  position{line: 151, col: 1, offset: 2950},
			expr: &actionExpr{
				pos: position{line: 151, col: 18, offset: 2967},
				run: (*parser).callonStringLiteral1,
				expr: &labeledExpr{
					pos:   position{line: 151, col: 18, offset: 2967},
					label: "s",
					expr: &choiceExpr{
						pos: position{line: 151, col: 21, offset: 2970},
						alternatives: []interface{}{
							&seqExpr{
								pos: position{line: 151, col: 21, offset: 2970},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 151, col: 21, offset: 2970},
										val:        "'",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 151, col: 25, offset: 2974},
										name: "_",
									},
									&ruleRefExpr{
										pos:  position{line: 151, col: 27, offset: 2976},
										name: "String",
									},
									&ruleRefExpr{
										pos:  position{line: 151, col: 34, offset: 2983},
										name: "_",
									},
									&zeroOrMoreExpr{
										pos: position{line: 151, col: 36, offset: 2985},
										expr: &seqExpr{
											pos: position{line: 151, col: 37, offset: 2986},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 151, col: 37, offset: 2986},
													name: "String",
												},
												&ruleRefExpr{
													pos:  position{line: 151, col: 44, offset: 2993},
													name: "_",
												},
											},
										},
									},
									&litMatcher{
										pos:        position{line: 151, col: 48, offset: 2997},
										val:        "'",
										ignoreCase: false,
									},
								},
							},
							&seqExpr{
								pos: position{line: 151, col: 54, offset: 3003},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 151, col: 54, offset: 3003},
										val:        "\"",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 151, col: 58, offset: 3007},
										name: "_",
									},
									&ruleRefExpr{
										pos:  position{line: 151, col: 60, offset: 3009},
										name: "String",
									},
									&ruleRefExpr{
										pos:  position{line: 151, col: 67, offset: 3016},
										name: "_",
									},
									&zeroOrMoreExpr{
										pos: position{line: 151, col: 69, offset: 3018},
										expr: &seqExpr{
											pos: position{line: 151, col: 70, offset: 3019},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 151, col: 70, offset: 3019},
													name: "String",
												},
												&ruleRefExpr{
													pos:  position{line: 151, col: 77, offset: 3026},
													name: "_",
												},
											},
										},
									},
									&litMatcher{
										pos:        position{line: 151, col: 81, offset: 3030},
										val:        "\"",
										ignoreCase: false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "String",
			pos:  position{line: 155, col: 1, offset: 3089},
			expr: &actionExpr{
				pos: position{line: 155, col: 11, offset: 3099},
				run: (*parser).callonString1,
				expr: &oneOrMoreExpr{
					pos: position{line: 155, col: 11, offset: 3099},
					expr: &charClassMatcher{
						pos:        position{line: 155, col: 11, offset: 3099},
						val:        "[a-zA-Z0-9]",
						ranges:     []rune{'a', 'z', 'A', 'Z', '0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 159, col: 1, offset: 3148},
			expr: &actionExpr{
				pos: position{line: 159, col: 12, offset: 3159},
				run: (*parser).callonInteger1,
				expr: &oneOrMoreExpr{
					pos: position{line: 159, col: 12, offset: 3159},
					expr: &charClassMatcher{
						pos:        position{line: 159, col: 12, offset: 3159},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "BoolLiteral",
			pos:  position{line: 163, col: 1, offset: 3223},
			expr: &choiceExpr{
				pos: position{line: 163, col: 16, offset: 3238},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 163, col: 16, offset: 3238},
						run: (*parser).callonBoolLiteral2,
						expr: &seqExpr{
							pos: position{line: 163, col: 16, offset: 3238},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 163, col: 16, offset: 3238},
									name: "T",
								},
								&litMatcher{
									pos:        position{line: 163, col: 18, offset: 3240},
									val:        "rue",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 163, col: 47, offset: 3269},
						run: (*parser).callonBoolLiteral6,
						expr: &seqExpr{
							pos: position{line: 163, col: 47, offset: 3269},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 163, col: 47, offset: 3269},
									name: "F",
								},
								&litMatcher{
									pos:        position{line: 163, col: 49, offset: 3271},
									val:        "alse",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 165, col: 1, offset: 3300},
			expr: &zeroOrMoreExpr{
				pos: position{line: 165, col: 18, offset: 3319},
				expr: &charClassMatcher{
					pos:        position{line: 165, col: 18, offset: 3319},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "A",
			pos:  position{line: 167, col: 1, offset: 3331},
			expr: &choiceExpr{
				pos: position{line: 167, col: 7, offset: 3337},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 167, col: 7, offset: 3337},
						val:        "A",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 167, col: 13, offset: 3343},
						val:        "a",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "B",
			pos:  position{line: 168, col: 1, offset: 3348},
			expr: &choiceExpr{
				pos: position{line: 168, col: 7, offset: 3354},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 168, col: 7, offset: 3354},
						val:        "B",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 168, col: 13, offset: 3360},
						val:        "b",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "C",
			pos:  position{line: 169, col: 1, offset: 3365},
			expr: &choiceExpr{
				pos: position{line: 169, col: 7, offset: 3371},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 169, col: 7, offset: 3371},
						val:        "C",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 169, col: 13, offset: 3377},
						val:        "c",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "D",
			pos:  position{line: 170, col: 1, offset: 3382},
			expr: &choiceExpr{
				pos: position{line: 170, col: 7, offset: 3388},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 170, col: 7, offset: 3388},
						val:        "D",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 170, col: 13, offset: 3394},
						val:        "d",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "E",
			pos:  position{line: 171, col: 1, offset: 3399},
			expr: &choiceExpr{
				pos: position{line: 171, col: 7, offset: 3405},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 171, col: 7, offset: 3405},
						val:        "E",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 171, col: 13, offset: 3411},
						val:        "e",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "F",
			pos:  position{line: 172, col: 1, offset: 3416},
			expr: &choiceExpr{
				pos: position{line: 172, col: 7, offset: 3422},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 172, col: 7, offset: 3422},
						val:        "F",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 172, col: 13, offset: 3428},
						val:        "f",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "G",
			pos:  position{line: 173, col: 1, offset: 3433},
			expr: &choiceExpr{
				pos: position{line: 173, col: 7, offset: 3439},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 173, col: 7, offset: 3439},
						val:        "G",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 173, col: 13, offset: 3445},
						val:        "g",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "H",
			pos:  position{line: 174, col: 1, offset: 3450},
			expr: &choiceExpr{
				pos: position{line: 174, col: 7, offset: 3456},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 174, col: 7, offset: 3456},
						val:        "H",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 174, col: 13, offset: 3462},
						val:        "h",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "I",
			pos:  position{line: 175, col: 1, offset: 3467},
			expr: &choiceExpr{
				pos: position{line: 175, col: 7, offset: 3473},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 175, col: 7, offset: 3473},
						val:        "I",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 175, col: 13, offset: 3479},
						val:        "i",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "K",
			pos:  position{line: 176, col: 1, offset: 3484},
			expr: &choiceExpr{
				pos: position{line: 176, col: 7, offset: 3490},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 176, col: 7, offset: 3490},
						val:        "K",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 176, col: 13, offset: 3496},
						val:        "k",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "L",
			pos:  position{line: 177, col: 1, offset: 3501},
			expr: &choiceExpr{
				pos: position{line: 177, col: 7, offset: 3507},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 177, col: 7, offset: 3507},
						val:        "L",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 177, col: 13, offset: 3513},
						val:        "l",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "M",
			pos:  position{line: 178, col: 1, offset: 3518},
			expr: &choiceExpr{
				pos: position{line: 178, col: 7, offset: 3524},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 178, col: 7, offset: 3524},
						val:        "M",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 178, col: 13, offset: 3530},
						val:        "m",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "N",
			pos:  position{line: 179, col: 1, offset: 3535},
			expr: &choiceExpr{
				pos: position{line: 179, col: 7, offset: 3541},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 179, col: 7, offset: 3541},
						val:        "N",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 179, col: 13, offset: 3547},
						val:        "n",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "O",
			pos:  position{line: 180, col: 1, offset: 3552},
			expr: &choiceExpr{
				pos: position{line: 180, col: 7, offset: 3558},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 180, col: 7, offset: 3558},
						val:        "O",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 180, col: 13, offset: 3564},
						val:        "o",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "P",
			pos:  position{line: 181, col: 1, offset: 3569},
			expr: &choiceExpr{
				pos: position{line: 181, col: 7, offset: 3575},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 181, col: 7, offset: 3575},
						val:        "P",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 181, col: 13, offset: 3581},
						val:        "p",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Q",
			pos:  position{line: 182, col: 1, offset: 3586},
			expr: &choiceExpr{
				pos: position{line: 182, col: 7, offset: 3592},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 182, col: 7, offset: 3592},
						val:        "Q",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 182, col: 13, offset: 3598},
						val:        "q",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "R",
			pos:  position{line: 183, col: 1, offset: 3603},
			expr: &choiceExpr{
				pos: position{line: 183, col: 7, offset: 3609},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 183, col: 7, offset: 3609},
						val:        "R",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 183, col: 13, offset: 3615},
						val:        "r",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "S",
			pos:  position{line: 184, col: 1, offset: 3620},
			expr: &choiceExpr{
				pos: position{line: 184, col: 7, offset: 3626},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 184, col: 7, offset: 3626},
						val:        "S",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 184, col: 13, offset: 3632},
						val:        "s",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "T",
			pos:  position{line: 185, col: 1, offset: 3637},
			expr: &choiceExpr{
				pos: position{line: 185, col: 7, offset: 3643},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 185, col: 7, offset: 3643},
						val:        "T",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 185, col: 13, offset: 3649},
						val:        "t",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "U",
			pos:  position{line: 186, col: 1, offset: 3654},
			expr: &choiceExpr{
				pos: position{line: 186, col: 7, offset: 3660},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 186, col: 7, offset: 3660},
						val:        "U",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 186, col: 13, offset: 3666},
						val:        "u",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "V",
			pos:  position{line: 187, col: 1, offset: 3671},
			expr: &choiceExpr{
				pos: position{line: 187, col: 7, offset: 3677},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 187, col: 7, offset: 3677},
						val:        "V",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 187, col: 13, offset: 3683},
						val:        "v",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "W",
			pos:  position{line: 188, col: 1, offset: 3688},
			expr: &choiceExpr{
				pos: position{line: 188, col: 7, offset: 3694},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 188, col: 7, offset: 3694},
						val:        "W",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 188, col: 13, offset: 3700},
						val:        "w",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "X",
			pos:  position{line: 189, col: 1, offset: 3705},
			expr: &choiceExpr{
				pos: position{line: 189, col: 7, offset: 3711},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 189, col: 7, offset: 3711},
						val:        "X",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 189, col: 13, offset: 3717},
						val:        "x",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Y",
			pos:  position{line: 190, col: 1, offset: 3722},
			expr: &choiceExpr{
				pos: position{line: 190, col: 7, offset: 3728},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 190, col: 7, offset: 3728},
						val:        "Y",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 190, col: 13, offset: 3734},
						val:        "y",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 192, col: 1, offset: 3740},
			expr: &notExpr{
				pos: position{line: 192, col: 8, offset: 3747},
				expr: &anyMatcher{
					line: 192, col: 9, offset: 3748,
				},
			},
		},
	},
}

func (c *current) onStatement1(query interface{}) (interface{}, error) {
	q := QueryPlan{
		ReadingClause: []MatchQueryPlan{query.(MatchQueryPlan)},
	}

	return q, nil
}

func (p *parser) callonStatement1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStatement1(stack["query"])
}

func (c *current) onQuery1(regularQuery interface{}) (interface{}, error) {
	return regularQuery, nil
}

func (p *parser) callonQuery1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onQuery1(stack["regularQuery"])
}

func (c *current) onRegularQuery1(singleQuery interface{}) (interface{}, error) {
	return singleQuery, nil
}

func (p *parser) callonRegularQuery1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRegularQuery1(stack["singleQuery"])
}

func (c *current) onSingleQuery1(clause interface{}) (interface{}, error) {
	return clause, nil
}

func (p *parser) callonSingleQuery1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSingleQuery1(stack["clause"])
}

func (c *current) onReadingClause1(match interface{}) (interface{}, error) {
	return match, nil
}

func (p *parser) callonReadingClause1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReadingClause1(stack["match"])
}

func (c *current) onMatch1(pattern interface{}) (interface{}, error) {
	plan := MatchQueryPlan{Nodes: []NodeQueryPlan{pattern.(NodeQueryPlan)}}
	return plan, nil
}

func (p *parser) callonMatch1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMatch1(stack["pattern"])
}

func (c *current) onNodePattern1(labels, props interface{}) (interface{}, error) {
	plan := NodeQueryPlan{}

	if labels != nil {
		plan.Labels = labels.([]string)
	}

	if props != nil {
		plan.Properties = props.(map[string][]byte)
	}

	return plan, nil
}

func (p *parser) callonNodePattern1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNodePattern1(stack["labels"], stack["props"])
}

func (c *current) onNodeLabels1(label, labels interface{}) (interface{}, error) {
	l := make([]string, 1+len(labels.([]interface{})))
	l[0] = label.(string)

	for i, label := range labels.([]interface{}) {
		l[i+1] = label.(string)
	}

	return l, nil
}

func (p *parser) callonNodeLabels1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNodeLabels1(stack["label"], stack["labels"])
}

func (c *current) onNodeLabel1(label interface{}) (interface{}, error) {
	return label, nil
}

func (p *parser) callonNodeLabel1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNodeLabel1(stack["label"])
}

func (c *current) onProperyKV1(key, value interface{}) (interface{}, error) {
	switch value.(type) {
	case string:
		return KV{key.(string), []byte(value.(string))}, nil
	case int64:
		return KV{key.(string), []byte(fmt.Sprintf("%d", value.(int64)))}, nil
	case bool:
		return KV{key.(string), []byte(fmt.Sprintf("%t", value.(bool)))}, nil
	default:
		return nil, fmt.Errorf("Don't know what to do with %#v", value)
	}
}

func (p *parser) callonProperyKV1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onProperyKV1(stack["key"], stack["value"])
}

func (c *current) onMapLiteral1(kv interface{}) (interface{}, error) {
	log.Printf("%s", c.text)
	props := make(map[string][]byte)

	propsList := toIfaceSlice(kv)

	if len(propsList) == 0 {
		return props, nil
	}

	first := propsList[0].(KV)
	props[first.Key] = first.Value

	rest := toIfaceSlice(propsList[2])
	for _, pkv := range rest {
		kv := toIfaceSlice(pkv)[2].(KV)
		props[kv.Key] = kv.Value
	}

	return props, nil
}

func (p *parser) callonMapLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMapLiteral1(stack["kv"])
}

func (c *current) onStringLiteral1(s interface{}) (interface{}, error) {
	return string(c.text[1 : len(c.text)-1]), nil
}

func (p *parser) callonStringLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringLiteral1(stack["s"])
}

func (c *current) onString1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onString1()
}

func (c *current) onInteger1() (interface{}, error) {
	return strconv.ParseInt(string(c.text), 10, 32)
}

func (p *parser) callonInteger1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInteger1()
}

func (c *current) onBoolLiteral2() (interface{}, error) {
	return true, nil
}

func (p *parser) callonBoolLiteral2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBoolLiteral2()
}

func (c *current) onBoolLiteral6() (interface{}, error) {
	return false, nil
}

func (p *parser) callonBoolLiteral6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBoolLiteral6()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n == 1 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
