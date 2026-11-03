# Install the yarbrain CLI from a GitHub release on Windows, after verifying
# its SHA-256 against the release's checksums.txt.
#
#   YARBRAIN_VERSION      version to install, such as 0.3.0 (default: latest)
#   YARBRAIN_INSTALL_DIR  target directory (default: %LOCALAPPDATA%\Programs\yarbrain)
$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'

$repo = 'yarlson/yarstack'
$tagPrefix = 'plugins/yarbrain/cli/v'
# YARBRAIN_RELEASE_BASE and YARBRAIN_RELEASE_API exist for tests only.
$base = if ($env:YARBRAIN_RELEASE_BASE) { $env:YARBRAIN_RELEASE_BASE } else { "https://github.com/$repo/releases/download" }
$api = if ($env:YARBRAIN_RELEASE_API) { $env:YARBRAIN_RELEASE_API } else { "https://api.github.com/repos/$repo/releases?per_page=50" }
$installDir = if ($env:YARBRAIN_INSTALL_DIR) { $env:YARBRAIN_INSTALL_DIR } else { Join-Path $env:LOCALAPPDATA 'Programs\yarbrain' }

function Fail([string]$message) {
    [Console]::Error.WriteLine("yarbrain install: $message")
    exit 1
}

$version = $env:YARBRAIN_VERSION
if (-not $version) {
    try {
        $releases = Invoke-RestMethod -Uri $api -TimeoutSec 30
    } catch {
        Fail "could not read the release list from GitHub: $_"
    }
    $release = $releases | Where-Object { $_.tag_name -like "$tagPrefix*" } | Select-Object -First 1
    if (-not $release) { Fail 'no yarbrain release is published yet' }
    $version = $release.tag_name.Substring($tagPrefix.Length)
}

$tag = "$tagPrefix$version"
# Only an amd64 build is published; Windows on arm64 runs it under emulation.
$asset = "yarbrain_${version}_windows_amd64.exe"
$work = Join-Path ([System.IO.Path]::GetTempPath()) ("yarbrain-" + [System.Guid]::NewGuid())
New-Item -ItemType Directory -Path $work | Out-Null

try {
    try {
        Invoke-WebRequest -Uri "$base/$tag/$asset" -OutFile (Join-Path $work $asset) -TimeoutSec 60 -UseBasicParsing
        Invoke-WebRequest -Uri "$base/$tag/checksums.txt" -OutFile (Join-Path $work 'checksums.txt') -TimeoutSec 60 -UseBasicParsing
    } catch {
        Fail "could not download $asset from release ${tag}: $_"
    }

    $expected = $null
    foreach ($line in Get-Content (Join-Path $work 'checksums.txt')) {
        $fields = $line -split '\s+'
        if ($fields.Length -ge 2 -and $fields[1] -eq $asset) { $expected = $fields[0].ToLowerInvariant() }
    }
    if (-not $expected) { Fail "checksums.txt of release $tag lists no $asset" }

    $actual = (Get-FileHash -Algorithm SHA256 (Join-Path $work $asset)).Hash.ToLowerInvariant()
    if ($expected -ne $actual) { Fail "checksum mismatch for $asset; nothing was installed" }

    New-Item -ItemType Directory -Force -Path $installDir | Out-Null
    $target = Join-Path $installDir 'yarbrain.exe'

    if ((Test-Path $target) -and ((Get-FileHash -Algorithm SHA256 $target).Hash.ToLowerInvariant() -eq $expected)) {
        Write-Output "yarbrain $version is already installed at $target"
    } else {
        $staged = Join-Path $installDir ".yarbrain.$PID.exe"
        Copy-Item (Join-Path $work $asset) $staged -Force
        Move-Item $staged $target -Force
        Write-Output "Installed yarbrain $version at $target"
    }
} finally {
    Remove-Item -Recurse -Force $work -ErrorAction SilentlyContinue
}

$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if (-not (($userPath -split ';') -contains $installDir)) {
    $parts = @($userPath, $installDir) | Where-Object { $_ }
    [Environment]::SetEnvironmentVariable('Path', ($parts -join ';'), 'User')
    Write-Output "Added $installDir to your user PATH; open a new terminal to use it."
}

Write-Output 'Next: ask your agent to set up Yarbrain, or run: yarbrain init <vault-path>'
