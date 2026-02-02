@echo off
REM WeKnora Docker 镜像重构和重启脚本 (Windows)
REM 用法: scripts\rebuild_and_restart.bat [选项]
REM 选项:
REM   dev         使用开发环境配置
REM   prod        使用生产环境配置 [默认]
REM   app         只重构 app 服务
REM   frontend    只重构 frontend 服务
REM   docreader   只重构 docreader 服务
REM   all         重构所有服务 [默认]
REM   no-cache    不使用缓存构建

setlocal enabledelayedexpansion

REM 默认配置
set COMPOSE_FILE=docker-compose.yml
set SERVICE=all
set NO_CACHE=
set PROFILE=
set ENV_MODE=生产

REM 解析命令行参数
:parse_args
if "%~1"=="" goto start_process
if /i "%~1"=="dev" (
    set COMPOSE_FILE=docker-compose.dev.yml
    set ENV_MODE=开发
    shift
    goto parse_args
)
if /i "%~1"=="prod" (
    set COMPOSE_FILE=docker-compose.yml
    set ENV_MODE=生产
    shift
    goto parse_args
)
if /i "%~1"=="app" (
    set SERVICE=app
    shift
    goto parse_args
)
if /i "%~1"=="frontend" (
    set SERVICE=frontend
    shift
    goto parse_args
)
if /i "%~1"=="docreader" (
    set SERVICE=docreader
    shift
    goto parse_args
)
if /i "%~1"=="all" (
    set SERVICE=all
    shift
    goto parse_args
)
if /i "%~1"=="no-cache" (
    set NO_CACHE=--no-cache
    shift
    goto parse_args
)
if /i "%~1"=="minio" (
    set PROFILE=--profile minio
    shift
    goto parse_args
)
if /i "%~1"=="neo4j" (
    set PROFILE=--profile neo4j
    shift
    goto parse_args
)
if /i "%~1"=="qdrant" (
    set PROFILE=--profile qdrant
    shift
    goto parse_args
)
if /i "%~1"=="jaeger" (
    set PROFILE=--profile jaeger
    shift
    goto parse_args
)
if /i "%~1"=="full" (
    set PROFILE=--profile full
    shift
    goto parse_args
)
echo 未知选项: %~1
echo 用法: %~nx0 [dev^|prod] [app^|frontend^|docreader^|all] [no-cache] [minio^|neo4j^|qdrant^|jaeger^|full]
exit /b 1

:start_process
echo ==========================================
echo WeKnora Docker 镜像重构和重启
echo ==========================================
echo 环境模式: %ENV_MODE%
echo 配置文件: %COMPOSE_FILE%
echo 服务范围: %SERVICE%
echo ==========================================

REM 停止并删除容器
echo.
echo 步骤 1/4: 停止现有容器...
if "%SERVICE%"=="all" (
    docker compose -f %COMPOSE_FILE% %PROFILE% down
) else (
    docker compose -f %COMPOSE_FILE% %PROFILE% stop %SERVICE%
    docker compose -f %COMPOSE_FILE% %PROFILE% rm -f %SERVICE%
)

REM 删除旧镜像
echo.
echo 步骤 2/4: 清理旧镜像...
if "%SERVICE%"=="all" (
    echo 清理所有 WeKnora 相关镜像...
    for /f "tokens=3" %%i in ('docker images ^| findstr "wechatopenai/weknora"') do (
        docker rmi -f %%i 2>nul
    )
) else (
    echo 清理 %SERVICE% 镜像...
    for /f "tokens=3" %%i in ('docker images ^| findstr "wechatopenai/weknora-%SERVICE%"') do (
        docker rmi -f %%i 2>nul
    )
)

REM 重新构建镜像
echo.
echo 步骤 3/4: 重新构建镜像...
if "%SERVICE%"=="all" (
    docker compose -f %COMPOSE_FILE% %PROFILE% build %NO_CACHE%
) else (
    docker compose -f %COMPOSE_FILE% %PROFILE% build %NO_CACHE% %SERVICE%
)

REM 启动服务
echo.
echo 步骤 4/4: 启动服务...
if "%SERVICE%"=="all" (
    docker compose -f %COMPOSE_FILE% %PROFILE% up -d
) else (
    docker compose -f %COMPOSE_FILE% %PROFILE% up -d %SERVICE%
)

REM 显示运行状态
echo.
echo ==========================================
echo 服务状态:
echo ==========================================
docker compose -f %COMPOSE_FILE% %PROFILE% ps

echo.
echo ==========================================
echo 重构和重启完成！
echo ==========================================
echo.
echo 查看日志: docker compose -f %COMPOSE_FILE% logs -f [服务名]
echo 停止服务: docker compose -f %COMPOSE_FILE% down
echo.

endlocal
