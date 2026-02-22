@echo off
REM SaveAny-Bot Docker Build Script for Windows
REM 使用方法: build.bat [tag]

setlocal

set TAG=%1
if "%TAG%"=="" set TAG=latest

echo ========================================
echo Building SaveAny-Bot Docker Image
echo Tag: %TAG%
echo ========================================

echo.
echo [1/3] Building Docker image...
docker build -t saveany-bot:%TAG% .

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo [ERROR] Build failed!
    pause
    exit /b 1
)

echo.
echo [2/3] Tagging as latest...
docker tag saveany-bot:%TAG% saveany-bot:latest

echo.
echo [3/3] Build complete!
echo.
echo Usage:
echo   docker run -d ^
echo     -v ./config.toml:/app/config.toml ^
echo     -v ./data:/app/data ^
echo     -v ./downloads:/app/downloads ^
echo     -p 8080:8080 ^
echo     saveany-bot:latest
echo.
echo Or use docker-compose:
echo   docker-compose -f docker-compose.dev.yml up -d
echo.
echo Access Web UI at: http://localhost:8080
echo.

pause
