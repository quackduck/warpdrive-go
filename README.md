# WarpDrive: the Go version.

![image](https://user-images.githubusercontent.com/38882631/160743437-f5ae7827-bbfb-4dcb-bc01-ecef6a814162.png)


## What does this do?

Instead of having a huge cd routine to get where you want, with WarpDrive you use short keywords to warp to the directory in ~5 ms.

## Help

```yaml
WarpDrive - Warp across the filesystem instantly

Usage: wd [<pattern> ...]
       wd --list/-l | --help/-h | --version/-v
       wd {--add/-a | --remove/-r} <path>
Options:
   --list/-l     List currently tracked paths along with their frecency scores
   --add/-a      Add a path to the data file (paths will be added automatically)
   --remove/-r   Remove a path from the data file
   --help/-h     Print this help message
   --version/-v  Print the version of WarpDrive installed
Examples:
   wd -l                          # list all tracked paths and scores
   wd                             # cd to home directory
   wd s                           # tries to match 's'
   wd someDir                     # tries to match 'someDir'
   wd some subDir                 # ensures matched path also contains 'some'
   wd /absolute/path/to/someDir   # absolute paths work too
Note:
   When specifying multiple patterns, order does not matter, except for the last pattern.
   WarpDrive will always cd to a directory that matches the last pattern.
```

Related: https://github.com/quackduck/WarpDrive. That repo has the earlier, slower, Java version of this.

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

### If you're a windows command prompt user:
1. Put `windows-support\wd-go_on_prompt.bat` and `windows-support\wd.bat` in a directory of your choice.
2. Compile the go module into a .exe file and put that in the same directory.
3. Add the directory to the `PATH` environment variable
4. Make a script that runs when command prompt is opened: https://superuser.com/a/916478.
5. In the script add the line `doskey cd=wd-go_on_prompt $*`
6. Profit!

### If you use some other shell:
Make a PR to add support for your shell. Adding support for other shells is easier than for similar programs. All you need is:
1. A way for `wd` to add a new directory when you use `cd`
         
   You can do this with a prompt hook that detects a change in the pwd.
2. A way for your shell to `cd` to `wd`'s result, or show its output as needed.

    The `wd-go` binary returns an exit code of `3` if it wants you to `cd` to its output. Else, please show the user the output.

See the `bash-zsh-support/` and `fish-support/` directories for examples of how to do this.
