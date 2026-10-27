# Regression test for O-04 (SAR-AUTOSETUP-001-O04-DEV): Get-Registered must pick the KnowsLink entry from a host plugin list of 0, 1 or many items.
# It runs the installer's own Norm, Get-HostJson and Get-Registered against a stub host CLI; run it under powershell.exe (5.1) and pwsh (7).
$ErrorActionPreference = 'Stop'
$installer = Join-Path $PSScriptRoot 'knowslink-setup.ps1'
# ParseFile would read the BOM-less UTF-8 installer as ANSI under Windows PowerShell 5.1 and garble its Korean quoting.
$ast = [Management.Automation.Language.Parser]::ParseInput([IO.File]::ReadAllText($installer, [Text.Encoding]::UTF8), [ref]$null, [ref]$null)
foreach ($name in 'Norm', 'Get-HostJson', 'Get-Registered') {
    . ([scriptblock]::Create(($ast.FindAll({ param($n) $n -is [Management.Automation.Language.FunctionDefinitionAst] -and $n.Name -eq $name }, $true) | Select-Object -First 1).Extent.Text))
}

# The stub host is a .ps1 so it runs on any OS; the list it prints is rewritten per case.
$work = Join-Path ([IO.Path]::GetTempPath()) ("knowslink-registered-test-" + [guid]::NewGuid().ToString('n'))
New-Item -ItemType Directory -Path $work | Out-Null
$stub = Join-Path $work 'host.ps1'
Set-Content -LiteralPath $stub -Value 'Get-Content -Raw -LiteralPath (Join-Path $PSScriptRoot "list.json"); exit 0'
$HostCommand = [pscustomobject]@{ Source = $stub }
$Marketplace = 'knowslink-plugins'
$Plugin = 'knowslink@knowslink-plugins'
$failed = 0

function Test-Case([string]$AppName, [string]$Name, [string]$Json, $Expected) {
    $script:App = $AppName
    Set-Content -LiteralPath (Join-Path $work 'list.json') -Value $Json
    $got = Get-Registered
    if ($got -eq $Expected) { Write-Host "PASS $AppName $Name" } else { Write-Host "FAIL $AppName $Name : expected [$Expected] got [$got]"; $script:failed++ }
}

$mine = Norm 'C:/x/claude-01'
$me = '{"id":"knowslink@knowslink-plugins","readFromFolder":"C:/x/claude-01/knowslink"}'
$other = '{"id":"other@inline"}'
$nofolder = '{"id":"orca-agent-view@inline","readFromFolder":null}'
Test-Case claude '0 plugins' '[]' $null
Test-Case claude '1 plugin (knowslink)' "[$me]" $mine
Test-Case claude '1 plugin (other)' "[$other]" $null
Test-Case claude '2 plugins (knowslink first)' "[$me,$other]" $mine
Test-Case claude '2 plugins (knowslink last)' "[$other,$me]" $mine
Test-Case claude '3 plugins incl. null readFromFolder' "[$nofolder,$other,$me]" $mine
Test-Case claude '2 plugins without knowslink' "[$other,$nofolder]" $null
$src = '{"marketplaceName":"knowslink-plugins","marketplaceSource":{"source":"C:/x/codex-01"}}'
$mineCodex = Norm 'C:/x/codex-01'
Test-Case codex '0 entries' '{"installed":[],"available":[]}' $null
Test-Case codex '1 installed' "{`"installed`":[$src],`"available`":[]}" $mineCodex
Test-Case codex 'many entries' "{`"installed`":[{`"marketplaceName`":`"other`"},$src],`"available`":[{`"marketplaceName`":`"x`"}]}" $mineCodex

Remove-Item -LiteralPath $work -Recurse -Force
if ($failed) { Write-Host "$failed case(s) failed"; exit 1 }
Write-Host "OK ($($PSVersionTable.PSEdition) $($PSVersionTable.PSVersion))"
