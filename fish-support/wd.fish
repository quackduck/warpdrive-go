function wd --description 'Warp across directories'
    set exitCodeCD 3

    set output (wd-go $argv)
    if test $status -eq $exitCodeCD # if exit code is zero, print output
        cd $output; or return
    else
        for line in $output
            echo $line
        end
    end
    set -g wd_last_added_dir (pwd)
end
