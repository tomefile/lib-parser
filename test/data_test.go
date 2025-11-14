package libparser_test

import libparser "github.com/tomefile/lib-parser"

type DataTestCase struct {
	Filename string
	Expect   *libparser.NodeTree
}

var ExpectedData = []DataTestCase{
	{
		Filename: "01_basic.tome",
		Expect: &libparser.NodeTree{
			Tomes: map[string]libparser.Node{},
			NodeChildren: libparser.NodeChildren{
				&libparser.CommentNode{Contents: " Example program, привет мир 👨‍🚀!"},
				&libparser.DirectiveNode{
					Name: "include",
					NodeArgs: libparser.NodeArgs{
						&libparser.StringNode{Contents: "@std"},
					},
					NodeChildren: libparser.NodeChildren{},
				},
				&libparser.ExecNode{
					Binary: "echo",
					NodeArgs: libparser.NodeArgs{
						&libparser.StringNode{Contents: "Hello World!"},
						&libparser.StringNode{Contents: "and another line"},
						&libparser.StringNode{Contents: "and another."},
					},
				},
			},
		},
	},
	{
		Filename: "02_directive_body.tome",
		Expect: &libparser.NodeTree{
			Tomes: map[string]libparser.Node{},
			NodeChildren: libparser.NodeChildren{
				&libparser.ExecNode{
					Binary: "echo",
					NodeArgs: libparser.NodeArgs{
						&libparser.StringNode{Contents: "1"},
					},
				},
				&libparser.DirectiveNode{
					Name: "section",
					NodeArgs: libparser.NodeArgs{
						&libparser.StringNode{Contents: "Hello World!"},
					},
					NodeChildren: libparser.NodeChildren{
						&libparser.ExecNode{
							Binary: "echo",
							NodeArgs: libparser.NodeArgs{
								&libparser.LiteralNode{Contents: "1.1"},
							},
						},
						&libparser.ExecNode{
							Binary: "echo",
							NodeArgs: libparser.NodeArgs{
								&libparser.StringNode{Contents: "1.2"},
							},
						},
					},
				},
			},
		},
	},
	{
		Filename: "03_directive_nested.tome",
		Expect: &libparser.NodeTree{
			Tomes: map[string]libparser.Node{},
			NodeChildren: libparser.NodeChildren{
				&libparser.ExecNode{
					Binary: "echo",
					NodeArgs: libparser.NodeArgs{
						&libparser.StringNode{Contents: "1"},
					},
				},
				&libparser.DirectiveNode{
					Name: "section",
					NodeArgs: libparser.NodeArgs{
						&libparser.StringNode{Contents: "Hello World!"},
					},
					NodeChildren: libparser.NodeChildren{
						&libparser.ExecNode{
							Binary: "echo",
							NodeArgs: libparser.NodeArgs{
								&libparser.StringNode{Contents: "1.1"},
							},
						},
						&libparser.ExecNode{
							Binary: "echo",
							NodeArgs: libparser.NodeArgs{
								&libparser.StringNode{Contents: "1.2"},
							},
						},
						&libparser.DirectiveNode{
							Name: "section",
							NodeArgs: libparser.NodeArgs{
								&libparser.StringNode{Contents: "Nested"},
							},
							NodeChildren: libparser.NodeChildren{
								&libparser.CommentNode{Contents: " This is nested inside"},
								&libparser.ExecNode{
									Binary: "echo",
									NodeArgs: libparser.NodeArgs{
										&libparser.StringNode{Contents: "2.1"},
									},
								},
								&libparser.ExecNode{
									Binary: "echo",
									NodeArgs: libparser.NodeArgs{
										&libparser.StringNode{Contents: "2.2"},
									},
								},
							},
						},
						&libparser.ExecNode{
							Binary: "echo",
							NodeArgs: libparser.NodeArgs{
								&libparser.StringNode{Contents: "1.3"},
							},
						},
					},
				},
			},
		},
	},
	{
		Filename: "04_subcommand.tome",
		Expect: &libparser.NodeTree{
			Tomes: map[string]libparser.Node{},
			NodeChildren: libparser.NodeChildren{
				&libparser.CallNode{
					Macro: "my_macro",
					NodeArgs: libparser.NodeArgs{
						&libparser.StringNode{Contents: "123"},
						&libparser.ExecNode{
							Binary: "readlink",
							NodeArgs: libparser.NodeArgs{
								&libparser.StringNode{Contents: "-p"},
								&libparser.StringNode{Contents: "$MY_LINK"},
							},
						},
						&libparser.StringNode{Contents: "456"},
					},
				},
			},
		},
	},
	{
		Filename: "05_tomes.tome",
		Expect: &libparser.NodeTree{
			Tomes: map[string]libparser.Node{
				"first":  nil,
				"second": nil,
			},
			NodeChildren: libparser.NodeChildren{
				&libparser.ExecNode{
					Binary: "echo",
					NodeArgs: libparser.NodeArgs{
						&libparser.StringNode{Contents: "0"},
					},
				},
				&libparser.DirectiveNode{
					Name: "tome",
					NodeArgs: libparser.NodeArgs{
						&libparser.StringNode{Contents: "first"},
					},
					NodeChildren: libparser.NodeChildren{
						&libparser.ExecNode{
							Binary: "echo",
							NodeArgs: libparser.NodeArgs{
								&libparser.StringNode{Contents: "1.1"},
							},
						},
					},
				},
				&libparser.DirectiveNode{
					Name: "tome",
					NodeArgs: libparser.NodeArgs{
						&libparser.StringNode{Contents: "second"},
						&libparser.StringNode{Contents: "With a description"},
					},
					NodeChildren: libparser.NodeChildren{
						&libparser.ExecNode{
							Binary: "echo",
							NodeArgs: libparser.NodeArgs{
								&libparser.StringNode{Contents: "2.1"},
							},
						},
					},
				},
			},
		},
	},
	{
		Filename: "06_semicolon.tome",
		Expect: &libparser.NodeTree{
			Tomes: map[string]libparser.Node{},
			NodeChildren: libparser.NodeChildren{
				&libparser.ExecNode{
					Binary: "echo",
					NodeArgs: libparser.NodeArgs{
						&libparser.StringNode{Contents: "1"},
					},
				},
				&libparser.ExecNode{
					Binary: "echo",
					NodeArgs: libparser.NodeArgs{
						&libparser.StringNode{Contents: "2"},
					},
				},
				&libparser.ExecNode{
					Binary: "echo",
					NodeArgs: libparser.NodeArgs{
						&libparser.StringNode{Contents: "3"},
					},
				},
				&libparser.ExecNode{
					Binary: "echo",
					NodeArgs: libparser.NodeArgs{
						&libparser.StringNode{Contents: "4"},
					},
				},
			},
		},
	},
	{
		Filename: "07_pipes.tome",
		Expect: &libparser.NodeTree{
			Tomes: map[string]libparser.Node{},
			NodeChildren: libparser.NodeChildren{
				&libparser.PipeNode{
					Source: &libparser.ExecNode{
						Binary: "echo",
						NodeArgs: libparser.NodeArgs{
							&libparser.StringNode{Contents: "-e"},
							&libparser.StringNode{Contents: `Hello World!\n`},
						},
					},
					Dest: &libparser.ExecNode{
						Binary: "bat",
						NodeArgs: libparser.NodeArgs{
							&libparser.StringNode{Contents: "--lang"},
							&libparser.StringNode{Contents: "html"},
						},
					},
				},
				&libparser.PipeNode{
					Source: &libparser.ExecNode{
						Binary: "echo",
						NodeArgs: libparser.NodeArgs{
							&libparser.StringNode{Contents: "123"},
						},
					},
					Dest: &libparser.PipeNode{
						Source: &libparser.ExecNode{
							Binary: "program2",
							NodeArgs: libparser.NodeArgs{
								&libparser.StringNode{Contents: "input"},
							},
						},
						Dest: &libparser.PipeNode{
							Source: &libparser.ExecNode{
								Binary: "program3",
								NodeArgs: libparser.NodeArgs{
									&libparser.StringNode{Contents: "input"},
								},
							},
							Dest: &libparser.ExecNode{
								Binary:   "bat",
								NodeArgs: libparser.NodeArgs{},
							},
						},
					},
				},
			},
		},
	},
}
