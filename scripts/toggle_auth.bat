@echo off
REM 切换认证功能的脚本 (Windows版本)
REM 用法: scripts\toggle_auth.bat [enable|disable|status]

setlocal enabledelayedexpansion

set ENV_FILE=.env

if not exist "%ENV_FILE%" (
    echo 错误: .env 文件不存在
    echo 请先复制 .env.example 为 .env
    exit /b 1
)

if "%1"=="disable" (
    echo 禁用用户认证...
    findstr /C:"DISABLE_AUTH=" "%ENV_FILE%" >nul
    if !errorlevel! equ 0 (
        powershell -Command "(Get-Content '%ENV_FILE%') -replace '^DISABLE_AUTH=.*', 'DISABLE_AUTH=true' | Set-Content '%ENV_FILE%'"
    ) else (
        echo DISABLE_AUTH=true >> "%ENV_FILE%"
    )
    
    REM 同时隐藏 Ollama 设置
    findstr /C:"HIDE_OLLAMA=" "%ENV_FILE%" >nul
    if !errorlevel! equ 0 (
        powershell -Command "(Get-Content '%ENV_FILE%') -replace '^HIDE_OLLAMA=.*', 'HIDE_OLLAMA=true' | Set-Content '%ENV_FILE%'"
    ) else (
        echo HIDE_OLLAMA=true >> "%ENV_FILE%"
    )
    
    echo ✓ 已禁用用户认证
    echo ✓ 已隐藏 Ollama 相关设置
    echo 提示: 需要重启服务才能生效
    echo 运行: scripts\start_all.bat --stop ^&^& scripts\start_all.bat
    goto :end
)

if "%1"=="enable" (
    echo 启用用户认证...
    findstr /C:"DISABLE_AUTH=" "%ENV_FILE%" >nul
    if !errorlevel! equ 0 (
        powershell -Command "(Get-Content '%ENV_FILE%') -replace '^DISABLE_AUTH=.*', 'DISABLE_AUTH=false' | Set-Content '%ENV_FILE%'"
    ) else (
        echo DISABLE_AUTH=false >> "%ENV_FILE%"
    )
    
    REM 同时显示 Ollama 设置
    findstr /C:"HIDE_OLLAMA=" "%ENV_FILE%" >nul
    if !errorlevel! equ 0 (
        powershell -Command "(Get-Content '%ENV_FILE%') -replace '^HIDE_OLLAMA=.*', 'HIDE_OLLAMA=false' | Set-Content '%ENV_FILE%'"
    ) else (
        echo HIDE_OLLAMA=false >> "%ENV_FILE%"
    )
    
    echo ✓ 已启用用户认证
    echo ✓ 已显示 Ollama 相关设置
    echo 提示: 需要重启服务才能生效
    echo 运行: scripts\start_all.bat --stop ^&^& scripts\start_all.bat
    goto :end
)

if "%1"=="status" (
    findstr /C:"DISABLE_AUTH=true" "%ENV_FILE%" >nul
    if !errorlevel! equ 0 (
        echo 当前状态: 认证已禁用
    ) else (
        findstr /C:"DISABLE_AUTH=false" "%ENV_FILE%" >nul
        if !errorlevel! equ 0 (
            echo 当前状态: 认证已启用
        ) else (
            echo 当前状态: 未配置 ^(默认启用认证^)
        )
    )
    
    findstr /C:"HIDE_OLLAMA=true" "%ENV_FILE%" >nul
    if !errorlevel! equ 0 (
        echo Ollama 设置: 已隐藏
    ) else (
        findstr /C:"HIDE_OLLAMA=false" "%ENV_FILE%" >nul
        if !errorlevel! equ 0 (
            echo Ollama 设置: 已显示
        ) else (
            echo Ollama 设置: 未配置 ^(默认显示^)
        )
    )
    goto :end
)

echo 用法: %0 [enable^|disable^|status]
echo.
echo 命令:
echo   enable   - 启用用户认证
echo   disable  - 禁用用户认证
echo   status   - 查看当前状态
echo.
echo 示例:
echo   %0 disable  # 禁用认证，适合内网部署
echo   %0 enable   # 启用认证，适合生产环境
echo   %0 status   # 查看当前认证状态
exit /b 1

:end
endlocal
