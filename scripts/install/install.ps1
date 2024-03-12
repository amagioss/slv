#!/usr/bin/env pwsh

$ErrorActionPreference = 'Stop'

if ($v) {
  $Version = "v${v}"
}
if ($Args.Length -eq 1) {
  $Version = $Args.Get(0)
}

$SLVInstall = $env:SLV_INSTALL
$BinDir = if ($SLVInstall) {
  "${SLVInstall}\bin"
} else {
  "${Home}\.slv\bin"
}

$SLVZip = "$BinDir\slv.zip"
$SLVExe = "$BinDir\slv.exe"
$Target = 'windows_amd64'

$DownloadUrl = if (!$Version) {
  "https://github.com/savesecrets/slv/releases/latest/download/slv_${Target}.zip"
} else {
  "https://github.com/savesecrets/slv/releases/download/${Version}/slv_${Target}.zip"
}

if (!(Test-Path $BinDir)) {
  New-Item $BinDir -ItemType Directory | Out-Null
}

curl.exe -Lo $SLVZip $DownloadUrl

tar.exe xf $SLVZip -C $BinDir

Remove-Item $SLVZip

$User = [System.EnvironmentVariableTarget]::User
$Path = [System.Environment]::GetEnvironmentVariable('Path', $User)
if (!(";${Path};".ToLower() -like "*;${BinDir};*".ToLower())) {
  [System.Environment]::SetEnvironmentVariable('Path', "${Path};${BinDir}", $User)
  $Env:Path += ";${BinDir}"
}

Write-Output "SLV was installed successfully to ${SLVExe}"
Write-Output "Run 'slv --help' to get started"
