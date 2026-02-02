# Session Rename Feature Verification Script
# Verifies all necessary files and code changes

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Session Rename Feature Verification" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$errors = 0

# Check if file exists
function Test-FileExists {
    param($path, $description)
    if (Test-Path $path) {
        Write-Host "[OK] $description" -ForegroundColor Green
        return $true
    } else {
        Write-Host "[FAIL] $description - File not found: $path" -ForegroundColor Red
        $script:errors++
        return $false
    }
}

# Check file content
function Test-FileContent {
    param($path, $pattern, $description)
    if (Test-Path $path) {
        $content = Get-Content $path -Raw
        if ($content -match $pattern) {
            Write-Host "[OK] $description" -ForegroundColor Green
            return $true
        } else {
            Write-Host "[FAIL] $description - Pattern not found" -ForegroundColor Red
            $script:errors++
            return $false
        }
    } else {
        Write-Host "[FAIL] $description - File not found" -ForegroundColor Red
        $script:errors++
        return $false
    }
}

Write-Host "1. Checking core files..." -ForegroundColor Yellow
Write-Host ""

# Check API file
Test-FileExists "frontend/src/api/chat/index.ts" "API file exists"
Test-FileContent "frontend/src/api/chat/index.ts" "updateSession" "API contains updateSession function"

# Check menu component
Test-FileExists "frontend/src/components/menu.vue" "Menu component exists"
Test-FileContent "frontend/src/components/menu.vue" "renameSession" "Menu contains renameSession function"
Test-FileContent "frontend/src/components/menu.vue" "handleRenameConfirm" "Menu contains handleRenameConfirm function"
Test-FileContent "frontend/src/components/menu.vue" "DialogPlugin" "Menu imports DialogPlugin"

Write-Host ""
Write-Host "2. Checking i18n files..." -ForegroundColor Yellow
Write-Host ""

# Check Chinese translations
Test-FileExists "frontend/src/i18n/locales/zh-CN.ts" "Chinese i18n file exists"
Test-FileContent "frontend/src/i18n/locales/zh-CN.ts" "renameSession" "Chinese i18n contains renameSession"

# Check English translations
Test-FileExists "frontend/src/i18n/locales/en-US.ts" "English i18n file exists"
Test-FileContent "frontend/src/i18n/locales/en-US.ts" "renameSession" "English i18n contains renameSession"

Write-Host ""
Write-Host "3. Checking test files..." -ForegroundColor Yellow
Write-Host ""

Test-FileExists "frontend/src/api/chat/__tests__/session.test.ts" "Unit test file exists"

Write-Host ""
Write-Host "4. Checking documentation..." -ForegroundColor Yellow
Write-Host ""

Test-FileExists "docs/features/session-rename.md" "Feature documentation exists"
Test-FileExists "docs/quick-start/rename-session-guide.md" "Quick start guide exists"
Test-FileExists "CHANGELOG_SESSION_RENAME.md" "Changelog exists"
Test-FileExists "SESSION_RENAME_SUMMARY.md" "Summary exists"

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Verification Results" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

if ($errors -eq 0) {
    Write-Host ""
    Write-Host "All checks passed!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Session rename feature successfully implemented." -ForegroundColor Green
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Yellow
    Write-Host "1. Run frontend dev server: cd frontend && npm run dev" -ForegroundColor White
    Write-Host "2. Test rename feature in browser" -ForegroundColor White
    Write-Host "3. Run unit tests: cd frontend && npm test" -ForegroundColor White
    Write-Host ""
    exit 0
} else {
    Write-Host ""
    Write-Host "Found $errors error(s)" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please check and fix the errors above." -ForegroundColor Red
    Write-Host ""
    exit 1
}
