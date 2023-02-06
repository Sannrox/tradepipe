### CODING STYLE

Unless explicitly stated, we follow all coding guidelines from the Go community. While some of these standards may seem arbitrary, they somehow seem to result in a solid, consistent codebase.

The rules:

1. All code should be formatted with `gofmt -s`.
2. All code should pass the default levels of `golint`.
3. All code should follow the guidelines covered in [Effective Go](https://golang.org/doc/effective_go) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).
4. Comment the code. Tell us the why, the history and the context.
5. Document all declarations and methods, even private ones. Declare expectations, caveats and anything else that may be important. If a type gets exported, having the comments already there will ensure it's ready.
6. Variable name length should be proportional to its context and no longer. `noCommaALongVariableNameLikeThisIsNotMoreClearWhenASimpleCommentWouldDo`. In practice, short methods will have short variable names and globals will have longer names.
7. No underscores in package names. If you need a compound name, step back, and re-examine why you need a compound name. If you still think you need a compound name, lose the underscore.
8. No utils or helpers packages. If a function is not general enough to warrant its own package, it has not been written generally enough to be a part of a util package. Just leave it unexported and well-documented.
9. All tests should run with `go test` (`make test`) and outside tooling should not be required. No, we don't need another unit testing framework. Assertion packages are acceptable if they provide real incremental value.
10. Even though we call these "rules" above, they are actually just guidelines. Since you've read all the rules, you now know that.

If you are having trouble getting into the mood of idiomatic Go, we recommend reading through [Effective Go](https://golang.org/doc/effective_go). The [Go Blog](https://blog.golang.org/) is also a great resource.

