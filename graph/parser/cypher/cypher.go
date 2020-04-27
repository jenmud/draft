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
	ReadingClause []ReadingClause
}

type ReadingClause struct {
	Match   MatchQuery
	Returns []string
}

type MatchQuery struct {
	Nodes []NodeQuery
}

type NodeQuery struct {
	Variable   string
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
			pos:  position{line: 55, col: 1, offset: 864},
			expr: &actionExpr{
				pos: position{line: 55, col: 14, offset: 877},
				run: (*parser).callonStatement1,
				expr: &seqExpr{
					pos: position{line: 55, col: 14, offset: 877},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 55, col: 14, offset: 877},
							label: "query",
							expr: &ruleRefExpr{
								pos:  position{line: 55, col: 20, offset: 883},
								name: "Query",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 55, col: 26, offset: 889},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Query",
			pos:  position{line: 64, col: 1, offset: 1044},
			expr: &actionExpr{
				pos: position{line: 64, col: 10, offset: 1053},
				run: (*parser).callonQuery1,
				expr: &labeledExpr{
					pos:   position{line: 64, col: 10, offset: 1053},
					label: "regularQuery",
					expr: &ruleRefExpr{
						pos:  position{line: 64, col: 23, offset: 1066},
						name: "RegularQuery",
					},
				},
			},
		},
		{
			name: "RegularQuery",
			pos:  position{line: 68, col: 1, offset: 1117},
			expr: &actionExpr{
				pos: position{line: 68, col: 18, offset: 1134},
				run: (*parser).callonRegularQuery1,
				expr: &labeledExpr{
					pos:   position{line: 68, col: 18, offset: 1134},
					label: "singleQuery",
					expr: &ruleRefExpr{
						pos:  position{line: 68, col: 30, offset: 1146},
						name: "SingleQuery",
					},
				},
			},
		},
		{
			name: "SingleQuery",
			pos:  position{line: 72, col: 1, offset: 1195},
			expr: &actionExpr{
				pos: position{line: 72, col: 16, offset: 1210},
				run: (*parser).callonSingleQuery1,
				expr: &seqExpr{
					pos: position{line: 72, col: 16, offset: 1210},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 72, col: 16, offset: 1210},
							label: "clause",
							expr: &ruleRefExpr{
								pos:  position{line: 72, col: 23, offset: 1217},
								name: "ReadingClause",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 72, col: 37, offset: 1231},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 72, col: 39, offset: 1233},
							label: "ret",
							expr: &ruleRefExpr{
								pos:  position{line: 72, col: 43, offset: 1237},
								name: "Return",
							},
						},
					},
				},
			},
		},
		{
			name: "ReadingClause",
			pos:  position{line: 78, col: 1, offset: 1347},
			expr: &actionExpr{
				pos: position{line: 78, col: 18, offset: 1364},
				run: (*parser).callonReadingClause1,
				expr: &labeledExpr{
					pos:   position{line: 78, col: 18, offset: 1364},
					label: "match",
					expr: &ruleRefExpr{
						pos:  position{line: 78, col: 24, offset: 1370},
						name: "Match",
					},
				},
			},
		},
		{
			name: "Return",
			pos:  position{line: 86, col: 1, offset: 1483},
			expr: &actionExpr{
				pos: position{line: 86, col: 11, offset: 1493},
				run: (*parser).callonReturn1,
				expr: &seqExpr{
					pos: position{line: 86, col: 11, offset: 1493},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 86, col: 11, offset: 1493},
							name: "R",
						},
						&ruleRefExpr{
							pos:  position{line: 86, col: 13, offset: 1495},
							name: "E",
						},
						&ruleRefExpr{
							pos:  position{line: 86, col: 15, offset: 1497},
							name: "T",
						},
						&ruleRefExpr{
							pos:  position{line: 86, col: 17, offset: 1499},
							name: "U",
						},
						&ruleRefExpr{
							pos:  position{line: 86, col: 19, offset: 1501},
							name: "R",
						},
						&ruleRefExpr{
							pos:  position{line: 86, col: 21, offset: 1503},
							name: "N",
						},
						&ruleRefExpr{
							pos:  position{line: 86, col: 24, offset: 1506},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 86, col: 26, offset: 1508},
							label: "v",
							expr: &ruleRefExpr{
								pos:  position{line: 86, col: 28, offset: 1510},
								name: "Variable",
							},
						},
					},
				},
			},
		},
		{
			name: "Match",
			pos:  position{line: 90, col: 1, offset: 1555},
			expr: &actionExpr{
				pos: position{line: 90, col: 10, offset: 1564},
				run: (*parser).callonMatch1,
				expr: &seqExpr{
					pos: position{line: 90, col: 10, offset: 1564},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 90, col: 10, offset: 1564},
							name: "M",
						},
						&ruleRefExpr{
							pos:  position{line: 90, col: 12, offset: 1566},
							name: "A",
						},
						&ruleRefExpr{
							pos:  position{line: 90, col: 14, offset: 1568},
							name: "T",
						},
						&ruleRefExpr{
							pos:  position{line: 90, col: 16, offset: 1570},
							name: "C",
						},
						&ruleRefExpr{
							pos:  position{line: 90, col: 18, offset: 1572},
							name: "H",
						},
						&ruleRefExpr{
							pos:  position{line: 90, col: 20, offset: 1574},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 90, col: 22, offset: 1576},
							label: "pattern",
							expr: &ruleRefExpr{
								pos:  position{line: 90, col: 30, offset: 1584},
								name: "Pattern",
							},
						},
					},
				},
			},
		},
		{
			name: "Pattern",
			pos:  position{line: 95, col: 1, offset: 1687},
			expr: &ruleRefExpr{
				pos:  position{line: 95, col: 12, offset: 1698},
				name: "PatternPart",
			},
		},
		{
			name: "PatternPart",
			pos:  position{line: 97, col: 1, offset: 1714},
			expr: &ruleRefExpr{
				pos:  position{line: 97, col: 16, offset: 1729},
				name: "AnonymousPatternPart",
			},
		},
		{
			name: "AnonymousPatternPart",
			pos:  position{line: 99, col: 1, offset: 1753},
			expr: &ruleRefExpr{
				pos:  position{line: 99, col: 25, offset: 1777},
				name: "PatternElement",
			},
		},
		{
			name: "PatternElement",
			pos:  position{line: 101, col: 1, offset: 1795},
			expr: &ruleRefExpr{
				pos:  position{line: 101, col: 19, offset: 1813},
				name: "NodePattern",
			},
		},
		{
			name: "NodePattern",
			pos:  position{line: 103, col: 1, offset: 1828},
			expr: &actionExpr{
				pos: position{line: 103, col: 16, offset: 1843},
				run: (*parser).callonNodePattern1,
				expr: &seqExpr{
					pos: position{line: 103, col: 16, offset: 1843},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 103, col: 16, offset: 1843},
							val:        "(",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 103, col: 20, offset: 1847},
							label: "variable",
							expr: &ruleRefExpr{
								pos:  position{line: 103, col: 29, offset: 1856},
								name: "Variable",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 103, col: 38, offset: 1865},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 103, col: 40, offset: 1867},
							label: "labels",
							expr: &zeroOrOneExpr{
								pos: position{line: 103, col: 47, offset: 1874},
								expr: &ruleRefExpr{
									pos:  position{line: 103, col: 47, offset: 1874},
									name: "NodeLabels",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 103, col: 59, offset: 1886},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 103, col: 61, offset: 1888},
							label: "props",
							expr: &zeroOrOneExpr{
								pos: position{line: 103, col: 67, offset: 1894},
								expr: &ruleRefExpr{
									pos:  position{line: 103, col: 68, offset: 1895},
									name: "Properties",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 103, col: 81, offset: 1908},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 103, col: 83, offset: 1910},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "NodeLabels",
			pos:  position{line: 119, col: 1, offset: 2174},
			expr: &actionExpr{
				pos: position{line: 119, col: 15, offset: 2188},
				run: (*parser).callonNodeLabels1,
				expr: &seqExpr{
					pos: position{line: 119, col: 15, offset: 2188},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 119, col: 15, offset: 2188},
							label: "label",
							expr: &ruleRefExpr{
								pos:  position{line: 119, col: 21, offset: 2194},
								name: "NodeLabel",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 119, col: 31, offset: 2204},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 119, col: 33, offset: 2206},
							label: "labels",
							expr: &zeroOrMoreExpr{
								pos: position{line: 119, col: 40, offset: 2213},
								expr: &ruleRefExpr{
									pos:  position{line: 119, col: 41, offset: 2214},
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
			pos:  position{line: 130, col: 1, offset: 2433},
			expr: &actionExpr{
				pos: position{line: 130, col: 14, offset: 2446},
				run: (*parser).callonNodeLabel1,
				expr: &seqExpr{
					pos: position{line: 130, col: 14, offset: 2446},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 130, col: 14, offset: 2446},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 130, col: 18, offset: 2450},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 130, col: 20, offset: 2452},
							label: "label",
							expr: &ruleRefExpr{
								pos:  position{line: 130, col: 26, offset: 2458},
								name: "String",
							},
						},
					},
				},
			},
		},
		{
			name: "Variable",
			pos:  position{line: 134, col: 1, offset: 2496},
			expr: &ruleRefExpr{
				pos:  position{line: 134, col: 13, offset: 2508},
				name: "SymbolicName",
			},
		},
		{
			name: "SymbolicName",
			pos:  position{line: 136, col: 1, offset: 2524},
			expr: &ruleRefExpr{
				pos:  position{line: 136, col: 17, offset: 2540},
				name: "String",
			},
		},
		{
			name: "Properties",
			pos:  position{line: 138, col: 1, offset: 2550},
			expr: &ruleRefExpr{
				pos:  position{line: 138, col: 15, offset: 2564},
				name: "MapLiteral",
			},
		},
		{
			name: "ProperyKV",
			pos:  position{line: 139, col: 1, offset: 2576},
			expr: &actionExpr{
				pos: position{line: 139, col: 14, offset: 2589},
				run: (*parser).callonProperyKV1,
				expr: &seqExpr{
					pos: position{line: 139, col: 14, offset: 2589},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 139, col: 14, offset: 2589},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 139, col: 18, offset: 2593},
								name: "String",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 139, col: 25, offset: 2600},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 139, col: 27, offset: 2602},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 139, col: 31, offset: 2606},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 139, col: 33, offset: 2608},
							label: "value",
							expr: &choiceExpr{
								pos: position{line: 139, col: 40, offset: 2615},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 139, col: 40, offset: 2615},
										name: "StringLiteral",
									},
									&ruleRefExpr{
										pos:  position{line: 139, col: 54, offset: 2629},
										name: "Integer",
									},
									&ruleRefExpr{
										pos:  position{line: 139, col: 62, offset: 2637},
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
			pos:  position{line: 152, col: 1, offset: 3050},
			expr: &actionExpr{
				pos: position{line: 152, col: 15, offset: 3064},
				run: (*parser).callonMapLiteral1,
				expr: &seqExpr{
					pos: position{line: 152, col: 15, offset: 3064},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 152, col: 15, offset: 3064},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 152, col: 19, offset: 3068},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 152, col: 21, offset: 3070},
							label: "kv",
							expr: &zeroOrOneExpr{
								pos: position{line: 152, col: 24, offset: 3073},
								expr: &seqExpr{
									pos: position{line: 152, col: 25, offset: 3074},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 152, col: 25, offset: 3074},
											name: "ProperyKV",
										},
										&ruleRefExpr{
											pos:  position{line: 152, col: 35, offset: 3084},
											name: "_",
										},
										&zeroOrMoreExpr{
											pos: position{line: 152, col: 37, offset: 3086},
											expr: &seqExpr{
												pos: position{line: 152, col: 38, offset: 3087},
												exprs: []interface{}{
													&litMatcher{
														pos:        position{line: 152, col: 38, offset: 3087},
														val:        ",",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 152, col: 42, offset: 3091},
														name: "_",
													},
													&ruleRefExpr{
														pos:  position{line: 152, col: 44, offset: 3093},
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
							pos:  position{line: 152, col: 59, offset: 3108},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 152, col: 61, offset: 3110},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "StringLiteral",
			pos:  position{line: 176, col: 1, offset: 3646},
			expr: &actionExpr{
				pos: position{line: 176, col: 18, offset: 3663},
				run: (*parser).callonStringLiteral1,
				expr: &choiceExpr{
					pos: position{line: 176, col: 19, offset: 3664},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 176, col: 19, offset: 3664},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 176, col: 19, offset: 3664},
									val:        "\"",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 176, col: 23, offset: 3668},
									expr: &choiceExpr{
										pos: position{line: 176, col: 25, offset: 3670},
										alternatives: []interface{}{
											&seqExpr{
												pos: position{line: 176, col: 25, offset: 3670},
												exprs: []interface{}{
													&notExpr{
														pos: position{line: 176, col: 25, offset: 3670},
														expr: &ruleRefExpr{
															pos:  position{line: 176, col: 26, offset: 3671},
															name: "EscapedChar",
														},
													},
													&anyMatcher{
														line: 176, col: 38, offset: 3683,
													},
												},
											},
											&seqExpr{
												pos: position{line: 176, col: 42, offset: 3687},
												exprs: []interface{}{
													&litMatcher{
														pos:        position{line: 176, col: 42, offset: 3687},
														val:        "\\",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 176, col: 47, offset: 3692},
														name: "EscapeSequence",
													},
												},
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 176, col: 65, offset: 3710},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 176, col: 71, offset: 3716},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 176, col: 71, offset: 3716},
									val:        "'",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 176, col: 75, offset: 3720},
									expr: &choiceExpr{
										pos: position{line: 176, col: 77, offset: 3722},
										alternatives: []interface{}{
											&seqExpr{
												pos: position{line: 176, col: 77, offset: 3722},
												exprs: []interface{}{
													&notExpr{
														pos: position{line: 176, col: 77, offset: 3722},
														expr: &ruleRefExpr{
															pos:  position{line: 176, col: 78, offset: 3723},
															name: "EscapedChar",
														},
													},
													&anyMatcher{
														line: 176, col: 90, offset: 3735,
													},
												},
											},
											&seqExpr{
												pos: position{line: 176, col: 94, offset: 3739},
												exprs: []interface{}{
													&litMatcher{
														pos:        position{line: 176, col: 94, offset: 3739},
														val:        "\\",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 176, col: 99, offset: 3744},
														name: "EscapeSequence",
													},
												},
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 176, col: 117, offset: 3762},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "EscapedChar",
			pos:  position{line: 191, col: 1, offset: 4249},
			expr: &charClassMatcher{
				pos:        position{line: 191, col: 16, offset: 4264},
				val:        "[\\x00-\\x1f'\"\\\\]",
				chars:      []rune{'\'', '"', '\\'},
				ranges:     []rune{'\x00', '\x1f'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EscapeSequence",
			pos:  position{line: 193, col: 1, offset: 4283},
			expr: &choiceExpr{
				pos: position{line: 193, col: 19, offset: 4301},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 193, col: 19, offset: 4301},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 193, col: 38, offset: 4320},
						name: "UnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 195, col: 1, offset: 4337},
			expr: &charClassMatcher{
				pos:        position{line: 195, col: 21, offset: 4357},
				val:        "['\"\\\\/bfnrt]",
				chars:      []rune{'\'', '"', '\\', '/', 'b', 'f', 'n', 'r', 't'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeEscape",
			pos:  position{line: 197, col: 1, offset: 4373},
			expr: &seqExpr{
				pos: position{line: 197, col: 18, offset: 4390},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 197, col: 18, offset: 4390},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 22, offset: 4394},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 31, offset: 4403},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 40, offset: 4412},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 49, offset: 4421},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "String",
			pos:  position{line: 199, col: 1, offset: 4433},
			expr: &actionExpr{
				pos: position{line: 199, col: 11, offset: 4443},
				run: (*parser).callonString1,
				expr: &oneOrMoreExpr{
					pos: position{line: 199, col: 11, offset: 4443},
					expr: &charClassMatcher{
						pos:        position{line: 199, col: 11, offset: 4443},
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
			pos:  position{line: 203, col: 1, offset: 4496},
			expr: &actionExpr{
				pos: position{line: 203, col: 12, offset: 4507},
				run: (*parser).callonInteger1,
				expr: &oneOrMoreExpr{
					pos: position{line: 203, col: 12, offset: 4507},
					expr: &charClassMatcher{
						pos:        position{line: 203, col: 12, offset: 4507},
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
			pos:  position{line: 207, col: 1, offset: 4575},
			expr: &choiceExpr{
				pos: position{line: 207, col: 16, offset: 4590},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 207, col: 16, offset: 4590},
						run: (*parser).callonBoolLiteral2,
						expr: &seqExpr{
							pos: position{line: 207, col: 16, offset: 4590},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 207, col: 16, offset: 4590},
									name: "T",
								},
								&litMatcher{
									pos:        position{line: 207, col: 18, offset: 4592},
									val:        "rue",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 207, col: 47, offset: 4621},
						run: (*parser).callonBoolLiteral6,
						expr: &seqExpr{
							pos: position{line: 207, col: 47, offset: 4621},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 207, col: 47, offset: 4621},
									name: "F",
								},
								&litMatcher{
									pos:        position{line: 207, col: 49, offset: 4623},
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
			pos:         position{line: 209, col: 1, offset: 4654},
			expr: &zeroOrMoreExpr{
				pos: position{line: 209, col: 18, offset: 4673},
				expr: &charClassMatcher{
					pos:        position{line: 209, col: 18, offset: 4673},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "A",
			pos:  position{line: 211, col: 1, offset: 4687},
			expr: &choiceExpr{
				pos: position{line: 211, col: 7, offset: 4693},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 211, col: 7, offset: 4693},
						val:        "A",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 211, col: 13, offset: 4699},
						val:        "a",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "B",
			pos:  position{line: 212, col: 1, offset: 4705},
			expr: &choiceExpr{
				pos: position{line: 212, col: 7, offset: 4711},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 212, col: 7, offset: 4711},
						val:        "B",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 212, col: 13, offset: 4717},
						val:        "b",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "C",
			pos:  position{line: 213, col: 1, offset: 4723},
			expr: &choiceExpr{
				pos: position{line: 213, col: 7, offset: 4729},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 213, col: 7, offset: 4729},
						val:        "C",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 213, col: 13, offset: 4735},
						val:        "c",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "D",
			pos:  position{line: 214, col: 1, offset: 4741},
			expr: &choiceExpr{
				pos: position{line: 214, col: 7, offset: 4747},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 214, col: 7, offset: 4747},
						val:        "D",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 214, col: 13, offset: 4753},
						val:        "d",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "E",
			pos:  position{line: 215, col: 1, offset: 4759},
			expr: &choiceExpr{
				pos: position{line: 215, col: 7, offset: 4765},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 215, col: 7, offset: 4765},
						val:        "E",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 215, col: 13, offset: 4771},
						val:        "e",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "F",
			pos:  position{line: 216, col: 1, offset: 4777},
			expr: &choiceExpr{
				pos: position{line: 216, col: 7, offset: 4783},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 216, col: 7, offset: 4783},
						val:        "F",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 216, col: 13, offset: 4789},
						val:        "f",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "G",
			pos:  position{line: 217, col: 1, offset: 4795},
			expr: &choiceExpr{
				pos: position{line: 217, col: 7, offset: 4801},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 217, col: 7, offset: 4801},
						val:        "G",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 217, col: 13, offset: 4807},
						val:        "g",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "H",
			pos:  position{line: 218, col: 1, offset: 4813},
			expr: &choiceExpr{
				pos: position{line: 218, col: 7, offset: 4819},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 218, col: 7, offset: 4819},
						val:        "H",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 218, col: 13, offset: 4825},
						val:        "h",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "I",
			pos:  position{line: 219, col: 1, offset: 4831},
			expr: &choiceExpr{
				pos: position{line: 219, col: 7, offset: 4837},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 219, col: 7, offset: 4837},
						val:        "I",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 219, col: 13, offset: 4843},
						val:        "i",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "K",
			pos:  position{line: 220, col: 1, offset: 4849},
			expr: &choiceExpr{
				pos: position{line: 220, col: 7, offset: 4855},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 220, col: 7, offset: 4855},
						val:        "K",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 220, col: 13, offset: 4861},
						val:        "k",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "L",
			pos:  position{line: 221, col: 1, offset: 4867},
			expr: &choiceExpr{
				pos: position{line: 221, col: 7, offset: 4873},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 221, col: 7, offset: 4873},
						val:        "L",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 221, col: 13, offset: 4879},
						val:        "l",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "M",
			pos:  position{line: 222, col: 1, offset: 4885},
			expr: &choiceExpr{
				pos: position{line: 222, col: 7, offset: 4891},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 222, col: 7, offset: 4891},
						val:        "M",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 222, col: 13, offset: 4897},
						val:        "m",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "N",
			pos:  position{line: 223, col: 1, offset: 4903},
			expr: &choiceExpr{
				pos: position{line: 223, col: 7, offset: 4909},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 223, col: 7, offset: 4909},
						val:        "N",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 223, col: 13, offset: 4915},
						val:        "n",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "O",
			pos:  position{line: 224, col: 1, offset: 4921},
			expr: &choiceExpr{
				pos: position{line: 224, col: 7, offset: 4927},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 224, col: 7, offset: 4927},
						val:        "O",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 224, col: 13, offset: 4933},
						val:        "o",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "P",
			pos:  position{line: 225, col: 1, offset: 4939},
			expr: &choiceExpr{
				pos: position{line: 225, col: 7, offset: 4945},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 225, col: 7, offset: 4945},
						val:        "P",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 225, col: 13, offset: 4951},
						val:        "p",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Q",
			pos:  position{line: 226, col: 1, offset: 4957},
			expr: &choiceExpr{
				pos: position{line: 226, col: 7, offset: 4963},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 226, col: 7, offset: 4963},
						val:        "Q",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 226, col: 13, offset: 4969},
						val:        "q",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "R",
			pos:  position{line: 227, col: 1, offset: 4975},
			expr: &choiceExpr{
				pos: position{line: 227, col: 7, offset: 4981},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 227, col: 7, offset: 4981},
						val:        "R",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 227, col: 13, offset: 4987},
						val:        "r",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "S",
			pos:  position{line: 228, col: 1, offset: 4993},
			expr: &choiceExpr{
				pos: position{line: 228, col: 7, offset: 4999},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 228, col: 7, offset: 4999},
						val:        "S",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 228, col: 13, offset: 5005},
						val:        "s",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "T",
			pos:  position{line: 229, col: 1, offset: 5011},
			expr: &choiceExpr{
				pos: position{line: 229, col: 7, offset: 5017},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 229, col: 7, offset: 5017},
						val:        "T",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 229, col: 13, offset: 5023},
						val:        "t",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "U",
			pos:  position{line: 230, col: 1, offset: 5029},
			expr: &choiceExpr{
				pos: position{line: 230, col: 7, offset: 5035},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 230, col: 7, offset: 5035},
						val:        "U",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 230, col: 13, offset: 5041},
						val:        "u",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "V",
			pos:  position{line: 231, col: 1, offset: 5047},
			expr: &choiceExpr{
				pos: position{line: 231, col: 7, offset: 5053},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 231, col: 7, offset: 5053},
						val:        "V",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 231, col: 13, offset: 5059},
						val:        "v",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "W",
			pos:  position{line: 232, col: 1, offset: 5065},
			expr: &choiceExpr{
				pos: position{line: 232, col: 7, offset: 5071},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 232, col: 7, offset: 5071},
						val:        "W",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 232, col: 13, offset: 5077},
						val:        "w",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "X",
			pos:  position{line: 233, col: 1, offset: 5083},
			expr: &choiceExpr{
				pos: position{line: 233, col: 7, offset: 5089},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 233, col: 7, offset: 5089},
						val:        "X",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 233, col: 13, offset: 5095},
						val:        "x",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Y",
			pos:  position{line: 234, col: 1, offset: 5101},
			expr: &choiceExpr{
				pos: position{line: 234, col: 7, offset: 5107},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 234, col: 7, offset: 5107},
						val:        "Y",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 234, col: 13, offset: 5113},
						val:        "y",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 236, col: 1, offset: 5121},
			expr: &notExpr{
				pos: position{line: 236, col: 8, offset: 5128},
				expr: &anyMatcher{
					line: 236, col: 9, offset: 5129,
				},
			},
		},
	},
}

