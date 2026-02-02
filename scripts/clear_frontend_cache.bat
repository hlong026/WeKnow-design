@echo off
REM 清理前端缓存并重启前端服务

echo ==========================================
echo 清理前端缓存并重启
echo ==========================================

echo.
echo 步骤 1/2: 重启前端服务...
docker compose -f docker-compose.yml restart frontend

echo.
echo 步骤 2/2: 等待服务启动...
timeout /t 5 /nobreak >nul

echo.
echo ==========================================
echo 完成！请按以下步骤操作：
echo ==========================================
echo.
echo 1. 在浏览器中按 Ctrl+Shift+Delete 打开清除浏览数据
echo 2. 选择"缓存的图片和文件"和"Cookie及其他网站数据"
echo 3. 点击"清除数据"
echo 4. 刷新页面（Ctrl+F5 强制刷新）
echo.
echo 或者使用无痕模式访问: Ctrl+Shift+N (Chrome) 或 Ctrl+Shift+P (Firefox)
echo.
echo 访问地址: http://localhost
echo.

pause
