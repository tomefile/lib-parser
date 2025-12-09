package libparser_test

import libparser "github.com/tomefile/lib-parser"

type DataTestCase struct {
	Filename string
	Expect   *libparser.NodeRoot
}

var ExpectedData = []DataTestCase{
	{
		Filename: "01_basic.tome",
		Expect: &libparser.NodeRoot{
			Tomes: map[string]libparser.Node{},
			NodeChildren: libparser.NodeChildren{
				&libparser.NodeComment{Contents: " Example program, привет мир 👨‍🚀!"},
				&libparser.NodeDirective{
					Name: "include",
					NodeArgs: libparser.NodeArgs{
						&libparser.NodeString{
							Segments: libparser.SegmentedString{
								&libparser.LiteralStringSegment{Contents: "@std"},
							},
						},
					},
					NodeChildren: libparser.NodeChildren{},
				},
				&libparser.NodeExec{
					Name: "echo",
					NodeArgs: libparser.NodeArgs{
						&libparser.NodeString{
							Segments: libparser.SegmentedString{
								&libparser.LiteralStringSegment{Contents: "Hello World!"},
							},
						},
						&libparser.NodeString{
							Segments: libparser.SegmentedString{
								&libparser.LiteralStringSegment{Contents: "and another line"},
							},
						},
						&libparser.NodeString{
							Segments: libparser.SegmentedString{
								&libparser.LiteralStringSegment{Contents: "and another."},
							},
						},
					},
				},
			},
		},
	},
	{
		Filename: "02_directive_body.tome",
		Expect: &libparser.NodeRoot{
			Tomes: map[string]libparser.Node{},
			NodeChildren: libparser.NodeChildren{
				&libparser.NodeExec{
					Name: "echo",
					NodeArgs: libparser.NodeArgs{
						&libparser.NodeString{
							Segments: libparser.SegmentedString{
								&libparser.LiteralStringSegment{Contents: "1"},
							},
						},
					},
				},
				&libparser.NodeDirective{
					Name: "section",
					NodeArgs: libparser.NodeArgs{
						&libparser.NodeString{
							Segments: libparser.SegmentedString{
								&libparser.LiteralStringSegment{Contents: "Hello World!"},
							},
						},
					},
					NodeChildren: libparser.NodeChildren{
						&libparser.NodeExec{
							Name: "echo",
							NodeArgs: libparser.NodeArgs{
								&libparser.NodeLiteral{Contents: "1.1"},
							},
						},
						&libparser.NodeExec{
							Name: "echo",
							NodeArgs: libparser.NodeArgs{
								&libparser.NodeString{
									Segments: libparser.SegmentedString{
										&libparser.LiteralStringSegment{Contents: "1.2"},
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
		Filename: "03_directive_nested.tome",
		Expect: &libparser.NodeRoot{
			Tomes: map[string]libparser.Node{},
			NodeChildren: libparser.NodeChildren{
				&libparser.NodeExec{
					Name: "echo",
					NodeArgs: libparser.NodeArgs{
						&libparser.NodeString{
							Segments: libparser.SegmentedString{
								&libparser.LiteralStringSegment{Contents: "1"},
							},
						},
					},
				},
				&libparser.NodeDirective{
					Name: "section",
					NodeArgs: libparser.NodeArgs{
						&libparser.NodeString{
							Segments: libparser.SegmentedString{
								&libparser.LiteralStringSegment{Contents: "Hello World!"},
							},
						},
					},
					NodeChildren: libparser.NodeChildren{
						&libparser.NodeExec{
							Name: "echo",
							NodeArgs: libparser.NodeArgs{
								&libparser.NodeLiteral{Contents: "1.1"},
							},
						},
						&libparser.NodeExec{
							Name: "echo",
							NodeArgs: libparser.NodeArgs{
								&libparser.NodeString{
									Segments: libparser.SegmentedString{
										&libparser.LiteralStringSegment{Contents: "1.2"},
									},
								},
							},
						},
						&libparser.NodeDirective{
							Name: "section",
							NodeArgs: libparser.NodeArgs{
								&libparser.NodeString{
									Segments: libparser.SegmentedString{
										&libparser.LiteralStringSegment{Contents: "Nested"},
									},
								},
							},
							NodeChildren: libparser.NodeChildren{
								&libparser.NodeComment{Contents: " This is nested inside"},
								&libparser.NodeExec{
									Name: "echo",
									NodeArgs: libparser.NodeArgs{
										&libparser.NodeString{
											Segments: libparser.SegmentedString{
												&libparser.LiteralStringSegment{Contents: "2.1"},
											},
										},
									},
								},
								&libparser.NodeExec{
									Name: "echo",
									NodeArgs: libparser.NodeArgs{
										&libparser.NodeString{
											Segments: libparser.SegmentedString{
												&libparser.LiteralStringSegment{Contents: "2.2"},
											},
										},
									},
								},
							},
						},
						&libparser.NodeExec{
							Name: "echo",
							NodeArgs: libparser.NodeArgs{
								&libparser.NodeString{
									Segments: libparser.SegmentedString{
										&libparser.LiteralStringSegment{Contents: "1.3"},
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
		Filename: "04_subcommand.tome",
		Expect: &libparser.NodeRoot{
			Tomes: map[string]libparser.Node{},
			NodeChildren: libparser.NodeChildren{
				&libparser.NodeCall{
					Macro: "my_macro",
					NodeArgs: libparser.NodeArgs{
						&libparser.NodeString{
							Segments: libparser.SegmentedString{
								&libparser.LiteralStringSegment{Contents: "123"},
							},
						},
						&libparser.NodeExec{
							Name: "readlink",
							NodeArgs: libparser.NodeArgs{
								&libparser.NodeString{
									Segments: libparser.SegmentedString{
										&libparser.LiteralStringSegment{Contents: "-p"},
									},
								},
								&libparser.NodeString{
									Segments: libparser.SegmentedString{
										&libparser.VariableStringSegment{
											Name:       "MY_LINK",
											Modifiers:  []libparser.StringModifier{},
											IsOptional: false,
										},
									},
								},
							},
						},
						&libparser.NodeString{
							Segments: libparser.SegmentedString{
								&libparser.LiteralStringSegment{Contents: "456"},
							},
						},
					},
				},
			},
		},
	},
	{
		Filename: "05_tomes.tome",
		Expect: &libparser.NodeRoot{
			Tomes: map[string]libparser.Node{
				"first":  nil,
				"second": nil,
			},
			NodeChildren: libparser.NodeChildren{
				&libparser.NodeExec{
					Name: "echo",
					NodeArgs: libparser.NodeArgs{
						&libparser.NodeString{
							Segments: libparser.SegmentedString{
								&libparser.LiteralStringSegment{Contents: "0"},
							},
						},
					},
				},
				&libparser.NodeDirective{
					Name: "tome",
					NodeArgs: libparser.NodeArgs{
						&libparser.NodeString{
							Segments: libparser.SegmentedString{
								&libparser.LiteralStringSegment{Contents: "first"},
							},
						},
					},
					NodeChildren: libparser.NodeChildren{
						&libparser.NodeExec{
							Name: "echo",
							NodeArgs: libparser.NodeArgs{
								&libparser.NodeString{
									Segments: libparser.SegmentedString{
										&libparser.LiteralStringSegment{Contents: "1.1"},
									},
								},
							},
						},
					},
				},
				&libparser.NodeDirective{
					Name: "tome",
					NodeArgs: libparser.NodeArgs{
						&libparser.NodeString{
							Segments: libparser.SegmentedString{
								&libparser.LiteralStringSegment{Contents: "second"},
							},
						},
						&libparser.NodeString{
							Segments: libparser.SegmentedString{
								&libparser.LiteralStringSegment{Contents: "With a description"},
							},
						},
					},
					NodeChildren: libparser.NodeChildren{
						&libparser.NodeExec{
							Name: "echo",
							NodeArgs: libparser.NodeArgs{
								&libparser.NodeString{
									Segments: libparser.SegmentedString{
										&libparser.LiteralStringSegment{Contents: "2.1"},
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
		Filename: "06_semicolon.tome",
		Expect: &libparser.NodeRoot{
			Tomes: map[string]libparser.Node{},
			NodeChildren: libparser.NodeChildren{
				&libparser.NodeExec{
					Name: "echo",
					NodeArgs: libparser.NodeArgs{
						&libparser.NodeString{
							Segments: libparser.SegmentedString{
								&libparser.LiteralStringSegment{Contents: "1"},
							},
						},
					},
				},
				&libparser.NodeExec{
					Name: "echo",
					NodeArgs: libparser.NodeArgs{
						&libparser.NodeString{
							Segments: libparser.SegmentedString{
								&libparser.LiteralStringSegment{Contents: "2"},
							},
						},
					},
				},
				&libparser.NodeExec{
					Name: "echo",
					NodeArgs: libparser.NodeArgs{
						&libparser.NodeString{
							Segments: libparser.SegmentedString{
								&libparser.LiteralStringSegment{Contents: "3"},
							},
						},
					},
				},
				&libparser.NodeExec{
					Name: "echo",
					NodeArgs: libparser.NodeArgs{
						&libparser.NodeString{
							Segments: libparser.SegmentedString{
								&libparser.LiteralStringSegment{Contents: "4"},
							},
						},
					},
				},
			},
		},
	},
	{
		Filename: "07_pipes.tome",
		Expect: &libparser.NodeRoot{
			Tomes: map[string]libparser.Node{},
			NodeChildren: libparser.NodeChildren{
				&libparser.NodePipe{
					Source: &libparser.NodeExec{
						Name: "echo",
						NodeArgs: libparser.NodeArgs{
							&libparser.NodeString{
								Segments: libparser.SegmentedString{
									&libparser.LiteralStringSegment{Contents: "-e"},
								},
							},
							&libparser.NodeString{
								Segments: libparser.SegmentedString{
									&libparser.LiteralStringSegment{Contents: `Hello World!\n`},
								},
							},
						},
					},
					Dest: &libparser.NodeExec{
						Name: "bat",
						NodeArgs: libparser.NodeArgs{
							&libparser.NodeString{
								Segments: libparser.SegmentedString{
									&libparser.LiteralStringSegment{Contents: "--lang"},
								},
							},
							&libparser.NodeString{
								Segments: libparser.SegmentedString{
									&libparser.LiteralStringSegment{Contents: "html"},
								},
							},
						},
					},
				},
				&libparser.NodePipe{
					Source: &libparser.NodeExec{
						Name: "echo",
						NodeArgs: libparser.NodeArgs{
							&libparser.NodeString{
								Segments: libparser.SegmentedString{
									&libparser.LiteralStringSegment{Contents: "123"},
								},
							},
						},
					},
					Dest: &libparser.NodePipe{
						Source: &libparser.NodeExec{
							Name: "program2",
							NodeArgs: libparser.NodeArgs{
								&libparser.NodeString{
									Segments: libparser.SegmentedString{
										&libparser.LiteralStringSegment{Contents: "input"},
									},
								},
							},
						},
						Dest: &libparser.NodePipe{
							Source: &libparser.NodeExec{
								Name: "program3",
								NodeArgs: libparser.NodeArgs{
									&libparser.NodeString{
										Segments: libparser.SegmentedString{
											&libparser.LiteralStringSegment{Contents: "input"},
										},
									},
								},
							},
							Dest: &libparser.NodeExec{
								Name:     "bat",
								NodeArgs: libparser.NodeArgs{},
							},
						},
					},
				},
			},
		},
	},
	{
		Filename: "08_redirects.tome",
		Expect: &libparser.NodeRoot{
			Tomes: map[string]libparser.Node{},
			NodeChildren: libparser.NodeChildren{
				&libparser.NodeRedirect{
					Source: &libparser.NodePipe{
						Source: &libparser.NodeExec{
							Name:     "echo",
							NodeArgs: libparser.NodeArgs{},
						},
						Dest: &libparser.NodeExec{
							Name:     "bat",
							NodeArgs: libparser.NodeArgs{},
						},
					},
					Stdin: &libparser.NodeString{
						Segments: libparser.SegmentedString{
							&libparser.LiteralStringSegment{Contents: "stdin.txt"},
						},
					},
					Stdout: &libparser.NodeString{
						Segments: libparser.SegmentedString{
							&libparser.LiteralStringSegment{Contents: "stdout.txt"},
						},
					},
					Stderr: &libparser.NodeString{
						Segments: libparser.SegmentedString{
							&libparser.LiteralStringSegment{Contents: "stderr.txt"},
						},
					},
				},
			},
		},
	},
}
