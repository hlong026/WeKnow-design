@echo off
REM WeKnora Docker 快速重启脚本 (Windows)
REM 用法: scripts\restart.bat [dev|prod] [服务名]

setlocal

set COMPOSE_FILE=docker-compose.yml
set SERVICE=

if /i "%~1"=="dev" (
    set COMPOSE_FILE=docker-compose.dev.yml
    set SERVICE=%~2
) else if /i "%~1"=="prod" (
    set COMPOSE_FILE=docker-compose.yml
    set SERVICE=%~2
) else (
    set SERVICE=%~1
)

echo ==========================================
echo WeKnora Docker 快速重启
echo ==========================================
echo 配置文件: %COMPOSE_FILE%
if not "%SERVICE%"=="" (
    echo 服务: %SERVICE%
) else (
    echo 服务: 所有服务
)
echo ==========================================

echo.
echo 重启服务中...
docker compose -f %COMPOSE_FILE% restart %SERVICE%

echo.
echo ==========================================
echo 服务状态:
echo ==========================================
docker compose -f %COMPOSE_FILE% ps

echo.
echo 重启完成！
echo.

endlocal
