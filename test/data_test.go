package libparser_test

import libparser "github.com/tomefile/lib-parser"

type DataTestCase struct {
	Filename string
	Expect   *libparser.NodeRoot
}

var ExpectedData = []DataTestCase{
	{
		Filename: "01_syntax.tome",
		Expect: &libparser.NodeRoot{
			Tomes: map[string]libparser.Node{},
			NodeChildren: libparser.NodeChildren{
				&libparser.NodeComment{Contents: " Example program, привет мир 👨‍🚀!"},
				&libparser.NodeDirective{
					Name: "include",
					NodeArgs: libparser.NodeArgs{
						libparser.NewSimpleNodeString("@std"),
					},
					NodeChildren: libparser.NodeChildren{},
				},
				&libparser.NodeExec{
					Name: "echo",
					NodeArgs: libparser.NodeArgs{
						libparser.NewSimpleNodeString("Hello World!"),
						libparser.NewSimpleNodeString("and another line"),
						libparser.NewSimpleNodeString("and another."),
					},
				},
				&libparser.NodeDirective{
					Name: "section",
					NodeArgs: libparser.NodeArgs{
						&libparser.NodeExec{
							Name: "echo",
							NodeArgs: libparser.NodeArgs{
								libparser.NewSimpleNodeString("/tmp/filename.png"),
							},
						},
						&libparser.NodeLiteral{
							Contents: "Some literal string",
						},
					},
					NodeChildren: libparser.NodeChildren{
						&libparser.NodeDirective{
							Name: "assert",
							NodeArgs: libparser.NodeArgs{
								&libparser.NodeString{
									Segments: libparser.SegmentedString{
										&libparser.VariableStringSegment{
											Name: "build_dir",
											Modifiers: []libparser.StringModifier{
												getModifierSafe(libparser.MOD_IS_DIR),
												getModifierSafe(libparser.MOD_NOT),
											},
											IsOptional: true,
										},
									},
								},
							},
							NodeChildren: libparser.NodeChildren{},
						},
						&libparser.NodeDirective{
							Name: "for",
							NodeArgs: libparser.NodeArgs{
								&libparser.NodeString{
									Segments: libparser.SegmentedString{
										&libparser.VariableStringSegment{
											Name:       "file",
											Modifiers:  []libparser.StringModifier{},
											IsOptional: false,
										},
									},
								},
								libparser.NewSimpleNodeString("="),
								&libparser.NodeString{
									Segments: libparser.SegmentedString{
										&libparser.VariableStringSegment{
											Name:       "in_dir",
											Modifiers:  []libparser.StringModifier{},
											IsOptional: false,
										},
										&libparser.LiteralStringSegment{Contents: "/"},
										&libparser.VariableStringSegment{
											Name:       "pattern",
											Modifiers:  []libparser.StringModifier{},
											IsOptional: false,
										},
										&libparser.LiteralStringSegment{Contents: ".json.patch"},
									},
								},
							},
							NodeChildren: libparser.NodeChildren{
								&libparser.NodeRedirect{
									Source: &libparser.NodeExec{
										Name: "patch",
										NodeArgs: libparser.NodeArgs{
											libparser.NewSimpleNodeString("-s"),
											libparser.NewSimpleNodeString("-o"),
											libparser.NewSimpleNodeString("/tmp/patched-file.json"),
											&libparser.NodeExec{
												Name: "realpath",
												NodeArgs: libparser.NodeArgs{
													&libparser.NodeString{
														Segments: libparser.SegmentedString{
															&libparser.LiteralStringSegment{
																Contents: "../something/something/",
															},
															&libparser.VariableStringSegment{
																Name:       "basename",
																Modifiers:  []libparser.StringModifier{},
																IsOptional: false,
															},
														},
													},
												},
											},
											&libparser.NodeString{
												Segments: libparser.SegmentedString{
													&libparser.VariableStringSegment{
														Name:       "file",
														Modifiers:  []libparser.StringModifier{},
														IsOptional: false,
													},
												},
											},
										},
									},
									Stdout: libparser.NewSimpleNodeString("/some/output"),
								},
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
						libparser.NewSimpleNodeString("1"),
					},
				},
				&libparser.NodeDirective{
					Name: "section",
					NodeArgs: libparser.NodeArgs{
						libparser.NewSimpleNodeString("Hello World!"),
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
								libparser.NewSimpleNodeString("1.2"),
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
						libparser.NewSimpleNodeString("1"),
					},
				},
				&libparser.NodeDirective{
					Name: "section",
					NodeArgs: libparser.NodeArgs{
						libparser.NewSimpleNodeString("Hello World!"),
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
								libparser.NewSimpleNodeString("1.2"),
							},
						},
						&libparser.NodeDirective{
							Name: "section",
							NodeArgs: libparser.NodeArgs{
								libparser.NewSimpleNodeString("Nested"),
							},
							NodeChildren: libparser.NodeChildren{
								&libparser.NodeComment{Contents: " This is nested inside"},
								&libparser.NodeExec{
									Name: "echo",
									NodeArgs: libparser.NodeArgs{
										libparser.NewSimpleNodeString("2.1"),
									},
								},
								&libparser.NodeExec{
									Name: "echo",
									NodeArgs: libparser.NodeArgs{
										libparser.NewSimpleNodeString("2.2"),
									},
								},
							},
						},
						&libparser.NodeExec{
							Name: "echo",
							NodeArgs: libparser.NodeArgs{
								libparser.NewSimpleNodeString("1.3"),
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
						libparser.NewSimpleNodeString("123"),
						&libparser.NodeExec{
							Name: "readlink",
							NodeArgs: libparser.NodeArgs{
								libparser.NewSimpleNodeString("-p"),
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
						libparser.NewSimpleNodeString("456"),
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
						libparser.NewSimpleNodeString("0"),
					},
				},
				&libparser.NodeDirective{
					Name: "tome",
					NodeArgs: libparser.NodeArgs{
						libparser.NewSimpleNodeString("first"),
					},
					NodeChildren: libparser.NodeChildren{
						&libparser.NodeExec{
							Name: "echo",
							NodeArgs: libparser.NodeArgs{
								libparser.NewSimpleNodeString("1.1"),
							},
						},
					},
				},
				&libparser.NodeDirective{
					Name: "tome",
					NodeArgs: libparser.NodeArgs{
						libparser.NewSimpleNodeString("second"),
						libparser.NewSimpleNodeString("With a description"),
					},
					NodeChildren: libparser.NodeChildren{
						&libparser.NodeExec{
							Name: "echo",
							NodeArgs: libparser.NodeArgs{
								libparser.NewSimpleNodeString("2.1"),
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
						libparser.NewSimpleNodeString("1"),
					},
				},
				&libparser.NodeExec{
					Name: "echo",
					NodeArgs: libparser.NodeArgs{
						libparser.NewSimpleNodeString("2"),
					},
				},
				&libparser.NodeExec{
					Name: "echo",
					NodeArgs: libparser.NodeArgs{
						libparser.NewSimpleNodeString("3"),
					},
				},
				&libparser.NodeExec{
					Name: "echo",
					NodeArgs: libparser.NodeArgs{
						libparser.NewSimpleNodeString("4"),
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
							libparser.NewSimpleNodeString("-e"),
							libparser.NewSimpleNodeString(`Hello World!\n`),
						},
					},
					Dest: &libparser.NodeExec{
						Name: "bat",
						NodeArgs: libparser.NodeArgs{
							libparser.NewSimpleNodeString("--lang"),
							libparser.NewSimpleNodeString("html"),
						},
					},
				},
				&libparser.NodePipe{
					Source: &libparser.NodeExec{
						Name: "echo",
						NodeArgs: libparser.NodeArgs{
							libparser.NewSimpleNodeString("123"),
						},
					},
					Dest: &libparser.NodePipe{
						Source: &libparser.NodeExec{
							Name: "program2",
							NodeArgs: libparser.NodeArgs{
								libparser.NewSimpleNodeString("input"),
							},
						},
						Dest: &libparser.NodePipe{
							Source: &libparser.NodeExec{
								Name: "program3",
								NodeArgs: libparser.NodeArgs{
									libparser.NewSimpleNodeString("input"),
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
					Stdin:  libparser.NewSimpleNodeString("stdin.txt"),
					Stdout: libparser.NewSimpleNodeString("stdout.txt"),
					Stderr: libparser.NewSimpleNodeString("stderr.txt"),
				},
			},
		},
	},
}

func getModifierSafe(name libparser.ModifierName) libparser.StringModifier {
	modifier, _ := libparser.GetModifier(name, []*libparser.NodeString{})
	return modifier
}