func (c *current) onStatement1(query interface{}) (interface{}, error) {

	log.Printf("%s", c.text)
	q := QueryPlan{
		ReadingClause: []ReadingClause{query.(ReadingClause)},
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

func (c *current) onSingleQuery1(clause, ret interface{}) (interface{}, error) {

	cl := clause.(ReadingClause)
	cl.Returns = []string{ret.(string)}
	return cl, nil
}

func (p *parser) callonSingleQuery1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSingleQuery1(stack["clause"], stack["ret"])
}

func (c *current) onReadingClause1(match interface{}) (interface{}, error) {

	clause := ReadingClause{
		Match: match.(MatchQuery),
	}

	return clause, nil
}

func (p *parser) callonReadingClause1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReadingClause1(stack["match"])
}

func (c *current) onReturn1(v interface{}) (interface{}, error) {

	return v.(string), nil
}

func (p *parser) callonReturn1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReturn1(stack["v"])
}

func (c *current) onMatch1(pattern interface{}) (interface{}, error) {

	plan := MatchQuery{Nodes: []NodeQuery{pattern.(NodeQuery)}}
	return plan, nil
}

func (p *parser) callonMatch1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMatch1(stack["pattern"])
}

func (c *current) onNodePattern1(variable, labels, props interface{}) (interface{}, error) {

	plan := NodeQuery{
		Variable: variable.(string),
	}

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
	return p.cur.onNodePattern1(stack["variable"], stack["labels"], stack["props"])
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
		if _, ok := props[kv.Key]; ok {
			return nil, fmt.Errorf("Duplicate map key %s not allowed", kv.Key)
		}
		props[kv.Key] = kv.Value
	}

	return props, nil
}

func (p *parser) callonMapLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMapLiteral1(stack["kv"])
}

func (c *current) onStringLiteral1() (interface{}, error) {

	c.text = bytes.Replace(c.text, []byte(`\/`), []byte(`/`), -1)

	// deal with single quates
	if strings.HasPrefix(string(c.text), "'") {
		c.text = bytes.ReplaceAll(c.text, []byte(`"`), []byte(`'`))
		c.text = bytes.ReplaceAll(c.text, []byte(`'`), []byte(``))
		return string(c.text), nil
	}

	// deal with double quates
	c.text = bytes.ReplaceAll(c.text, []byte(`'`), []byte(`\'`))
	return strconv.Unquote(string(c.text))
}

func (p *parser) callonStringLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringLiteral1()
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
