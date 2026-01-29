# Local CI Simulation for Windows (Strict)
$ErrorActionPreference = "Stop"
# Toolchain Paths
$GO_BIN = "C:\Users\Stran\go\pkg\mod\golang.org\toolchain@v0.0.1-go1.24.0.windows-amd64\bin\go.exe"
# Use GOFUMPT strictly instead of GOFMT
$GOFUMPT_BIN = "C:\Users\Stran\go\bin\gofumpt.exe"
# Local Bin Paths
$GCI_BIN = "C:\Users\Stran\go\bin\gci.exe"

Write-Host "[START] Local CI Simulation (Strict Mode)..." -ForegroundColor Cyan

# 1. Formatting Check (gofumpt) - STRICTER than gofmt
Write-Host "[CHECK] Formatting (gofumpt)..." -ForegroundColor Yellow
$fmtOut = & $GOFUMPT_BIN -l .
if ($fmtOut) {
    Write-Host "[FAIL] Formatting errors found (Gofumpt Strictness):" -ForegroundColor Red
    $fmtOut
    exit 1
}
Write-Host "[PASS] Formatting OK." -ForegroundColor Green

# 2. Import Sorting Check (gci)
Write-Host "[CHECK] Import Sorting (gci)..." -ForegroundColor Yellow
if (Test-Path $GCI_BIN) {
    $gciOut = & $GCI_BIN diff -s standard -s default -s "Prefix(github.com/google/go-github/v82)" github otel example
    if ($gciOut) {
        Write-Host "[FAIL] Import sorting errors found:" -ForegroundColor Red
        $gciOut
        exit 1
    }
    Write-Host "[PASS] Import Sorting OK." -ForegroundColor Green
}
else {
    Write-Host "[FAIL] 'gci.exe' not found at $GCI_BIN. Verification incomplete." -ForegroundColor Red
    exit 1
}

# 3. Vet (Static Analysis)
Write-Host "[CHECK] Go Vet..." -ForegroundColor Yellow
& $GO_BIN vet ./...
if ($LASTEXITCODE -ne 0) {
    Write-Host "[FAIL] Go Vet Failed." -ForegroundColor Red
    exit 1
}
Write-Host "[PASS] Go Vet OK." -ForegroundColor Green

# 4. Generation Consistency Check
Write-Host "[CHECK] Verifying Generated Files..." -ForegroundColor Yellow
if (Test-Path "github") {
    Push-Location "github"
    & $GO_BIN run gen-accessors.go
    if ($LASTEXITCODE -ne 0) { Write-Host "[FAIL] gen-accessors failed" -ForegroundColor Red; exit 1 }

    & $GO_BIN run gen-stringify-test.go
    if ($LASTEXITCODE -ne 0) { Write-Host "[FAIL] gen-stringify-test failed" -ForegroundColor Red; exit 1 }
    
    Pop-Location
}
else {
    Write-Host "[FAIL] github directory not found." -ForegroundColor Red
    exit 1
}

# Check for git changes (dirty state) IGNORING local_ci.ps1
$status = & git status --porcelain | Select-String -NotMatch "local_ci.ps1"
if ($status) {
    Write-Host "[FAIL] Generated files are out of sync (Dirty Git Status):" -ForegroundColor Red
    $status
    Write-Host "[FIX] Commit these changes before pushing." -ForegroundColor Yellow
    exit 1
}
Write-Host "[PASS] Generated Files Clean." -ForegroundColor Green

# 5. Build/Test Check
Write-Host "[CHECK] Build Test (Dry Run)..." -ForegroundColor Yellow
& $GO_BIN build ./...
if ($LASTEXITCODE -ne 0) {
    Write-Host "[FAIL] Build Failed." -ForegroundColor Red
    exit 1
}
Write-Host "[PASS] Build OK." -ForegroundColor Green

Write-Host "[SUCCESS] Ready for PR Review." -ForegroundColor Green
exit 0
