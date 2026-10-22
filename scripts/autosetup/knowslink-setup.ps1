# KnowsLink automatic setup for Claude Code and Codex CLI on Windows 11 (SAR-AUTOSETUP-001, D02 AS-01-AS-13).
# scripts/package_plugin.py --release fills the pinned version, URLs and SHA-256 values; never edit a released copy.
# Run: & ([scriptblock]::Create((irm <base>/knowslink-setup.ps1))) -App claude|codex -Lang ko|en
# The installer never reads or prints keys, tokens or host sign-in files, and changes only the KnowsLink host entries.
param(
    [Parameter(Mandatory = $true)][ValidateSet('claude', 'codex')][string]$App,
    [ValidateSet('ko', 'en')][string]$Lang = 'ko',
    [switch]$NewAgent
)
$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'
$Version = '__VERSION__'
$Base = '__BASE__'
$Relay = '__RELAY__'
$Sums = @{ claude = '__SHA_CLAUDE__'; codex = '__SHA_CODEX__' }
$OpenBrowser = $true # __OPEN_BROWSER__
$NodeVersion = '22.22.2'
$NodeSums = @{ x64 = '7c93e9d92bf68c07182b471aa187e35ee6cd08ef0f24ab060dfff605fcc1c57c'; arm64 = '380d375cf650c5a7f2ef3ce29ac6ea9a1c9d2ec8ea8e8391e1a34fd543886ab3' }
$Tested = @{ claude = '2.1.302'; codex = '0.163.0' }
$Names = @{ claude = 'Claude Code'; codex = 'Codex CLI' }
$Marketplace = 'knowslink-plugins'
$Plugin = 'knowslink@knowslink-plugins'

