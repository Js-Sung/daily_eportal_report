@title -= daily report of health =-

@cd /d %~dp0

:: modify the following fields 
@set id=201821019876
@set passwd=1234567890
@set chrome_exe=D:\Program Files\CentBrowser\Application\chrome.exe


@set tempdir=%~dp0temp

@if exist %tempdir% (
   echo folder %tempdir% exits
) else (
   echo making folder %tempdir%
   md %tempdir%
)


@go run func.go -id="%id%" -passwd="%passwd%" -exe_path="%chrome_exe%" -temp_path="%tempdir%"


@cd /d %tempdir%

@timeout /T 2 /NOBREAK
@for /F %%D in ('dir /ad/b chromedp-*') do @rd /s /q %%D

::exit
