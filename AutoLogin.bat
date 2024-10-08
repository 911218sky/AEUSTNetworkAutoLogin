@if (@CodeSection == @Batch) @then
@echo off
setlocal

rem Change the current directory to C:\AutoLogin
cd /d "C:\AutoLogin"

set "app=AEUSTNetworkAutoLogin.exe"

rem Create the invisible window
set "vbs=%temp%\invisible.vbs"
echo CreateObject^("WScript.Shell"^).Run """" ^& WScript.Arguments(0) ^& """", 0, False > "%vbs%"

rem Launch the app in the invisible window
cscript.exe "%vbs%" "%app%"

rem Clean up
del "%vbs%"

exit /b
@end
