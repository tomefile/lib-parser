# Tomefile Parser

Library to parse [Tomefile](https://github.com/tomefile) code and output statements similar to [AST](https://en.wikipedia.org/wiki/Abstract_syntax_tree) for use by another program.

Uses [just-in-time compilation](https://en.wikipedia.org/wiki/Just-in-time_compilation) to allow for execution of statements as soon as possible rather than waiting for the entire code to parse first.
**It comes with a drawback** of not catching errors ahead of time.

# Dev

```

```

something something

```
> [directive 'include'] -> [string '<log>']
> [comment '# This is an example']
> [exec 'echo'] -> [string '1']
> [directive 'section'] -> [string 'This is a section']
    > [attach '0'] -> [exec 'echo'] -> [literal '1.1']
    > [attach '0'] -> [exec 'echo'] -> [string '1.2']
    > [attach '0'] -> [directive 'section'] -> [scope '1']
        > [attach '1'] -> [comment '# This is a nested inside']
        > [attach '1'] -> [exec 'echo'] -> [string '2.1']
        > [attach '1'] -> [exec 'echo'] -> [exec 'readlink' [string '-p'] -> [string '$HOME']] -> [literal 'Hello World']
    > [attach '0'] -> [exec 'echo'] -> [string '1.3']
```
