# Go Module update command

A helper command to update go modules to latest versions.

## Why not just `go get -u`

Running `go get -u` works well unless you have dependencies that are also overridden with `replace`.
Specifically, this will fail to update any dependency if you are using private monorepo with internal libraries, where
modules cannot be looked up - Go Get will get 404 Not Found and entire command will fail.

This tool simply ignores such dependencies and runs `go get -u <mod>` on each **direct** dependency.

## A word of warning

While this tool does not do anything weird - blindly updating dependencies can cause unexpected behavior. Always test
your code properly after updates.
