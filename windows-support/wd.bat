@echo off
wd-go %* > temp
set /p output=<temp
if %ERRORLEVEL%==3 (
    del temp
    cd %output%
) else (
    @echo on
    type temp
    @echo off
    del temp
)
set wd_previous_dir=%cd%


