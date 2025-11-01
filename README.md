# Tomefile Parser

Library to parse [Tomefile](https://github.com/tomefile/specification) code and output statements similar to [AST](https://en.wikipedia.org/wiki/Abstract_syntax_tree) for use by another program.

Uses [just-in-time compilation](https://en.wikipedia.org/wiki/Just-in-time_compilation) to allow for execution of statements as soon as possible rather than waiting for the entire code to parse first.
**It comes with a drawback** of not catching errors ahead of time.
