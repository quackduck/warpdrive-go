@echo off
if NOT "%wd_previous_dir%"=="%cd%" if NOT "%cd%"=="%USERPROFILE%" (
    wd-go --add "%cd%"
)
set wd_previous_dir="%cd%"