# WarpDrive: the Go version.


## What does this do?

Instead of having a huge cd routine to get where you want, with WarpDrive you use short keywords to warp to the directory in ~8 ms.

See github.com/quackduck/WarpDrive for more info, that repo has the earlier, slower, Java version of this.

## Install directions

### If you're a bash or zsh user:
1. Put `bash-zsh-support/wd.sh` somewhere you want.
2. Put `. /path/to/wd.sh` in your relevant rc/profile file (you know what I mean)
3. Run `go install` in this directory
4. Profit!

### If you're a fish user:
1. Put `fish-support/wd-go_on_prompt.fish` in `~/.config/fish/conf.d/`
2. Put `fish-support/wd.fish` in `~/.config/fish/functions/`
3. Run `go install` in this directory
4. Profit!

### If you use some other shell:
Make a PR to add support for your shell. Adding support for other shells is easier than for similar programs. All you need is:
1. A way for `wd` to add a new directory when you use `cd`
         
   You can do this with a prompt hook that detects a change in the pwd.
2. A way for your shell to `cd` to `wd`'s result, or show its output as needed.

    The `wd-go` binary returns an exit code of `3` if it wants you to `cd` to its output. Else, please show the user the output.

See the `bash-zsh-support/` and `fish-support/` directories for examples of how to do this.
