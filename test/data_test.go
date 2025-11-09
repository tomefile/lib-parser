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
						&libparser.LiteralNode{Contents: "<std>"},
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
				&libparser.ExecNode{
					Binary: "echo",
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
}
