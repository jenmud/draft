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

var g = &grammar{
	rules: []*rule{
		{
			name: "Statement",
			pos:  position{line: 6, col: 1, offset: 84},
			expr: &actionExpr{
				pos: position{line: 6, col: 14, offset: 97},
				run: (*parser).callonStatement1,
				expr: &seqExpr{
					pos: position{line: 6, col: 14, offset: 97},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 6, col: 14, offset: 97},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 6, col: 16, offset: 99},
							label: "query",
							expr: &ruleRefExpr{
								pos:  position{line: 6, col: 22, offset: 105},
								name: "Query",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 6, col: 28, offset: 111},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 6, col: 30, offset: 113},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Query",
			pos:  position{line: 15, col: 1, offset: 275},
			expr: &actionExpr{
				pos: position{line: 15, col: 10, offset: 284},
				run: (*parser).callonQuery1,
				expr: &labeledExpr{
					pos:   position{line: 15, col: 10, offset: 284},
					label: "regularQuery",
					expr: &ruleRefExpr{
						pos:  position{line: 15, col: 23, offset: 297},
						name: "RegularQuery",
					},
				},
			},
		},
		{
			name: "RegularQuery",
			pos:  position{line: 20, col: 1, offset: 385},
			expr: &actionExpr{
				pos: position{line: 20, col: 18, offset: 402},
				run: (*parser).callonRegularQuery1,
				expr: &labeledExpr{
					pos:   position{line: 20, col: 18, offset: 402},
					label: "singleQuery",
					expr: &ruleRefExpr{
						pos:  position{line: 20, col: 30, offset: 414},
						name: "SingleQuery",
					},
				},
			},
		},
		{
			name: "SingleQuery",
			pos:  position{line: 25, col: 1, offset: 507},
			expr: &actionExpr{
				pos: position{line: 25, col: 16, offset: 522},
				run: (*parser).callonSingleQuery1,
				expr: &seqExpr{
					pos: position{line: 25, col: 16, offset: 522},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 25, col: 16, offset: 522},
							label: "matches",
							expr: &oneOrMoreExpr{
								pos: position{line: 25, col: 24, offset: 530},
								expr: &seqExpr{
									pos: position{line: 25, col: 25, offset: 531},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 25, col: 25, offset: 531},
											name: "ReadingClause",
										},
										&ruleRefExpr{
											pos:  position{line: 25, col: 39, offset: 545},
											name: "_",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 25, col: 43, offset: 549},
							label: "returns",
							expr: &ruleRefExpr{
								pos:  position{line: 25, col: 51, offset: 557},
								name: "Return",
							},
						},
					},
				},
			},
		},
		{
			name: "ReadingClause",
			pos:  position{line: 44, col: 1, offset: 1002},
			expr: &actionExpr{
				pos: position{line: 44, col: 18, offset: 1019},
				run: (*parser).callonReadingClause1,
				expr: &labeledExpr{
					pos:   position{line: 44, col: 18, offset: 1019},
					label: "match",
					expr: &ruleRefExpr{
						pos:  position{line: 44, col: 24, offset: 1025},
						name: "Match",
					},
				},
			},
		},
		{
			name: "Return",
			pos:  position{line: 49, col: 1, offset: 1115},
			expr: &actionExpr{
				pos: position{line: 49, col: 11, offset: 1125},
				run: (*parser).callonReturn1,
				expr: &seqExpr{
					pos: position{line: 49, col: 11, offset: 1125},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 49, col: 11, offset: 1125},
							name: "R",
						},
						&ruleRefExpr{
							pos:  position{line: 49, col: 13, offset: 1127},
							name: "E",
						},
						&ruleRefExpr{
							pos:  position{line: 49, col: 15, offset: 1129},
							name: "T",
						},
						&ruleRefExpr{
							pos:  position{line: 49, col: 17, offset: 1131},
							name: "U",
						},
						&ruleRefExpr{
							pos:  position{line: 49, col: 19, offset: 1133},
							name: "R",
						},
						&ruleRefExpr{
							pos:  position{line: 49, col: 21, offset: 1135},
							name: "N",
						},
						&ruleRefExpr{
							pos:  position{line: 49, col: 24, offset: 1138},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 49, col: 26, offset: 1140},
							label: "variable",
							expr: &ruleRefExpr{
								pos:  position{line: 49, col: 35, offset: 1149},
								name: "Variable",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 49, col: 44, offset: 1158},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 49, col: 46, offset: 1160},
							label: "extra",
							expr: &zeroOrMoreExpr{
								pos: position{line: 49, col: 52, offset: 1166},
								expr: &seqExpr{
									pos: position{line: 49, col: 53, offset: 1167},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 49, col: 53, offset: 1167},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 49, col: 57, offset: 1171},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 49, col: 59, offset: 1173},
											name: "Variable",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Match",
			pos:  position{line: 61, col: 1, offset: 1528},
			expr: &actionExpr{
				pos: position{line: 61, col: 10, offset: 1537},
				run: (*parser).callonMatch1,
				expr: &seqExpr{
					pos: position{line: 61, col: 10, offset: 1537},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 61, col: 10, offset: 1537},
							name: "M",
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 12, offset: 1539},
							name: "A",
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 14, offset: 1541},
							name: "T",
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 16, offset: 1543},
							name: "C",
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 18, offset: 1545},
							name: "H",
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 20, offset: 1547},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 61, col: 22, offset: 1549},
							label: "pattern",
							expr: &ruleRefExpr{
								pos:  position{line: 61, col: 30, offset: 1557},
								name: "Pattern",
							},
						},
					},
				},
			},
		},
		{
			name: "Pattern",
			pos:  position{line: 67, col: 1, offset: 1684},
			expr: &ruleRefExpr{
				pos:  position{line: 67, col: 12, offset: 1695},
				name: "PatternPart",
			},
		},
		{
			name: "PatternPart",
			pos:  position{line: 69, col: 1, offset: 1711},
			expr: &ruleRefExpr{
				pos:  position{line: 69, col: 16, offset: 1726},
				name: "AnonymousPatternPart",
			},
		},
		{
			name: "AnonymousPatternPart",
			pos:  position{line: 71, col: 1, offset: 1750},
			expr: &ruleRefExpr{
				pos:  position{line: 71, col: 25, offset: 1774},
				name: "PatternElement",
			},
		},
		{
			name: "PatternElement",
			pos:  position{line: 73, col: 1, offset: 1792},
			expr: &ruleRefExpr{
				pos:  position{line: 73, col: 19, offset: 1810},
				name: "NodePattern",
			},
		},
		{
			name: "NodePattern",
			pos:  position{line: 75, col: 1, offset: 1825},
			expr: &actionExpr{
				pos: position{line: 75, col: 16, offset: 1840},
				run: (*parser).callonNodePattern1,
				expr: &seqExpr{
					pos: position{line: 75, col: 16, offset: 1840},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 75, col: 16, offset: 1840},
							val:        "(",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 75, col: 20, offset: 1844},
							label: "variable",
							expr: &ruleRefExpr{
								pos:  position{line: 75, col: 29, offset: 1853},
								name: "Variable",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 75, col: 38, offset: 1862},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 75, col: 40, offset: 1864},
							label: "labels",
							expr: &zeroOrOneExpr{
								pos: position{line: 75, col: 47, offset: 1871},
								expr: &ruleRefExpr{
									pos:  position{line: 75, col: 47, offset: 1871},
									name: "NodeLabels",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 75, col: 59, offset: 1883},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 75, col: 61, offset: 1885},
							label: "props",
							expr: &zeroOrOneExpr{
								pos: position{line: 75, col: 67, offset: 1891},
								expr: &ruleRefExpr{
									pos:  position{line: 75, col: 68, offset: 1892},
									name: "Properties",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 75, col: 81, offset: 1905},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 75, col: 83, offset: 1907},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "NodeLabels",
			pos:  position{line: 91, col: 1, offset: 2166},
			expr: &actionExpr{
				pos: position{line: 91, col: 15, offset: 2180},
				run: (*parser).callonNodeLabels1,
				expr: &seqExpr{
					pos: position{line: 91, col: 15, offset: 2180},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 91, col: 15, offset: 2180},
							label: "label",
							expr: &ruleRefExpr{
								pos:  position{line: 91, col: 21, offset: 2186},
								name: "NodeLabel",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 91, col: 31, offset: 2196},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 91, col: 33, offset: 2198},
							label: "labels",
							expr: &zeroOrMoreExpr{
								pos: position{line: 91, col: 40, offset: 2205},
								expr: &ruleRefExpr{
									pos:  position{line: 91, col: 41, offset: 2206},
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
			pos:  position{line: 108, col: 1, offset: 2548},
			expr: &actionExpr{
				pos: position{line: 108, col: 14, offset: 2561},
				run: (*parser).callonNodeLabel1,
				expr: &seqExpr{
					pos: position{line: 108, col: 14, offset: 2561},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 108, col: 14, offset: 2561},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 108, col: 18, offset: 2565},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 108, col: 20, offset: 2567},
							label: "label",
							expr: &ruleRefExpr{
								pos:  position{line: 108, col: 26, offset: 2573},
								name: "String",
							},
						},
					},
				},
			},
		},
		{
			name: "Variable",
			pos:  position{line: 112, col: 1, offset: 2611},
			expr: &ruleRefExpr{
				pos:  position{line: 112, col: 13, offset: 2623},
				name: "SymbolicName",
			},
		},
		{
			name: "SymbolicName",
			pos:  position{line: 114, col: 1, offset: 2639},
			expr: &ruleRefExpr{
				pos:  position{line: 114, col: 17, offset: 2655},
				name: "String",
			},
		},
		{
			name: "Properties",
			pos:  position{line: 116, col: 1, offset: 2665},
			expr: &ruleRefExpr{
				pos:  position{line: 116, col: 15, offset: 2679},
				name: "MapLiteral",
			},
		},
		{
			name: "ProperyKV",
			pos:  position{line: 117, col: 1, offset: 2691},
			expr: &actionExpr{
				pos: position{line: 117, col: 14, offset: 2704},
				run: (*parser).callonProperyKV1,
				expr: &seqExpr{
					pos: position{line: 117, col: 14, offset: 2704},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 117, col: 14, offset: 2704},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 117, col: 18, offset: 2708},
								name: "String",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 117, col: 25, offset: 2715},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 117, col: 27, offset: 2717},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 117, col: 31, offset: 2721},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 117, col: 33, offset: 2723},
							label: "value",
							expr: &choiceExpr{
								pos: position{line: 117, col: 40, offset: 2730},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 117, col: 40, offset: 2730},
										name: "StringLiteral",
									},
									&ruleRefExpr{
										pos:  position{line: 117, col: 54, offset: 2744},
										name: "Integer",
									},
									&ruleRefExpr{
										pos:  position{line: 117, col: 62, offset: 2752},
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
			pos:  position{line: 130, col: 1, offset: 3165},
			expr: &actionExpr{
				pos: position{line: 130, col: 15, offset: 3179},
				run: (*parser).callonMapLiteral1,
				expr: &seqExpr{
					pos: position{line: 130, col: 15, offset: 3179},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 130, col: 15, offset: 3179},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 130, col: 19, offset: 3183},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 130, col: 21, offset: 3185},
							label: "kv",
							expr: &zeroOrOneExpr{
								pos: position{line: 130, col: 24, offset: 3188},
								expr: &seqExpr{
									pos: position{line: 130, col: 25, offset: 3189},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 130, col: 25, offset: 3189},
											name: "ProperyKV",
										},
										&ruleRefExpr{
											pos:  position{line: 130, col: 35, offset: 3199},
											name: "_",
										},
										&zeroOrMoreExpr{
											pos: position{line: 130, col: 37, offset: 3201},
											expr: &seqExpr{
												pos: position{line: 130, col: 38, offset: 3202},
												exprs: []interface{}{
													&litMatcher{
														pos:        position{line: 130, col: 38, offset: 3202},
														val:        ",",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 130, col: 42, offset: 3206},
														name: "_",
													},
													&ruleRefExpr{
														pos:  position{line: 130, col: 44, offset: 3208},
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
							pos:  position{line: 130, col: 59, offset: 3223},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 130, col: 61, offset: 3225},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "StringLiteral",
			pos:  position{line: 154, col: 1, offset: 3761},
			expr: &actionExpr{
				pos: position{line: 154, col: 18, offset: 3778},
				run: (*parser).callonStringLiteral1,
				expr: &choiceExpr{
					pos: position{line: 154, col: 19, offset: 3779},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 154, col: 19, offset: 3779},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 154, col: 19, offset: 3779},
									val:        "\"",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 154, col: 23, offset: 3783},
									expr: &choiceExpr{
										pos: position{line: 154, col: 25, offset: 3785},
										alternatives: []interface{}{
											&seqExpr{
												pos: position{line: 154, col: 25, offset: 3785},
												exprs: []interface{}{
													&notExpr{
														pos: position{line: 154, col: 25, offset: 3785},
														expr: &ruleRefExpr{
															pos:  position{line: 154, col: 26, offset: 3786},
															name: "EscapedChar",
														},
													},
													&anyMatcher{
														line: 154, col: 38, offset: 3798,
													},
												},
											},
											&seqExpr{
												pos: position{line: 154, col: 42, offset: 3802},
												exprs: []interface{}{
													&litMatcher{
														pos:        position{line: 154, col: 42, offset: 3802},
														val:        "\\",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 154, col: 47, offset: 3807},
														name: "EscapeSequence",
													},
												},
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 154, col: 65, offset: 3825},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 154, col: 71, offset: 3831},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 154, col: 71, offset: 3831},
									val:        "'",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 154, col: 75, offset: 3835},
									expr: &choiceExpr{
										pos: position{line: 154, col: 77, offset: 3837},
										alternatives: []interface{}{
											&seqExpr{
												pos: position{line: 154, col: 77, offset: 3837},
												exprs: []interface{}{
													&notExpr{
														pos: position{line: 154, col: 77, offset: 3837},
														expr: &ruleRefExpr{
															pos:  position{line: 154, col: 78, offset: 3838},
															name: "EscapedChar",
														},
													},
													&anyMatcher{
														line: 154, col: 90, offset: 3850,
													},
												},
											},
											&seqExpr{
												pos: position{line: 154, col: 94, offset: 3854},
												exprs: []interface{}{
													&litMatcher{
														pos:        position{line: 154, col: 94, offset: 3854},
														val:        "\\",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 154, col: 99, offset: 3859},
														name: "EscapeSequence",
													},
												},
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 154, col: 117, offset: 3877},
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
			pos:  position{line: 169, col: 1, offset: 4364},
			expr: &charClassMatcher{
				pos:        position{line: 169, col: 16, offset: 4379},
				val:        "[\\x00-\\x1f'\"\\\\]",
				chars:      []rune{'\'', '"', '\\'},
				ranges:     []rune{'\x00', '\x1f'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EscapeSequence",
			pos:  position{line: 171, col: 1, offset: 4398},
			expr: &choiceExpr{
				pos: position{line: 171, col: 19, offset: 4416},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 171, col: 19, offset: 4416},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 171, col: 38, offset: 4435},
						name: "UnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 173, col: 1, offset: 4452},
			expr: &charClassMatcher{
				pos:        position{line: 173, col: 21, offset: 4472},
				val:        "['\"\\\\/bfnrt]",
				chars:      []rune{'\'', '"', '\\', '/', 'b', 'f', 'n', 'r', 't'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeEscape",
			pos:  position{line: 175, col: 1, offset: 4488},
			expr: &seqExpr{
				pos: position{line: 175, col: 18, offset: 4505},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 175, col: 18, offset: 4505},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 175, col: 22, offset: 4509},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 175, col: 31, offset: 4518},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 175, col: 40, offset: 4527},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 175, col: 49, offset: 4536},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "String",
			pos:  position{line: 177, col: 1, offset: 4548},
			expr: &actionExpr{
				pos: position{line: 177, col: 11, offset: 4558},
				run: (*parser).callonString1,
				expr: &oneOrMoreExpr{
					pos: position{line: 177, col: 11, offset: 4558},
					expr: &charClassMatcher{
						pos:        position{line: 177, col: 11, offset: 4558},
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
			pos:  position{line: 181, col: 1, offset: 4611},
			expr: &actionExpr{
				pos: position{line: 181, col: 12, offset: 4622},
				run: (*parser).callonInteger1,
				expr: &oneOrMoreExpr{
					pos: position{line: 181, col: 12, offset: 4622},
					expr: &charClassMatcher{
						pos:        position{line: 181, col: 12, offset: 4622},
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
			pos:  position{line: 185, col: 1, offset: 4690},
			expr: &choiceExpr{
				pos: position{line: 185, col: 16, offset: 4705},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 185, col: 16, offset: 4705},
						run: (*parser).callonBoolLiteral2,
						expr: &seqExpr{
							pos: position{line: 185, col: 16, offset: 4705},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 185, col: 16, offset: 4705},
									name: "T",
								},
								&litMatcher{
									pos:        position{line: 185, col: 18, offset: 4707},
									val:        "rue",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 185, col: 47, offset: 4736},
						run: (*parser).callonBoolLiteral6,
						expr: &seqExpr{
							pos: position{line: 185, col: 47, offset: 4736},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 185, col: 47, offset: 4736},
									name: "F",
								},
								&litMatcher{
									pos:        position{line: 185, col: 49, offset: 4738},
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
			pos:         position{line: 187, col: 1, offset: 4769},
			expr: &zeroOrMoreExpr{
				pos: position{line: 187, col: 19, offset: 4787},
				expr: &charClassMatcher{
					pos:        position{line: 187, col: 19, offset: 4787},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "A",
			pos:  position{line: 189, col: 1, offset: 4801},
			expr: &choiceExpr{
				pos: position{line: 189, col: 7, offset: 4807},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 189, col: 7, offset: 4807},
						val:        "A",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 189, col: 13, offset: 4813},
						val:        "a",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "B",
			pos:  position{line: 190, col: 1, offset: 4819},
			expr: &choiceExpr{
				pos: position{line: 190, col: 7, offset: 4825},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 190, col: 7, offset: 4825},
						val:        "B",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 190, col: 13, offset: 4831},
						val:        "b",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "C",
			pos:  position{line: 191, col: 1, offset: 4837},
			expr: &choiceExpr{
				pos: position{line: 191, col: 7, offset: 4843},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 191, col: 7, offset: 4843},
						val:        "C",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 191, col: 13, offset: 4849},
						val:        "c",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "D",
			pos:  position{line: 192, col: 1, offset: 4855},
			expr: &choiceExpr{
				pos: position{line: 192, col: 7, offset: 4861},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 192, col: 7, offset: 4861},
						val:        "D",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 192, col: 13, offset: 4867},
						val:        "d",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "E",
			pos:  position{line: 193, col: 1, offset: 4873},
			expr: &choiceExpr{
				pos: position{line: 193, col: 7, offset: 4879},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 193, col: 7, offset: 4879},
						val:        "E",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 193, col: 13, offset: 4885},
						val:        "e",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "F",
			pos:  position{line: 194, col: 1, offset: 4891},
			expr: &choiceExpr{
				pos: position{line: 194, col: 7, offset: 4897},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 194, col: 7, offset: 4897},
						val:        "F",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 194, col: 13, offset: 4903},
						val:        "f",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "G",
			pos:  position{line: 195, col: 1, offset: 4909},
			expr: &choiceExpr{
				pos: position{line: 195, col: 7, offset: 4915},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 195, col: 7, offset: 4915},
						val:        "G",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 195, col: 13, offset: 4921},
						val:        "g",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "H",
			pos:  position{line: 196, col: 1, offset: 4927},
			expr: &choiceExpr{
				pos: position{line: 196, col: 7, offset: 4933},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 196, col: 7, offset: 4933},
						val:        "H",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 196, col: 13, offset: 4939},
						val:        "h",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "I",
			pos:  position{line: 197, col: 1, offset: 4945},
			expr: &choiceExpr{
				pos: position{line: 197, col: 7, offset: 4951},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 197, col: 7, offset: 4951},
						val:        "I",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 197, col: 13, offset: 4957},
						val:        "i",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "K",
			pos:  position{line: 198, col: 1, offset: 4963},
			expr: &choiceExpr{
				pos: position{line: 198, col: 7, offset: 4969},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 198, col: 7, offset: 4969},
						val:        "K",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 198, col: 13, offset: 4975},
						val:        "k",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "L",
			pos:  position{line: 199, col: 1, offset: 4981},
			expr: &choiceExpr{
				pos: position{line: 199, col: 7, offset: 4987},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 199, col: 7, offset: 4987},
						val:        "L",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 199, col: 13, offset: 4993},
						val:        "l",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "M",
			pos:  position{line: 200, col: 1, offset: 4999},
			expr: &choiceExpr{
				pos: position{line: 200, col: 7, offset: 5005},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 200, col: 7, offset: 5005},
						val:        "M",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 200, col: 13, offset: 5011},
						val:        "m",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "N",
			pos:  position{line: 201, col: 1, offset: 5017},
			expr: &choiceExpr{
				pos: position{line: 201, col: 7, offset: 5023},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 201, col: 7, offset: 5023},
						val:        "N",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 201, col: 13, offset: 5029},
						val:        "n",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "O",
			pos:  position{line: 202, col: 1, offset: 5035},
			expr: &choiceExpr{
				pos: position{line: 202, col: 7, offset: 5041},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 202, col: 7, offset: 5041},
						val:        "O",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 202, col: 13, offset: 5047},
						val:        "o",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "P",
			pos:  position{line: 203, col: 1, offset: 5053},
			expr: &choiceExpr{
				pos: position{line: 203, col: 7, offset: 5059},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 203, col: 7, offset: 5059},
						val:        "P",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 203, col: 13, offset: 5065},
						val:        "p",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Q",
			pos:  position{line: 204, col: 1, offset: 5071},
			expr: &choiceExpr{
				pos: position{line: 204, col: 7, offset: 5077},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 204, col: 7, offset: 5077},
						val:        "Q",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 204, col: 13, offset: 5083},
						val:        "q",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "R",
			pos:  position{line: 205, col: 1, offset: 5089},
			expr: &choiceExpr{
				pos: position{line: 205, col: 7, offset: 5095},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 205, col: 7, offset: 5095},
						val:        "R",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 205, col: 13, offset: 5101},
						val:        "r",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "S",
			pos:  position{line: 206, col: 1, offset: 5107},
			expr: &choiceExpr{
				pos: position{line: 206, col: 7, offset: 5113},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 206, col: 7, offset: 5113},
						val:        "S",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 206, col: 13, offset: 5119},
						val:        "s",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "T",
			pos:  position{line: 207, col: 1, offset: 5125},
			expr: &choiceExpr{
				pos: position{line: 207, col: 7, offset: 5131},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 207, col: 7, offset: 5131},
						val:        "T",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 207, col: 13, offset: 5137},
						val:        "t",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "U",
			pos:  position{line: 208, col: 1, offset: 5143},
			expr: &choiceExpr{
				pos: position{line: 208, col: 7, offset: 5149},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 208, col: 7, offset: 5149},
						val:        "U",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 208, col: 13, offset: 5155},
						val:        "u",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "V",
			pos:  position{line: 209, col: 1, offset: 5161},
			expr: &choiceExpr{
				pos: position{line: 209, col: 7, offset: 5167},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 209, col: 7, offset: 5167},
						val:        "V",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 209, col: 13, offset: 5173},
						val:        "v",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "W",
			pos:  position{line: 210, col: 1, offset: 5179},
			expr: &choiceExpr{
				pos: position{line: 210, col: 7, offset: 5185},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 210, col: 7, offset: 5185},
						val:        "W",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 210, col: 13, offset: 5191},
						val:        "w",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "X",
			pos:  position{line: 211, col: 1, offset: 5197},
			expr: &choiceExpr{
				pos: position{line: 211, col: 7, offset: 5203},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 211, col: 7, offset: 5203},
						val:        "X",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 211, col: 13, offset: 5209},
						val:        "x",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Y",
			pos:  position{line: 212, col: 1, offset: 5215},
			expr: &choiceExpr{
				pos: position{line: 212, col: 7, offset: 5221},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 212, col: 7, offset: 5221},
						val:        "Y",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 212, col: 13, offset: 5227},
						val:        "y",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 214, col: 1, offset: 5235},
			expr: &notExpr{
				pos: position{line: 214, col: 8, offset: 5242},
				expr: &anyMatcher{
					line: 214, col: 9, offset: 5243,
				},
			},
		},
	},
}

