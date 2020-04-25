// Stripped down OpenCypher parser for quering the graph.
package main

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
	Lables     []string
	Properties map[string][]byte
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
			pos:  position{line: 42, col: 1, offset: 604},
			expr: &seqExpr{
				pos: position{line: 42, col: 14, offset: 617},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 42, col: 14, offset: 617},
						name: "Query",
					},
					&ruleRefExpr{
						pos:  position{line: 42, col: 20, offset: 623},
						name: "EOF",
					},
				},
			},
		},
		{
			name: "Query",
			pos:  position{line: 44, col: 1, offset: 628},
			expr: &ruleRefExpr{
				pos:  position{line: 44, col: 10, offset: 637},
				name: "RegularQuery",
			},
		},
		{
			name: "RegularQuery",
			pos:  position{line: 46, col: 1, offset: 651},
			expr: &ruleRefExpr{
				pos:  position{line: 46, col: 18, offset: 668},
				name: "SingleQuery",
			},
		},
		{
			name: "SingleQuery",
			pos:  position{line: 48, col: 1, offset: 681},
			expr: &seqExpr{
				pos: position{line: 48, col: 16, offset: 696},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 48, col: 16, offset: 696},
						name: "ReadingClause",
					},
					&ruleRefExpr{
						pos:  position{line: 48, col: 30, offset: 710},
						name: "_",
					},
				},
			},
		},
		{
			name: "ReadingClause",
			pos:  position{line: 50, col: 1, offset: 713},
			expr: &actionExpr{
				pos: position{line: 50, col: 18, offset: 730},
				run: (*parser).callonReadingClause1,
				expr: &ruleRefExpr{
					pos:  position{line: 50, col: 18, offset: 730},
					name: "Match",
				},
			},
		},
		{
			name: "Match",
			pos:  position{line: 55, col: 1, offset: 805},
			expr: &actionExpr{
				pos: position{line: 55, col: 10, offset: 814},
				run: (*parser).callonMatch1,
				expr: &seqExpr{
					pos: position{line: 55, col: 10, offset: 814},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 55, col: 10, offset: 814},
							name: "M",
						},
						&ruleRefExpr{
							pos:  position{line: 55, col: 12, offset: 816},
							name: "A",
						},
						&ruleRefExpr{
							pos:  position{line: 55, col: 14, offset: 818},
							name: "T",
						},
						&ruleRefExpr{
							pos:  position{line: 55, col: 16, offset: 820},
							name: "C",
						},
						&ruleRefExpr{
							pos:  position{line: 55, col: 18, offset: 822},
							name: "H",
						},
						&ruleRefExpr{
							pos:  position{line: 55, col: 20, offset: 824},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 55, col: 22, offset: 826},
							label: "pattern",
							expr: &ruleRefExpr{
								pos:  position{line: 55, col: 30, offset: 834},
								name: "Pattern",
							},
						},
					},
				},
			},
		},
		{
			name: "Pattern",
			pos:  position{line: 61, col: 1, offset: 987},
			expr: &ruleRefExpr{
				pos:  position{line: 61, col: 12, offset: 998},
				name: "PatternPart",
			},
		},
		{
			name: "PatternPart",
			pos:  position{line: 63, col: 1, offset: 1012},
			expr: &ruleRefExpr{
				pos:  position{line: 63, col: 16, offset: 1027},
				name: "AnonymousPatternPart",
			},
		},
		{
			name: "AnonymousPatternPart",
			pos:  position{line: 65, col: 1, offset: 1049},
			expr: &ruleRefExpr{
				pos:  position{line: 65, col: 25, offset: 1073},
				name: "PatternElement",
			},
		},
		{
			name: "PatternElement",
			pos:  position{line: 67, col: 1, offset: 1089},
			expr: &ruleRefExpr{
				pos:  position{line: 67, col: 19, offset: 1107},
				name: "NodePattern",
			},
		},
		{
			name: "NodePattern",
			pos:  position{line: 69, col: 1, offset: 1120},
			expr: &actionExpr{
				pos: position{line: 69, col: 16, offset: 1135},
				run: (*parser).callonNodePattern1,
				expr: &seqExpr{
					pos: position{line: 69, col: 16, offset: 1135},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 69, col: 16, offset: 1135},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 69, col: 20, offset: 1139},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 69, col: 22, offset: 1141},
							label: "labels",
							expr: &zeroOrOneExpr{
								pos: position{line: 69, col: 29, offset: 1148},
								expr: &ruleRefExpr{
									pos:  position{line: 69, col: 29, offset: 1148},
									name: "NodeLabels",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 69, col: 41, offset: 1160},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 69, col: 43, offset: 1162},
							label: "props",
							expr: &zeroOrOneExpr{
								pos: position{line: 69, col: 49, offset: 1168},
								expr: &ruleRefExpr{
									pos:  position{line: 69, col: 50, offset: 1169},
									name: "Properties",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 69, col: 63, offset: 1182},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 69, col: 65, offset: 1184},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "NodeLabels",
			pos:  position{line: 83, col: 1, offset: 1394},
			expr: &actionExpr{
				pos: position{line: 83, col: 15, offset: 1408},
				run: (*parser).callonNodeLabels1,
				expr: &seqExpr{
					pos: position{line: 83, col: 15, offset: 1408},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 83, col: 15, offset: 1408},
							label: "label",
							expr: &ruleRefExpr{
								pos:  position{line: 83, col: 21, offset: 1414},
								name: "NodeLabel",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 83, col: 31, offset: 1424},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 83, col: 33, offset: 1426},
							label: "labels",
							expr: &zeroOrMoreExpr{
								pos: position{line: 83, col: 40, offset: 1433},
								expr: &ruleRefExpr{
									pos:  position{line: 83, col: 41, offset: 1434},
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
			pos:  position{line: 94, col: 1, offset: 1642},
			expr: &actionExpr{
				pos: position{line: 94, col: 14, offset: 1655},
				run: (*parser).callonNodeLabel1,
				expr: &seqExpr{
					pos: position{line: 94, col: 14, offset: 1655},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 94, col: 14, offset: 1655},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 94, col: 18, offset: 1659},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 94, col: 20, offset: 1661},
							label: "label",
							expr: &ruleRefExpr{
								pos:  position{line: 94, col: 26, offset: 1667},
								name: "String",
							},
						},
					},
				},
			},
		},
		{
			name: "Properties",
			pos:  position{line: 98, col: 1, offset: 1701},
			expr: &ruleRefExpr{
				pos:  position{line: 98, col: 15, offset: 1715},
				name: "MapLiteral",
			},
		},
		{
			name: "ProperyKV",
			pos:  position{line: 99, col: 1, offset: 1726},
			expr: &actionExpr{
				pos: position{line: 99, col: 14, offset: 1739},
				run: (*parser).callonProperyKV1,
				expr: &seqExpr{
					pos: position{line: 99, col: 14, offset: 1739},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 99, col: 14, offset: 1739},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 99, col: 18, offset: 1743},
								name: "String",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 99, col: 25, offset: 1750},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 99, col: 27, offset: 1752},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 99, col: 31, offset: 1756},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 99, col: 33, offset: 1758},
							label: "value",
							expr: &choiceExpr{
								pos: position{line: 99, col: 40, offset: 1765},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 99, col: 40, offset: 1765},
										name: "StringLiteral",
									},
									&ruleRefExpr{
										pos:  position{line: 99, col: 54, offset: 1779},
										name: "Integer",
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
			pos:  position{line: 110, col: 1, offset: 2083},
			expr: &actionExpr{
				pos: position{line: 110, col: 15, offset: 2097},
				run: (*parser).callonMapLiteral1,
				expr: &seqExpr{
					pos: position{line: 110, col: 15, offset: 2097},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 110, col: 15, offset: 2097},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 110, col: 19, offset: 2101},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 110, col: 21, offset: 2103},
							label: "kv",
							expr: &ruleRefExpr{
								pos:  position{line: 110, col: 24, offset: 2106},
								name: "ProperyKV",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 110, col: 34, offset: 2116},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 110, col: 36, offset: 2118},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "StringLiteral",
			pos:  position{line: 117, col: 1, offset: 2230},
			expr: &actionExpr{
				pos: position{line: 117, col: 18, offset: 2247},
				run: (*parser).callonStringLiteral1,
				expr: &labeledExpr{
					pos:   position{line: 117, col: 18, offset: 2247},
					label: "s",
					expr: &choiceExpr{
						pos: position{line: 117, col: 21, offset: 2250},
						alternatives: []interface{}{
							&seqExpr{
								pos: position{line: 117, col: 21, offset: 2250},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 117, col: 21, offset: 2250},
										val:        "'",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 117, col: 25, offset: 2254},
										name: "_",
									},
									&ruleRefExpr{
										pos:  position{line: 117, col: 27, offset: 2256},
										name: "String",
									},
									&ruleRefExpr{
										pos:  position{line: 117, col: 34, offset: 2263},
										name: "_",
									},
									&litMatcher{
										pos:        position{line: 117, col: 36, offset: 2265},
										val:        "'",
										ignoreCase: false,
									},
								},
							},
							&seqExpr{
								pos: position{line: 117, col: 42, offset: 2271},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 117, col: 42, offset: 2271},
										val:        "\"",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 117, col: 46, offset: 2275},
										name: "_",
									},
									&ruleRefExpr{
										pos:  position{line: 117, col: 48, offset: 2277},
										name: "String",
									},
									&ruleRefExpr{
										pos:  position{line: 117, col: 55, offset: 2284},
										name: "_",
									},
									&litMatcher{
										pos:        position{line: 117, col: 57, offset: 2286},
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
			pos:  position{line: 122, col: 1, offset: 2356},
			expr: &actionExpr{
				pos: position{line: 122, col: 11, offset: 2366},
				run: (*parser).callonString1,
				expr: &oneOrMoreExpr{
					pos: position{line: 122, col: 11, offset: 2366},
					expr: &charClassMatcher{
						pos:        position{line: 122, col: 11, offset: 2366},
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
			pos:  position{line: 126, col: 1, offset: 2415},
			expr: &actionExpr{
				pos: position{line: 126, col: 12, offset: 2426},
				run: (*parser).callonInteger1,
				expr: &oneOrMoreExpr{
					pos: position{line: 126, col: 12, offset: 2426},
					expr: &charClassMatcher{
						pos:        position{line: 126, col: 12, offset: 2426},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 130, col: 1, offset: 2490},
			expr: &zeroOrMoreExpr{
				pos: position{line: 130, col: 18, offset: 2509},
				expr: &charClassMatcher{
					pos:        position{line: 130, col: 18, offset: 2509},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "A",
			pos:  position{line: 132, col: 1, offset: 2521},
			expr: &choiceExpr{
				pos: position{line: 132, col: 7, offset: 2527},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 132, col: 7, offset: 2527},
						val:        "A",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 132, col: 13, offset: 2533},
						val:        "a",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "B",
			pos:  position{line: 133, col: 1, offset: 2538},
			expr: &choiceExpr{
				pos: position{line: 133, col: 7, offset: 2544},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 133, col: 7, offset: 2544},
						val:        "B",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 133, col: 13, offset: 2550},
						val:        "b",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "C",
			pos:  position{line: 134, col: 1, offset: 2555},
			expr: &choiceExpr{
				pos: position{line: 134, col: 7, offset: 2561},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 134, col: 7, offset: 2561},
						val:        "C",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 134, col: 13, offset: 2567},
						val:        "c",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "D",
			pos:  position{line: 135, col: 1, offset: 2572},
			expr: &choiceExpr{
				pos: position{line: 135, col: 7, offset: 2578},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 135, col: 7, offset: 2578},
						val:        "D",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 135, col: 13, offset: 2584},
						val:        "d",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "E",
			pos:  position{line: 136, col: 1, offset: 2589},
			expr: &choiceExpr{
				pos: position{line: 136, col: 7, offset: 2595},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 136, col: 7, offset: 2595},
						val:        "E",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 136, col: 13, offset: 2601},
						val:        "e",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "F",
			pos:  position{line: 137, col: 1, offset: 2606},
			expr: &choiceExpr{
				pos: position{line: 137, col: 7, offset: 2612},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 137, col: 7, offset: 2612},
						val:        "F",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 137, col: 13, offset: 2618},
						val:        "f",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "G",
			pos:  position{line: 138, col: 1, offset: 2623},
			expr: &choiceExpr{
				pos: position{line: 138, col: 7, offset: 2629},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 138, col: 7, offset: 2629},
						val:        "G",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 138, col: 13, offset: 2635},
						val:        "g",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "H",
			pos:  position{line: 139, col: 1, offset: 2640},
			expr: &choiceExpr{
				pos: position{line: 139, col: 7, offset: 2646},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 139, col: 7, offset: 2646},
						val:        "H",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 139, col: 13, offset: 2652},
						val:        "h",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "I",
			pos:  position{line: 140, col: 1, offset: 2657},
			expr: &choiceExpr{
				pos: position{line: 140, col: 7, offset: 2663},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 140, col: 7, offset: 2663},
						val:        "I",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 140, col: 13, offset: 2669},
						val:        "i",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "K",
			pos:  position{line: 141, col: 1, offset: 2674},
			expr: &choiceExpr{
				pos: position{line: 141, col: 7, offset: 2680},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 141, col: 7, offset: 2680},
						val:        "K",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 141, col: 13, offset: 2686},
						val:        "k",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "L",
			pos:  position{line: 142, col: 1, offset: 2691},
			expr: &choiceExpr{
				pos: position{line: 142, col: 7, offset: 2697},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 142, col: 7, offset: 2697},
						val:        "L",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 142, col: 13, offset: 2703},
						val:        "l",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "M",
			pos:  position{line: 143, col: 1, offset: 2708},
			expr: &choiceExpr{
				pos: position{line: 143, col: 7, offset: 2714},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 143, col: 7, offset: 2714},
						val:        "M",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 143, col: 13, offset: 2720},
						val:        "m",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "N",
			pos:  position{line: 144, col: 1, offset: 2725},
			expr: &choiceExpr{
				pos: position{line: 144, col: 7, offset: 2731},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 144, col: 7, offset: 2731},
						val:        "N",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 144, col: 13, offset: 2737},
						val:        "n",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "O",
			pos:  position{line: 145, col: 1, offset: 2742},
			expr: &choiceExpr{
				pos: position{line: 145, col: 7, offset: 2748},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 145, col: 7, offset: 2748},
						val:        "O",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 145, col: 13, offset: 2754},
						val:        "o",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "P",
			pos:  position{line: 146, col: 1, offset: 2759},
			expr: &choiceExpr{
				pos: position{line: 146, col: 7, offset: 2765},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 146, col: 7, offset: 2765},
						val:        "P",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 146, col: 13, offset: 2771},
						val:        "p",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Q",
			pos:  position{line: 147, col: 1, offset: 2776},
			expr: &choiceExpr{
				pos: position{line: 147, col: 7, offset: 2782},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 147, col: 7, offset: 2782},
						val:        "Q",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 147, col: 13, offset: 2788},
						val:        "q",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "R",
			pos:  position{line: 148, col: 1, offset: 2793},
			expr: &choiceExpr{
				pos: position{line: 148, col: 7, offset: 2799},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 148, col: 7, offset: 2799},
						val:        "R",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 148, col: 13, offset: 2805},
						val:        "r",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "S",
			pos:  position{line: 149, col: 1, offset: 2810},
			expr: &choiceExpr{
				pos: position{line: 149, col: 7, offset: 2816},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 149, col: 7, offset: 2816},
						val:        "S",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 149, col: 13, offset: 2822},
						val:        "s",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "T",
			pos:  position{line: 150, col: 1, offset: 2827},
			expr: &choiceExpr{
				pos: position{line: 150, col: 7, offset: 2833},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 150, col: 7, offset: 2833},
						val:        "T",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 150, col: 13, offset: 2839},
						val:        "t",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "U",
			pos:  position{line: 151, col: 1, offset: 2844},
			expr: &choiceExpr{
				pos: position{line: 151, col: 7, offset: 2850},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 151, col: 7, offset: 2850},
						val:        "U",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 151, col: 13, offset: 2856},
						val:        "u",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "V",
			pos:  position{line: 152, col: 1, offset: 2861},
			expr: &choiceExpr{
				pos: position{line: 152, col: 7, offset: 2867},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 152, col: 7, offset: 2867},
						val:        "V",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 152, col: 13, offset: 2873},
						val:        "v",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "W",
			pos:  position{line: 153, col: 1, offset: 2878},
			expr: &choiceExpr{
				pos: position{line: 153, col: 7, offset: 2884},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 153, col: 7, offset: 2884},
						val:        "W",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 153, col: 13, offset: 2890},
						val:        "w",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "X",
			pos:  position{line: 154, col: 1, offset: 2895},
			expr: &choiceExpr{
				pos: position{line: 154, col: 7, offset: 2901},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 154, col: 7, offset: 2901},
						val:        "X",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 154, col: 13, offset: 2907},
						val:        "x",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Y",
			pos:  position{line: 155, col: 1, offset: 2912},
			expr: &choiceExpr{
				pos: position{line: 155, col: 7, offset: 2918},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 155, col: 7, offset: 2918},
						val:        "Y",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 155, col: 13, offset: 2924},
						val:        "y",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 157, col: 1, offset: 2930},
			expr: &notExpr{
				pos: position{line: 157, col: 8, offset: 2937},
				expr: &anyMatcher{
					line: 157, col: 9, offset: 2938,
				},
			},
		},
	},
}

func (c *current) onReadingClause1() (interface{}, error) {
	log.Printf("ReadingClause: %s", c.text)
	return nil, nil
}

func (p *parser) callonReadingClause1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReadingClause1()
}

func (c *current) onMatch1(pattern interface{}) (interface{}, error) {
	plan := MatchQueryPlan{Nodes: []NodeQueryPlan{pattern.(NodeQueryPlan)}}
	log.Printf("match pattern: %#v", plan)
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
		plan.Lables = labels.([]string)
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
	props := make(map[string][]byte)
	k := kv.(KV)
	props[k.Key] = k.Value
	return props, nil
}

func (p *parser) callonMapLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMapLiteral1(stack["kv"])
}

func (c *current) onStringLiteral1(s interface{}) (interface{}, error) {
	ss := s.([]interface{})
	return ss[2].(string), nil
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
