# ppath: A Precious Path Linter

[precious](https://github.com/houseabsolute/precious) allows you to include and exclude files on a per-command basis. Once a configuration becomes complex, it's possible to move/rename a file so that it no longer appears in either the include or exclude list. You may not realize this until much later, if at all, and you may have a file which is now no longer covered by your code quality tools.

What this utility aims to do is to protect you from this scenario. `ppath` will parse your `precious` config and return a non-zero exit code if a file path or pattern in the includes or excludes cannot be found.
After installing the published binary, you can enable it in your `precious.toml` in the following way:

```toml
[commands.ppath]
type = "lint"
include = ["precious.toml"]
run_mode = "files"
cmd = ["ppath"]
ok_exit_codes = 0
```

This will lint your config when your `precious.toml` changes, passing the path to your config file to `ppath`, which is the only argument which it expects. You could also run this in a pre-commit hook via:

`ppath precious.toml`

or even

`precious lint --command ppath precious.toml`

If you want to ignore certain paths or patterns, you can ignore them on a global or per-command basis. This configuration needs to exist in a file called `.ppath.toml`, which lives in the top level of your project workspace -- the same directory from which `precious` runs by default.

A `.ppath.toml` could look something like this:

```toml
ignore = [
    "**/node_modules/**/*",
]

[commands.omegasort-gitignore]
ignore = [
    "**/foo/**/*",
    "bar.txt",
]
```

The above `ppath` configuration will ignore the pattern `**/node_modules/**/*` if it comes up in the include or exclude list of any command in the `precious` configuration. It will also ignore the pattern `"**/foo/**/*"` and the file `bar.txt` if they appear in either the include or exclude list of the `omegasort-gitignore` command in your `precious` configuration.