func (c *current) onStatement1(query interface{}) (interface{}, error) {

	log.Printf("Query: %s", c.text)
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

	log.Printf("Query: %s", c.text)
	return regularQuery, nil
}

func (p *parser) callonQuery1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onQuery1(stack["regularQuery"])
}

func (c *current) onRegularQuery1(singleQuery interface{}) (interface{}, error) {

	log.Printf("RegularQuery: %s", c.text)
	return singleQuery, nil
}

func (p *parser) callonRegularQuery1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRegularQuery1(stack["singleQuery"])
}

func (c *current) onSingleQuery1(matches, returns interface{}) (interface{}, error) {

	log.Printf("SingleQuery: %s", c.text)
	if returns == nil {
		return nil, fmt.Errorf("RETURN missing and is required")
	}

	clause := ReadingClause{
		Matches: []Match{},
		Returns: returns.([]string),
	}

	for _, match := range toIfaceSlice(matches) {
		m := toIfaceSlice(match)
		clause.Matches = append(clause.Matches, m[0].(Match))
	}

	return clause, nil
}

func (p *parser) callonSingleQuery1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSingleQuery1(stack["matches"], stack["returns"])
}

func (c *current) onReadingClause1(match interface{}) (interface{}, error) {

	log.Printf("ReadingClause: %s", c.text)
	return match.(Match), nil
}

