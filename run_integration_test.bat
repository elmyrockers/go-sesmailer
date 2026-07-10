@echo off
REM run_integration_test.bat
REM Runs the sesmailer integration tests (real AWS SES calls).

echo Running sesmailer integration tests...
echo.

go test -tags=integration -v .

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ==============================
    echo  Integration tests FAILED
    echo ==============================
    pause
    exit /b %ERRORLEVEL%
) else (
    echo.
    echo ==============================
    echo  Integration tests PASSED
    echo ==============================
    pause
)
