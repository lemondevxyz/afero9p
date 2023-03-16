# afero9p
Export your favorite afero filesystems using the magic of 9p?

## why?
why not? afero is an amazing abstraction. 9p is an amazing protocol. mix em together and the possibilities become endless.

## possible usage(s)
tons of usages but not for the reasons you might think. one that jumps to my mind is the ability to use any Go library across programming languages.

Some cool languages, like Common Lisp, have very shotty library support. The solution? Make more libraries. Not by actually creating more libraries when there are already written well-designed libraries but by exporting libraries in other languages.

How? Make a go abstraction that makes developing filesystems easy(afero). Export filesystems through a battle tested simple protocol (9p).

You might be saying: "What if my language doesn't support 9p?" - I'd tell you either support it or mount that 9p filesystem through FUSE. Because FUSE is too linux dependent and 9p is almost universal.