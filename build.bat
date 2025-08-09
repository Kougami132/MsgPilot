@echo off
setlocal enabledelayedexpansion
echo ========================================
echo        MsgPilot Build Script
echo ========================================
echo.

:: Build configuration
set "BUILD_DIR=dist"
set "FRONTEND_DIR=frontend"
set "DIST_DIR=%FRONTEND_DIR%\dist"

:: Check required tools
echo Checking build environment...

where node >nul 2>&1
if errorlevel 1 (
    echo ERROR: Node.js not found
    pause
    exit /b 1
)

where npm >nul 2>&1
if errorlevel 1 (
    echo ERROR: npm not found
    pause
    exit /b 1
)

where go >nul 2>&1
if errorlevel 1 (
    echo ERROR: Go not found
    pause
    exit /b 1
)

echo Build environment check completed
echo.

:: Create build directory
if not exist "%BUILD_DIR%" (
    mkdir "%BUILD_DIR%"
    echo Created build directory: %BUILD_DIR%
)

:: Clean previous build files
echo Cleaning previous build files...
if exist "%BUILD_DIR%\*" (
    del /q "%BUILD_DIR%\*" 2>nul
)
echo Cleanup completed
echo.

:: Build frontend
echo ========================================
echo        Building Frontend Assets
echo ========================================

cd "%FRONTEND_DIR%"
if errorlevel 1 (
    echo ERROR: Cannot enter frontend directory
    pause
    exit /b 1
)

echo Installing frontend dependencies...
call npm install
if errorlevel 1 (
    echo ERROR: Frontend dependency installation failed
    cd ..
    pause
    exit /b 1
)

echo Building frontend assets...
call npm run build
if errorlevel 1 (
    echo ERROR: Frontend build failed
    cd ..
    pause
    exit /b 1
)

cd ..
echo Frontend build completed
echo.

:: Copy frontend static files to build directory
if exist "%DIST_DIR%" (
    echo Copying frontend static files...
    xcopy "%DIST_DIR%" "%BUILD_DIR%\static" /E /I /Y
    if errorlevel 1 (
        echo ERROR: Failed to copy frontend static files
        pause
        exit /b 1
    )
    echo Frontend static files copied successfully
) else (
    echo ERROR: Frontend build output does not exist
    pause
    exit /b 1
)
echo.

:: Build backend
echo ========================================
echo        Building Backend Executables
echo ========================================

echo Downloading Go dependencies...
go mod download
if errorlevel 1 (
    echo ERROR: Go dependency download failed
    pause
    exit /b 1
)

echo Building Windows version...
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w" -o "%BUILD_DIR%\msgpilot-windows-amd64.exe" .
if errorlevel 1 (
    echo ERROR: Windows version build failed
    pause
    exit /b 1
)
echo Windows version build completed: %BUILD_DIR%\msgpilot-windows-amd64.exe

echo Building Linux version...
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-s -w" -o "%BUILD_DIR%\msgpilot-linux-amd64" .
if errorlevel 1 (
    echo ERROR: Linux version build failed
    pause
    exit /b 1
)
echo Linux version build completed: %BUILD_DIR%\msgpilot-linux-amd64

:: Reset environment variables
set GOOS=
set GOARCH=

echo.
echo ========================================
echo        Build Completed
echo ========================================
echo.
echo Build artifacts location:
echo   - Frontend static files: %BUILD_DIR%\static\
echo   - Windows executable: %BUILD_DIR%\msgpilot-windows-amd64.exe
echo   - Linux executable: %BUILD_DIR%\msgpilot-linux-amd64
echo.

:: Display build artifact sizes
echo Build artifact information:
if exist "%BUILD_DIR%\msgpilot-windows-amd64.exe" (
    for %%A in ("%BUILD_DIR%\msgpilot-windows-amd64.exe") do (
        set "size=%%~zA"
        if defined size (
            set /a sizeMB=!size!/1024/1024
            echo   - Windows version size: !sizeMB! MB
        ) else (
            echo   - Windows version: Size unavailable
        )
    )
)

if exist "%BUILD_DIR%\msgpilot-linux-amd64" (
    for %%A in ("%BUILD_DIR%\msgpilot-linux-amd64") do (
        set "size=%%~zA"
        if defined size (
            set /a sizeMB=!size!/1024/1024
            echo   - Linux version size: !sizeMB! MB
        ) else (
            echo   - Linux version: Size unavailable
        )
    )
)

echo.
echo All build tasks completed successfully!
pause
