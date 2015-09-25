@echo off
cls

::———————————————————
IF "%~1"=="–FIX_CTRL_C" (
SHIFT
) ELSE (
CALL <NUL %0 –FIX_CTRL_C %*
GOTO EOF
)
::———————————————————

gin -p 12201 -a 12301 -i run "./config/server.gcfg"
pause