$Text = @{
    ko = @{
        title = 'KnowsLink 자동 설정 {0} — {1}'
        checking = '확인 중: 호스트·버전·OS·기존 설치'
        os = 'OS: {0} (빌드 {1})'
        hostVersion = '{0}: {1}'
        unverified = '검증하지 않은 환경입니다: {0}. 검증한 환경은 {1} {2}, Windows 11입니다.'
        unverifiedAsk = '검증하지 않은 환경에서 계속할까요? (y/N)'
        modeNew = '작업: 신규 설치'
        modeUpdate = '작업: 기존 설치 갱신 (현재 {0} → 새 {1})'
        modeResume = '작업: 이전 실행 이어서 진행 (남은 단계부터)'
        keepAgent = '유지하는 에이전트: {0}'
        keepAgentUnknown = '유지하는 에이전트: 이 폴더의 에이전트 (이름·ID는 웹 ''내 에이전트''에서 확인)'
        planDir = '설치 폴더: {0}'
        planPrivateNew = '에이전트 폴더(키는 이 컴퓨터의 이 폴더에만 저장): {0} — Google 연결 때 만듭니다'
        planPrivateKeep = '에이전트 폴더: {0} — 읽거나 바꾸지 않습니다'
        planNode = 'Node 실행 파일: {0} (없으면 내려받아 확인합니다)'
        planHost = '바꿀 {0} 설정: {1} — KnowsLink 항목만: {2}'
        planHostKeep = '다른 MCP 항목·승인 모드·모델·로그인 정보는 바꾸지 않습니다.'
        planReplace = '기존 KnowsLink 등록({0})을 이 설치로 바꿉니다.'
        planGoogle = '설치 뒤 Google 연결 링크를 열고 지문을 보여 줍니다.'
        planPrevious = '갱신 전 설치는 {0} 로 보존합니다.'
        consent = '이 계획대로 진행할까요? (y/N)'
        declined = '진행하지 않았습니다. 바뀐 것이 없습니다.'
        pick = '이 컴퓨터에 {0} 설치가 여러 개 있습니다:'
        pickAsk = '갱신할 설치 번호를 입력하세요 (그 밖의 입력은 종료):'
        pickUnsure = '[확인 필요] 어느 설치인지 확실하지 않아 멈췄습니다. 바뀐 것이 없습니다. 새 에이전트로 설치하려면 명령 끝에 -NewAgent 를 붙이세요.'
        stepOk = '완료: {0}'
        sDownload = '내려받기와 checksum 확인'
        sNode = 'Node {0} 준비'
        sFiles = '설치 폴더 준비'
        sRegister = '{0}에 KnowsLink 등록과 등록 확인'
        sGoogle = 'Google 연결과 키 저장 확인'
        googleOpen = '자기 브라우저에서 Google로 로그인하세요: {0}'
        googleResume = '이전 연결 시도를 이어서 기다립니다: {0}'
        googleNoBrowser = '브라우저를 열지 못했습니다. 위 링크를 자기 브라우저에서 직접 여세요.'
        googlePrint = '이 컴퓨터의 지문: {0}'
        googleCompare = '브라우저에 보이는 지문이 위 지문과 같을 때만 에이전트 이름을 정하고 동의하세요. 다른 사람이 보낸 연결 링크는 승인하지 마세요.'
        googleWait = '이 창에서 승인을 기다립니다 (기한 {0}).'
        googleDone = '키 저장 확인: 에이전트 {0} · 지문 {1}'
        f07Ask = '웹 ''내 에이전트''에서 에이전트 {0} 를 확인(쓰지 않으면 철회)했으면 new 를 입력해 새 연결을 시작하세요. 그 밖의 입력은 종료합니다:'
        changed = '이미 바뀐 것: {0}'
        next = '다음 행동: {0}'
        doneNew = '설치와 키 저장 완료'
        doneUpdate = '갱신 완료. 에이전트·키·친구는 그대로입니다.'
        doneAgent = '에이전트: {0}'
        restart = '{0}를 다시 열어야 KnowsLink가 로드됩니다. 실행 중인 {0} 창은 모두 닫고 다시 여세요.'
        sessionClaude = '세션 연결: 메시지를 받을 Claude Code 세션에서 /knowslink-session on 을 실행하세요. 보낼 때마다 발송 승인 대화상자에서 승인합니다.'
        sessionCodex = '세션 연결: 승인을 묻는 모드(예: approval_policy = "on-request")로 Codex를 열고 knowslink_session 을 action on 으로 실행하게 한 뒤 Codex 도구 승인과 4자리 확인 코드를 입력하세요. YOLO(전체 자동) 모드에서는 연결되지 않습니다. Codex 화면에는 도착 알림이 없습니다.'
        testedEnv = '시험 환경 확인: {0} {1}, Windows 11'
        removeTitle = '제거 안내 (제거 명령은 없습니다):'
        removeDir = '설치 폴더 {0}, 에이전트 폴더 {1}'
        removeHost = '{0} 설정 항목: {1}. 해제: {2}'
        removeKey = '키 철회: 웹 ''내 에이전트 → 연결 관리''에서 이 에이전트를 철회하세요.'
        record = '설치 기록: {0}'
        F01 = '{0}를 찾거나 실행할 수 없습니다.'
        F01n = '{0}를 설치하고 로그인한 뒤 새 PowerShell에서 ''{1} --version''이 버전을 출력하는지 확인하고 같은 명령을 다시 실행하세요.'
        F02 = '검증하지 않은 환경에서 계속하지 않았습니다.'
        F02n = '검증한 환경은 {0} {1}, Windows 11입니다. 그 전까지는 웹 ''웹에서 바로 사용''으로 대화하세요.'
        F03 = '내려받지 못했습니다: {0}'
        F03n = '인터넷 연결을 확인하고 같은 명령을 다시 실행하세요.'
        F04 = '내려받은 파일을 확인할 수 없어 멈췄습니다 (checksum 불일치).'
        F04n = '같은 명령을 다시 실행하세요. 계속되면 회원 화면의 명령을 다시 복사하세요.'
        F05 = '{0}에 KnowsLink를 등록하지 못했습니다: {1}'
        F05n = '''{0} plugin list''로 상태를 확인하고 같은 명령을 다시 실행하세요. 등록부터 이어서 진행합니다.'
        F06 = '연결하지 않았습니다. 다시 실행하면 새 연결을 시작합니다.'
        F06n = '같은 명령을 다시 실행하세요. Google 연결부터 이어서 진행합니다.'
        F07 = '확인 필요. 웹 ''내 에이전트''에서 에이전트 {0} 를 먼저 확인하세요. 키 저장 결과를 확인하지 못했습니다.'
        F07n = '같은 명령을 다시 실행하면 저장 결과부터 다시 확인합니다. 확인 전에는 새 연결을 시작하지 않습니다.'
        F08 = '이전 버전을 계속 사용합니다. 새 버전으로 바꾸지 못했습니다: {0}'
        F08n = '{0} 창을 모두 닫고 같은 명령을 다시 실행하세요. 갱신부터 이어서 진행합니다.'
        F09 = '이 설정에서는 연결할 수 없습니다. Codex 승인 모드가 YOLO(approval_policy = "never")입니다.'
        F09n = '{0} 의 approval_policy 를 "on-request"로 바꾸거나 해당 줄을 지운 뒤 Codex를 여세요. 설치기는 승인 모드를 바꾸지 않습니다.'
        cNone = '없음'
        cTemp = '없음 (임시 파일 삭제)'
        cFolder = '설치 폴더 {0} 있음, {1} 등록 없음 또는 이전 등록'
        cNoKey = '설치와 등록 있음, 키 없음'
        cUnsure = '불확실 (에이전트가 만들어졌을 수 있음)'
        cKept = '없음 (이전 설치 유지)'
        cDone = '설치 완료'
    }
    en = @{
        title = 'KnowsLink automatic setup {0} — {1}'
        checking = 'Checking host, version, OS and existing installs'
        os = 'OS: {0} (build {1})'
        hostVersion = '{0}: {1}'
        unverified = 'This environment is not verified: {0}. Verified: {1} {2}, Windows 11.'
        unverifiedAsk = 'Continue in an unverified environment? (y/N)'
        modeNew = 'Action: new install'
        modeUpdate = 'Action: update existing install (current {0} → new {1})'
        modeResume = 'Action: continue the previous run from the remaining steps'
        keepAgent = 'Agent kept: {0}'
        keepAgentUnknown = 'Agent kept: the agent of this folder (see name and ID in ''My agents'' on the web)'
        planDir = 'Install folder: {0}'
        planPrivateNew = 'Agent folder (keys are stored only here on this computer): {0} — created at Google connection'
        planPrivateKeep = 'Agent folder: {0} — not read or changed'
        planNode = 'Node runtime: {0} (downloaded and checked if missing)'
        planHost = '{0} settings to change: {1} — KnowsLink entries only: {2}'
        planHostKeep = 'Other MCP entries, approval mode, model and sign-in data are not changed.'
        planReplace = 'The existing KnowsLink registration ({0}) will point to this install.'
        planGoogle = 'After installing, it opens the Google connection link and shows the fingerprint.'
        planPrevious = 'The install before the update is kept as {0}.'
        consent = 'Proceed with this plan? (y/N)'
        declined = 'Not started. Nothing was changed.'
        pick = 'This computer has several {0} installs:'
        pickAsk = 'Enter the number of the install to update (anything else stops):'
        pickUnsure = '[Check needed] Stopped because it is unclear which install to update. Nothing was changed. To install a new agent, add -NewAgent to the end of the command.'
        stepOk = 'Done: {0}'
        sDownload = 'Download and checksum'
        sNode = 'Prepare Node {0}'
        sFiles = 'Prepare install folder'
        sRegister = 'Add KnowsLink to {0} and check the registration'
        sGoogle = 'Google connection and key storage check'
        googleOpen = 'Sign in with Google in your own browser: {0}'
        googleResume = 'Waiting for the previous connection attempt: {0}'
        googleNoBrowser = 'Could not open a browser. Open the link above in your own browser.'
        googlePrint = 'Fingerprint on this computer: {0}'
        googleCompare = 'Name the agent and agree only if the browser shows the same fingerprint. Never approve a connection link someone else sent you.'
        googleWait = 'Waiting here for your approval (until {0}).'
        googleDone = 'Key storage confirmed: agent {0} · fingerprint {1}'
        f07Ask = 'After checking agent {0} in ''My agents'' on the web (revoke it if unused), type new to start a new connection. Anything else stops:'
        changed = 'Already changed: {0}'
        next = 'Next: {0}'
        doneNew = 'Install and key storage complete'
        doneUpdate = 'Update complete. Your agent, key and friends stay the same.'
        doneAgent = 'Agent: {0}'
        restart = 'Reopen {0} so it loads KnowsLink. Close every running {0} window first.'
        sessionClaude = 'Session: in the Claude Code session that should receive messages, run /knowslink-session on. Approve every send in the send dialog.'
        sessionCodex = 'Session: open Codex in a mode that asks for approval (for example approval_policy = "on-request"), ask it to run knowslink_session with action on, approve the Codex tool and enter the 4-digit code. Connection does not work in YOLO (full auto) mode. Codex shows no arrival notice.'
        testedEnv = 'Checked in a test environment: {0} {1}, Windows 11'
        removeTitle = 'Removal (there is no remove command):'
        removeDir = 'Install folder {0}, agent folder {1}'
        removeHost = '{0} settings entries: {1}. Unregister: {2}'
        removeKey = 'Revoke the key: revoke this agent in ''My agents → Manage connections'' on the web.'
        record = 'Install record: {0}'
        F01 = '{0} was not found or cannot run.'
        F01n = 'Install and sign in to {0}, check that ''{1} --version'' prints a version in a new PowerShell, then run the same command again.'
        F02 = 'Did not continue in an unverified environment.'
        F02n = 'Verified: {0} {1}, Windows 11. Until then, use ''Use directly on the web''.'
        F03 = 'Download failed: {0}'
        F03n = 'Check your internet connection and run the same command again.'
        F04 = 'Stopped because the downloaded file could not be verified (checksum mismatch).'
        F04n = 'Run the same command again. If it repeats, copy the command from the member page again.'
        F05 = 'Could not add KnowsLink to {0}: {1}'
        F05n = 'Check with ''{0} plugin list'' and run the same command again. It continues from registration.'
        F06 = 'Not connected. Running again starts a new connection.'
        F06n = 'Run the same command again. It continues from the Google connection.'
        F07 = 'Check needed. First check agent {0} in ''My agents'' on the web. The key storage result could not be confirmed.'
        F07n = 'Running the same command again checks the stored result first. It does not start a new connection before that.'
        F08 = 'Keeping the previous version. Could not switch to the new version: {0}'
        F08n = 'Close every {0} window and run the same command again. It continues from the update.'
        F09 = 'You cannot connect with this setting. Codex approval mode is YOLO (approval_policy = "never").'
        F09n = 'Change approval_policy in {0} to "on-request" or delete that line, then open Codex. The installer does not change the approval mode.'
        cNone = 'nothing'
        cTemp = 'nothing (temporary files deleted)'
        cFolder = 'install folder {0} exists; {1} registration missing or previous'
        cNoKey = 'install and registration exist; no key'
        cUnsure = 'uncertain (an agent may have been created)'
        cKept = 'nothing (previous install kept)'
        cDone = 'install complete'
    }
}[$Lang]

function Say([string]$Key) { Write-Host ($Text[$Key] -f $args) }
function T([string]$Key) { $Text[$Key] -f $args }
function Ask([string]$Key) {
    Write-Host -NoNewline ((T $Key @args) + ' ')
    $answer = [Console]::In.ReadLine()
    if ($null -eq $answer) { return '' }
    $answer.Trim()
}
function Norm([string]$Path) { if ($Path) { [IO.Path]::GetFullPath(($Path -replace '^\\\\\?\\', '')).TrimEnd('\') } }
function Stamp { (Get-Date).ToString('yyyyMMddHHmmss') }

# Shared state lives in this hashtable: `exit` and script-scope variables would reach the caller's session when run as a script block.
$S = @{ Record = $null; RecordDir = $null; Temp = $null }
function Save-Record([string]$Step, [string]$Result) {
    if (-not $S.Record) { return }
    $S.Record.steps = @(@($S.Record.steps) + @([ordered]@{ step = $Step; result = $Result; version = $Version; at = (Get-Date).ToString('o') }) | Select-Object -Last 50)
    $S.Record.updated = (Get-Date).ToString('o')
    if ($S.RecordDir -and (Test-Path -LiteralPath $S.RecordDir)) {
        $json = $S.Record | ConvertTo-Json -Depth 6
        [IO.File]::WriteAllText((Join-Path $S.RecordDir 'knowslink-install.json'), $json, (New-Object Text.UTF8Encoding $false))
    }
}
function Remove-Temp { if ($S.Temp -and (Test-Path -LiteralPath $S.Temp)) { Remove-Item -LiteralPath $S.Temp -Recurse -Force } }

function Stop-Setup([string]$Id, [string]$Message, [string]$Changed, [string]$Next) {
    Write-Host ''
    Write-Host "[$Id] $Message" -ForegroundColor Yellow
    Say changed $Changed
    Say next $Next
    Save-Record $Id 'failed'
    Remove-Temp
    throw 'knowslink-setup-stopped'
}

# Host CLI: the npm .ps1 shim can be blocked by execution policy, so only .exe/.cmd applications run.
$HostCommand = Get-Command $App -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1
function Invoke-HostCli([string[]]$CliArgs) {
    $ErrorActionPreference = 'Continue'
    $out = & $HostCommand.Source @CliArgs 2>&1 | ForEach-Object { "$_" } | Out-String
    @{ Code = $LASTEXITCODE; Out = $out }
}
function Get-HostJson([string[]]$CliArgs) {
    $ErrorActionPreference = 'Continue'
    $out = & $HostCommand.Source @CliArgs 2>$null | Out-String
    if ($LASTEXITCODE -ne 0) { return $null }
    try { $out | ConvertFrom-Json } catch { $null }
}

$HostHome = if ($App -eq 'claude') { if ($env:CLAUDE_CONFIG_DIR) { $env:CLAUDE_CONFIG_DIR } else { Join-Path $HOME '.claude' } } else { if ($env:CODEX_HOME) { $env:CODEX_HOME } else { Join-Path $HOME '.codex' } }
$HostSettings = if ($App -eq 'claude') { Join-Path $HostHome 'settings.json' } else { Join-Path $HostHome 'config.toml' }
$HostEntries = if ($App -eq 'claude') { "extraKnownMarketplaces.$Marketplace, enabledPlugins.$Plugin" } else { "[marketplaces.$Marketplace], [plugins.""$Plugin""]" }
$Unregister = if ($App -eq 'claude') { "claude plugin uninstall $Plugin; claude plugin marketplace remove $Marketplace" } else { "codex plugin remove $Plugin; codex plugin marketplace remove $Marketplace" }

function Get-Registered {
    if ($App -eq 'claude') {
        # `plugin marketplace list` rewrites the user's model alias in Claude Code 2.1.308; `plugin list` leaves settings alone.
        $entry = @(Get-HostJson @('plugin', 'list', '--json')) | Where-Object { $_.id -eq $Plugin } | Select-Object -First 1
        if ($entry -and $entry.readFromFolder) { return Norm (Split-Path -Parent $entry.readFromFolder) }
        return $null
    }
    $list = Get-HostJson @('plugin', 'list', '--available', '--json')
    $entry = @($list.installed) + @($list.available) | Where-Object { $_.marketplaceName -eq $Marketplace } | Select-Object -First 1
    if ($entry) { Norm $entry.marketplaceSource.source }
}

function Invoke-Setup {
Say title $Version $Names[$App]
Say checking
# N2 / AS-09: OS and host before any change.
$build = [Environment]::OSVersion.Version.Build
$caption = try { (Get-CimInstance Win32_OperatingSystem).Caption } catch { 'Windows' }
$isWindows11 = ([Environment]::OSVersion.Platform -eq 'Win32NT') -and $build -ge 22000
if ([Environment]::OSVersion.Platform -ne 'Win32NT') { Stop-Setup 'F-02' (T F02) (T cNone) (T F02n $Names[$App] $Tested[$App]) }
Say os $caption $build
if (-not $HostCommand) { Stop-Setup 'F-01' (T F01 $Names[$App]) (T cNone) (T F01n $Names[$App] $App) }
$probe = Invoke-HostCli @('--version')
$hostVersion = if ($probe.Out -match '(\d+\.\d+\.\d+)') { $Matches[1] } else { '' }
if ($probe.Code -ne 0 -or -not $hostVersion) { Stop-Setup 'F-01' (T F01 $Names[$App]) (T cNone) (T F01n $Names[$App] $App) }
Say hostVersion $Names[$App] $hostVersion
$verified = $isWindows11 -and $hostVersion -eq $Tested[$App]
if (-not $verified) {
    $what = @()
    if (-not $isWindows11) { $what += "$caption ($build)" }
    if ($hostVersion -ne $Tested[$App]) { $what += "$($Names[$App]) $hostVersion" }
    Say unverified ($what -join ', ') $Names[$App] $Tested[$App]
    if ((Ask unverifiedAsk) -notmatch '^(y|yes)$') { Stop-Setup 'F-02' (T F02) (T cNone) (T F02n $Names[$App] $Tested[$App]) }
}

# Existing installs: host registration plus install records in the default root (AS-10).
$Root = Join-Path $env:LOCALAPPDATA 'KnowsLink'
$registered = Get-Registered
$candidates = @()
if (Test-Path -LiteralPath $Root) {
    $candidates += Get-ChildItem -LiteralPath $Root -Directory | Where-Object { $_.Name -match "^$App-\d{2}$" -and (Test-Path -LiteralPath (Join-Path $_.FullName 'knowslink-install.json')) } | ForEach-Object { Norm $_.FullName }
}
if ($registered -and (Test-Path -LiteralPath (Join-Path $registered 'knowslink'))) { $candidates += $registered }
$candidates = @($candidates | Sort-Object -Unique)
$Dir = $null
if (-not $NewAgent -and $candidates.Count -eq 1) { $Dir = $candidates[0] }
elseif (-not $NewAgent -and $candidates.Count -gt 1) {
    Say pick $Names[$App]
    for ($i = 0; $i -lt $candidates.Count; $i++) { Write-Host ("  {0}. {1}" -f ($i + 1), $candidates[$i]) }
    $n = 0
    if (-not [int]::TryParse((Ask pickAsk), [ref]$n) -or $n -lt 1 -or $n -gt $candidates.Count) { Write-Host (T pickUnsure); throw 'knowslink-setup-stopped' }
    $Dir = $candidates[$n - 1]
}
$ConfigName = if ($App -eq 'claude') { 'knowslink\config\local.mcp.json' } else { 'knowslink\.mcp.json' }
$Private = $null
$previousRecord = $null
if ($Dir) {
    $recordFile = Join-Path $Dir 'knowslink-install.json'
    if (Test-Path -LiteralPath $recordFile) { $previousRecord = Get-Content -Raw -LiteralPath $recordFile -Encoding UTF8 | ConvertFrom-Json }
    # Keep the agent folder the existing host entry already uses; the config is in the install folder, not the agent folder.
    $configFile = Join-Path $Dir $ConfigName
    if (Test-Path -LiteralPath $configFile) { $Private = (Get-Content -Raw -LiteralPath $configFile -Encoding UTF8 | ConvertFrom-Json).mcpServers.knowslink.env.KNOWSLINK_AGENT_FOLDER }
    if (-not $Private) { $Private = "$Dir-private" }
    $Mode = if (Test-Path -LiteralPath (Join-Path $Private 'agent.json')) { 'update' } else { 'resume' }
} else {
    $Mode = 'new'
    $n = 1
    # A fresh name never touches an existing install or agent folder (AS-03, AC-06).
    while ((Test-Path -LiteralPath (Join-Path $Root ('{0}-{1:d2}' -f $App, $n))) -or (Test-Path -LiteralPath (Join-Path $Root ('{0}-{1:d2}-private' -f $App, $n)))) { $n++ }
    $Dir = Join-Path $Root ('{0}-{1:d2}' -f $App, $n)
    $Private = "$Dir-private"
}
$Arch = if ($env:PROCESSOR_ARCHITECTURE -eq 'ARM64') { 'arm64' } else { 'x64' }
$NodeDir = Join-Path $Root "runtime\node-v$NodeVersion-win-$Arch"
$Node = Join-Path $NodeDir 'node.exe'
$Previous = "$Dir.previous-$(Stamp)"

# N3 / U1 / AS-02: the plan and consent come before any write.
Write-Host ''
switch ($Mode) {
    'new' { Say modeNew }
    'update' { Say modeUpdate $(if ($previousRecord) { $previousRecord.version } else { '?' }) $Version }
    'resume' { Say modeResume }
}
if ($Mode -eq 'update') {
    if ($previousRecord -and $previousRecord.agent.id) { Say keepAgent ("{0} ({1}) · {2}" -f $previousRecord.agent.name, $previousRecord.agent.id, $previousRecord.agent.fingerprint) } else { Say keepAgentUnknown }
}
Say planDir $Dir
if ($Mode -eq 'update') { Say planPrivateKeep $Private } else { Say planPrivateNew $Private }
Say planNode $Node
Say planHost $Names[$App] $HostSettings $HostEntries
Say planHostKeep
if ($registered -and $registered -ne (Norm $Dir)) { Say planReplace $registered }
if ($Mode -ne 'new') { Say planPrevious $Previous }
if ($Mode -ne 'update') { Say planGoogle }
if ((Ask consent) -notmatch '^(y|yes)$') { Say declined; return }

$S.Record = [ordered]@{
    product = 'knowslink'; app = $App; version = $Version; installed = (Get-Date).ToString('o'); updated = ''
    install = $Dir; agentFolder = $Private; hostSettings = "$HostSettings ($HostEntries)"; unregister = $Unregister
    agent = [ordered]@{ id = ''; name = ''; fingerprint = '' }; steps = @()
}
if ($previousRecord) {
    $S.Record.installed = $previousRecord.installed
    if ($previousRecord.agent) { $S.Record.agent = [ordered]@{ id = "$($previousRecord.agent.id)"; name = "$($previousRecord.agent.name)"; fingerprint = "$($previousRecord.agent.fingerprint)" } }
    $S.Record.steps = @($previousRecord.steps)
}

# N4 / U3: download and verify into a private temporary folder; nothing outside it changes on failure.
[Net.ServicePointManager]::SecurityProtocol = [Net.ServicePointManager]::SecurityProtocol -bor [Net.SecurityProtocolType]::Tls12
$S.Temp = Join-Path ([IO.Path]::GetTempPath()) ("knowslink-setup-" + [guid]::NewGuid().ToString('n'))
New-Item -ItemType Directory -Path $S.Temp | Out-Null
function Get-Verified([string]$Url, [string]$Sha, [string]$Name) {
    $file = Join-Path $S.Temp $Name
    try { (New-Object Net.WebClient).DownloadFile($Url, $file) } catch { Stop-Setup 'F-03' (T F03 $Url) (T cNone) (T F03n) }
    # .NET hashing: Get-FileHash may be missing when Windows PowerShell starts with another edition's module path.
    $stream = [IO.File]::OpenRead($file)
    try { $hash = -join ([Security.Cryptography.SHA256]::Create().ComputeHash($stream) | ForEach-Object { $_.ToString('x2') }) } finally { $stream.Dispose() }
    if ($hash -ne $Sha.ToLowerInvariant()) { Stop-Setup 'F-04' (T F04) (T cTemp) (T F04n) }
    $file
}
# A resumed run of the same version skips the finished download and file steps (AS-11).
$needFiles = $Mode -ne 'resume' -or -not $previousRecord -or $previousRecord.version -ne $Version
if ($needFiles) {
    $zip = Get-Verified "$Base/knowslink-$App-plugin.zip" $Sums[$App] 'plugin.zip'
    Say stepOk (T sDownload)
}

$nodeReady = (Test-Path -LiteralPath $Node) -and ((& $Node --version) -eq "v$NodeVersion")
if (-not $nodeReady) {
    $nodeZip = Get-Verified "https://nodejs.org/dist/v$NodeVersion/node-v$NodeVersion-win-$Arch.zip" $NodeSums[$Arch] 'node.zip'
    Add-Type -AssemblyName System.IO.Compression.FileSystem
    $partial = "$NodeDir.partial-$(Stamp)"
    New-Item -ItemType Directory -Path $partial | Out-Null
    $archive = [IO.Compression.ZipFile]::OpenRead($nodeZip)
    try { [IO.Compression.ZipFileExtensions]::ExtractToFile($archive.GetEntry("node-v$NodeVersion-win-$Arch/node.exe"), (Join-Path $partial 'node.exe')) } finally { $archive.Dispose() }
    if (Test-Path -LiteralPath $NodeDir) { [IO.Directory]::Move($NodeDir, "$NodeDir.previous-$(Stamp)") }
    [IO.Directory]::Move($partial, $NodeDir)
}
Say stepOk (T sNode $NodeVersion)

# Files: a staged copy is configured and checked before it replaces anything (AS-10, F-08).
if ($needFiles) {
$Stage = "$Dir.new-$(Stamp)"
Add-Type -AssemblyName System.IO.Compression.FileSystem
[IO.Compression.ZipFile]::ExtractToDirectory($zip, $Stage)
$configPath = Join-Path $Stage $ConfigName
$config = Get-Content -Raw -LiteralPath $configPath -Encoding UTF8 | ConvertFrom-Json
$server = $config.mcpServers.knowslink
$server.command = $Node
$server.args = @(Join-Path $Dir 'knowslink\dist\plugin.js')
$serverEnv = [ordered]@{ KNOWSLINK_MODE = 'public-node'; RELAY_URL = $Relay; KNOWSLINK_AGENT_FOLDER = $Private }
# Same host contract as package_plugin.py install(): Claude receives only in the owner's switched-on session; Codex binds per session.
if ($App -eq 'claude') { $serverEnv.KNOWSLINK_AUTO_RECEIVE = 'off' } else { $serverEnv.KNOWSLINK_SESSION_HOST = 'codex' }
$server | Add-Member -NotePropertyName env -NotePropertyValue ([pscustomobject]$serverEnv) -Force
[IO.File]::WriteAllText($configPath, ($config | ConvertTo-Json -Depth 8), (New-Object Text.UTF8Encoding $false))
$check = Get-Content -Raw -LiteralPath $configPath -Encoding UTF8 | ConvertFrom-Json
if (-not (Test-Path -LiteralPath (Join-Path $Stage 'knowslink\dist\plugin.js')) -or @($check.mcpServers.knowslink.args)[0] -ne (Join-Path $Dir 'knowslink\dist\plugin.js')) {
    Stop-Setup 'F-08' (T F08 'staged files') (T cKept) (T F08n $Names[$App])
}
$S.RecordDir = $Stage
Save-Record 'files' 'staged'
if ($Mode -eq 'new') {
    [IO.Directory]::Move($Stage, $Dir)
} else {
    try {
        [IO.Directory]::Move($Dir, $Previous)
    } catch {
        $S.RecordDir = $null
        Stop-Setup 'F-08' (T F08 $_.Exception.Message) (T cKept) (T F08n $Names[$App])
    }
    try {
        [IO.Directory]::Move($Stage, $Dir)
    } catch {
        [IO.Directory]::Move($Previous, $Dir)
        $S.RecordDir = $null
        Stop-Setup 'F-08' (T F08 $_.Exception.Message) (T cKept) (T F08n $Names[$App])
    }
}
$S.RecordDir = $Dir
Save-Record 'files' 'ok'
Say stepOk (T sFiles)
}
$S.RecordDir = $Dir

# Registration: only the KnowsLink marketplace and plugin entries (Q3 A, AS-06).
$pluginJs = Join-Path $Dir 'knowslink\dist\plugin.js'
$failure = $null
if ($App -eq 'codex' -and $registered -and $registered -ne (Norm $Dir)) {
    $r = Invoke-HostCli @('plugin', 'marketplace', 'remove', $Marketplace)
    if ($r.Code -ne 0) { $failure = "codex plugin marketplace remove: $($r.Out)" }
}
$steps = if ($App -eq 'claude') { @(@('plugin', 'marketplace', 'add', $Dir, '--json'), @('plugin', 'install', $Plugin, '--scope', 'user', '--json')) } else { @(@('plugin', 'marketplace', 'add', $Dir), @('plugin', 'add', $Plugin, '--json')) }
foreach ($cli in $steps) {
    if ($failure) { break }
    $r = Invoke-HostCli $cli
    if ($r.Code -ne 0) { $failure = "$App $($cli[0..2] -join ' '): $($r.Out)" }
}
if (-not $failure) {
    if ($App -eq 'claude') {
        $entry = @(Get-HostJson @('plugin', 'list', '--json')) | Where-Object { $_.id -eq $Plugin } | Select-Object -First 1
        $mcp = if ($entry) { $entry.mcpServers.knowslink }
        $ok = $entry -and $entry.enabled -and (Norm $entry.readFromFolder) -eq (Norm (Join-Path $Dir 'knowslink'))
    } else {
        $mcp = @(Get-HostJson @('mcp', 'list', '--json')) | Where-Object { $_.name -eq 'knowslink' } | Select-Object -First 1
        $ok = $mcp -and $mcp.enabled
        if ($mcp) { $mcp = $mcp.transport }
    }
    if (-not ($ok -and @($mcp.args)[0] -eq $pluginJs -and $mcp.env.KNOWSLINK_AGENT_FOLDER -eq $Private)) { $failure = "$App plugin list: KnowsLink entry missing or points elsewhere" }
}
if ($failure) {
    $summary = ($failure -replace '\s+', ' ').Trim()
    if ($summary.Length -gt 240) { $summary = $summary.Substring(0, 240) + '…' }
    Stop-Setup 'F-05' (T F05 $Names[$App] $summary) (T cFolder $Dir $Names[$App]) (T F05n $App)
}
Save-Record 'register' 'ok'
Say stepOk (T sRegister $Names[$App])

# N5 / AS-05: Google connection only when this install has no stored key yet; an update never touches the agent folder.
if (-not (Test-Path -LiteralPath (Join-Path $Private 'agent.json'))) {
    $fresh = $false
    while ($true) {
        $outcome = $null
        $ErrorActionPreference = 'Continue'
        $connectArgs = @($pluginJs, 'setup-connect', $Relay, $Private)
        if ($fresh) { $connectArgs += '--fresh' }
        & $Node @connectArgs 2>$null | ForEach-Object {
            $evt = try { $_ | ConvertFrom-Json } catch { $null }
            if (-not $evt) { return }
            switch ($evt.state) {
                'waiting' {
                    Write-Host ''
                    if ($evt.resumed) { Say googleResume $evt.url } else { Say googleOpen $evt.url }
                    if ($OpenBrowser) { try { Start-Process $evt.url } catch { Say googleNoBrowser } }
                    Say googlePrint $evt.fingerprint
                    Say googleCompare
                    Say googleWait ([DateTime]::Parse($evt.expires).ToLocalTime().ToString('HH:mm'))
                }
                'connected' {
                    if ($evt.agent) { $S.Record.agent = [ordered]@{ id = $evt.agent; name = $evt.name; fingerprint = $evt.fingerprint } }
                }
                default { $outcome = $evt }
            }
        }
        $code = $LASTEXITCODE
        $ErrorActionPreference = 'Stop'
        if ($code -eq 0) { break }
        if ($outcome -and $outcome.state -eq 'unconfirmed') {
            Write-Host ''
            $lost = ("{0} ({1})" -f $outcome.name, $outcome.agent).Trim()
            Write-Host ("[F-07] " + (T F07 $lost)) -ForegroundColor Yellow
            Save-Record 'F-07' 'unconfirmed'
            if ((Ask f07Ask $lost) -eq 'new') { $fresh = $true; continue }
            Stop-Setup 'F-07' (T F07 $lost) (T cUnsure) (T F07n)
        }
        Stop-Setup 'F-06' (T F06) (T cNoKey) (T F06n)
    }
    Save-Record 'google' 'ok'
    $a = $S.Record.agent
    Say googleDone ("{0} ({1})" -f $a.name, $a.id) $a.fingerprint
}

Remove-Temp
Save-Record 'complete' 'ok'
# N6 / U4 / AS-08 / AS-13.
Write-Host ''
if ($Mode -eq 'update') { Write-Host (T doneUpdate) -ForegroundColor Green } else { Write-Host (T doneNew) -ForegroundColor Green }
$a = $S.Record.agent
if ($a.id) { Say doneAgent ("{0} ({1}) · {2}" -f $a.name, $a.id, $a.fingerprint) }
if ($verified) { Say testedEnv $Names[$App] $Tested[$App] }
Say restart $Names[$App]
if ($App -eq 'claude') { Say sessionClaude } else { Say sessionCodex }
if ($App -eq 'codex') {
    $toml = if (Test-Path -LiteralPath $HostSettings) { Get-Content -LiteralPath $HostSettings -Encoding UTF8 } else { @() }
    $top = @($toml | ForEach-Object -Begin { $inTable = $false } -Process { if ($_ -match '^\s*\[') { $inTable = $true }; if (-not $inTable) { $_ } })
    if ($top -match '^\s*approval_policy\s*=\s*["'']never["'']') {
        Write-Host ''
        Write-Host ("[F-09] " + (T F09)) -ForegroundColor Yellow
        Say changed (T cDone)
        Say next (T F09n $HostSettings)
        Save-Record 'F-09' 'yolo'
    }
}
Write-Host ''
Say removeTitle
Write-Host ('  ' + (T removeDir $Dir $Private))
Write-Host ('  ' + (T removeHost $Names[$App] "$HostSettings ($HostEntries)" $Unregister))
Write-Host ('  ' + (T removeKey))
Say record (Join-Path $Dir 'knowslink-install.json')
}

try {
    Invoke-Setup
} catch {
    if ("$($_.Exception.Message)" -ne 'knowslink-setup-stopped') {
        # Unexpected errors still say where it stopped; the record keeps the step, never a key or token.
        Write-Host ''
        Write-Host ("[ERROR] " + $_.Exception.Message) -ForegroundColor Yellow
        Save-Record 'error' $_.Exception.Message
        Remove-Temp
    }
    $global:LASTEXITCODE = 1
}
