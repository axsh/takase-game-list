Param(
	[string]$Mode = "all"
)

$ErrorActionPreference = "Stop"

function Write-Info($Message) {
	Write-Host "[INFO] $Message"
}

function Write-FailAndExit($Message) {
	Write-Error "[ERROR] $Message"
	exit 1
}

function Invoke-Go {
	param(
		[string[]]$GoArgs
	)

	& go @GoArgs
	if ($LASTEXITCODE -ne 0) {
		throw "go $($GoArgs -join ' ') failed with exit code $LASTEXITCODE"
	}
}

function Ensure-Bin {
	param(
		[string]$BinDir
	)
	if (-not (Test-Path $BinDir)) {
		New-Item -ItemType Directory -Path $BinDir | Out-Null
	}
}

$ValidModes = @("all", "unit", "integration", "build-only")
if (-not ($ValidModes -contains $Mode)) {
	Write-FailAndExit "Invalid mode '$Mode'. Use one of: $($ValidModes -join ', ')"
}

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $ScriptDir

$BinDir = Join-Path $ScriptDir "bin"
$BinaryPath = Join-Path $BinDir "game-list.exe"

Write-Info "Mode: $Mode"
Ensure-Bin -BinDir $BinDir

try {
	Write-Info "Building binary -> $BinaryPath"
	Invoke-Go -GoArgs @("build", "-o", $BinaryPath, ".")
	Write-Info "Build succeeded"

	if ($Mode -eq "unit" -or $Mode -eq "all") {
		Write-Info "Running unit tests"
		Invoke-Go -GoArgs @("test", "-run", "_Unit$", "./...")
		Write-Info "Unit tests passed"
	}

	if ($Mode -eq "integration" -or $Mode -eq "all") {
		Write-Info "Running integration tests"
		Invoke-Go -GoArgs @("test", "-run", "_Integration$", "./tests/integration/...")
		Write-Info "Integration tests passed"
	}

	Write-Info "Completed mode '$Mode'"
} catch {
	Write-FailAndExit $_
}
