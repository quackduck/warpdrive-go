function wd() {
    exitCodeCD=3

    output=$(wd-go $@)
    if [ $? -eq $exitCodeCD ]; then
        cd "$output"
    else
        echo "$output"
    fi
    export wd_last_added_dir=$(pwd)
}

function __wd-go_on_prompt() {
    if [ "$wd_last_added_dir" != $(pwd) -a $(pwd) != "$HOME" ]; then
      wd-go --add $(pwd)
    fi
    export wd_last_added_dir=$(pwd)
}

# zsh: attach a hook to precmd so we don't clobber what's already in precmd
if [ ${ZSH_VERSION} ]; then
  autoload -Uz add-zsh-hook
  add-zsh-hook precmd __wd-go_on_prompt
fi

[ ${BASH_VERSION} ] && export PROMPT_COMMAND="__wd-go_on_prompt;$PROMPT_COMMAND"