@echo off
REM 配置受保护的内置智能体 (Windows版本)
REM 用法: scripts\configure_protected_agents.bat [default|all|none|custom|status]

setlocal enabledelayedexpansion

set ENV_FILE=.env

if not exist "%ENV_FILE%" (
    echo 错误: .env 文件不存在
    echo 请先复制 .env.example 为 .env
    exit /b 1
)

if "%1"=="default" (
    echo 设置为默认配置（只保护快速问答）...
    findstr /C:"PROTECTED_BUILTIN_AGENTS=" "%ENV_FILE%" >nul
    if !errorlevel! equ 0 (
        powershell -Command "(Get-Content '%ENV_FILE%') -replace '^PROTECTED_BUILTIN_AGENTS=.*', 'PROTECTED_BUILTIN_AGENTS=builtin-quick-answer' | Set-Content '%ENV_FILE%'"
    ) else (
        echo PROTECTED_BUILTIN_AGENTS=builtin-quick-answer >> "%ENV_FILE%"
    )
    echo ✓ 已设置为默认配置（只保护快速问答）
    echo.
    echo 受保护的智能体:
    echo   - builtin-quick-answer ^(快速问答^)
    goto :end
)

if "%1"=="all" (
    echo 保护所有内置智能体...
    findstr /C:"PROTECTED_BUILTIN_AGENTS=" "%ENV_FILE%" >nul
    if !errorlevel! equ 0 (
        powershell -Command "(Get-Content '%ENV_FILE%') -replace '^PROTECTED_BUILTIN_AGENTS=.*', 'PROTECTED_BUILTIN_AGENTS=builtin-quick-answer,builtin-smart-reasoning,builtin-data-analyst' | Set-Content '%ENV_FILE%'"
    ) else (
        echo PROTECTED_BUILTIN_AGENTS=builtin-quick-answer,builtin-smart-reasoning,builtin-data-analyst >> "%ENV_FILE%"
    )
    echo ✓ 已保护所有内置智能体
    echo.
    echo 受保护的智能体:
    echo   - builtin-quick-answer ^(快速问答^)
    echo   - builtin-smart-reasoning ^(智能推理^)
    echo   - builtin-data-analyst ^(数据分析师^)
    goto :end
)

if "%1"=="none" (
    echo 允许删除所有智能体...
    findstr /C:"PROTECTED_BUILTIN_AGENTS=" "%ENV_FILE%" >nul
    if !errorlevel! equ 0 (
        powershell -Command "(Get-Content '%ENV_FILE%') -replace '^PROTECTED_BUILTIN_AGENTS=.*', 'PROTECTED_BUILTIN_AGENTS=' | Set-Content '%ENV_FILE%'"
    ) else (
        echo PROTECTED_BUILTIN_AGENTS= >> "%ENV_FILE%"
    )
    echo ✓ 已允许删除所有智能体
    echo.
    echo ⚠️  警告：所有内置智能体都可以被删除
    echo ⚠️  请谨慎使用此配置
    goto :end
)

if "%1"=="custom" (
    echo 自定义配置
    echo.
    echo 请手动编辑 .env 文件中的 PROTECTED_BUILTIN_AGENTS 变量
    echo.
    echo 可用的智能体 ID:
    echo   - builtin-quick-answer      ^(快速问答^)
    echo   - builtin-smart-reasoning   ^(智能推理^)
    echo   - builtin-data-analyst      ^(数据分析师^)
    echo.
    echo 示例:
    echo   PROTECTED_BUILTIN_AGENTS=builtin-quick-answer,builtin-smart-reasoning
    exit /b 0
)

if "%1"=="status" (
    findstr /C:"PROTECTED_BUILTIN_AGENTS=" "%ENV_FILE%" >nul
    if !errorlevel! equ 0 (
        for /f "tokens=2 delims==" %%a in ('findstr /C:"PROTECTED_BUILTIN_AGENTS=" "%ENV_FILE%"') do set value=%%a
        if "!value!"=="" (
            echo 当前配置: 允许删除所有智能体
        ) else (
            echo 当前配置: 受保护的智能体
            echo !value! | findstr /C:"builtin-quick-answer" >nul && echo   - builtin-quick-answer ^(快速问答^)
            echo !value! | findstr /C:"builtin-smart-reasoning" >nul && echo   - builtin-smart-reasoning ^(智能推理^)
            echo !value! | findstr /C:"builtin-data-analyst" >nul && echo   - builtin-data-analyst ^(数据分析师^)
        )
    ) else (
        echo 当前配置: 未配置 ^(默认保护快速问答^)
    )
    exit /b 0
)

echo 用法: %0 [default^|all^|none^|custom^|status]
echo.
echo 命令:
echo   default  - 只保护快速问答（默认配置）
echo   all      - 保护所有内置智能体
echo   none     - 允许删除所有智能体（谨慎使用）
echo   custom   - 显示自定义配置说明
echo   status   - 查看当前配置
echo.
echo 示例:
echo   %0 default  # 设置为默认配置
echo   %0 all      # 保护所有内置智能体
echo   %0 status   # 查看当前配置
exit /b 1

:end
echo.
echo 提示: 需要重启服务才能生效
echo 运行: scripts\start_all.bat --stop ^&^& scripts\start_all.bat
endlocal
