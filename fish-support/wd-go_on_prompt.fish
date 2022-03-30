function __wd-go_on_prompt --on-event fish_prompt
    if test ! "$wd_last_added_dir"
        set -g wd_last_added_dir (pwd)
    end
    if test "$wd_last_added_dir" != (pwd) -a (pwd) != "$HOME"
        wd-go --add (pwd)
    end
    set -g wd_last_added_dir (pwd)
end