func (p *parser) callonReadingClause1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReadingClause1(stack["match"])
}

func (c *current) onReturn1(variable, extra interface{}) (interface{}, error) {

	log.Printf("Return: %s", c.text)
	extras := toIfaceSlice(extra)
	log.Printf("Vars: %#v, extra: %#v", variable, extra)
	variables := []string{variable.(string)}
	for _, r := range extras {
		v := toIfaceSlice(r)
		variables = append(variables, v[len(v)-1].(string))
	}
	return variables, nil
}

func (p *parser) callonReturn1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onReturn1(stack["variable"], stack["extra"])
}

func (c *current) onMatch1(pattern interface{}) (interface{}, error) {

	log.Printf("Match: %s", c.text)
	match := Match{Nodes: []Node{pattern.(Node)}}
	return match, nil
}

func (p *parser) callonMatch1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMatch1(stack["pattern"])
}

func (c *current) onNodePattern1(variable, labels, props interface{}) (interface{}, error) {

	plan := Node{
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

	labelsList := toIfaceSlice(labels)

	l := make([]string, 1+len(labelsList))
	l[0] = label.(string)

	for i, label := range labelsList {
		l[i+1] = label.(string)
	}

	if len(l) > 1 {
		return nil, fmt.Errorf("Draft nodes only support a single label")
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
