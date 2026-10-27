# Same two lines as Get-Registered (scripts/autosetup/knowslink-setup.ps1:244-245) on a synthetic `claude plugin list --json` payload.
function Get-HostJson($t) { $t | ConvertFrom-Json }
$Plugin = 'knowslink@knowslink-plugins'
foreach ($case in @(
  @{ name = '1 plugin (knowslink only)'; json = '[{"id":"knowslink@knowslink-plugins","readFromFolder":"C:/x/claude-01/knowslink"}]' },
  @{ name = '2 plugins (knowslink + another)'; json = '[{"id":"knowslink@knowslink-plugins","readFromFolder":"C:/x/claude-01/knowslink"},{"id":"other@inline"}]' })) {
  try {
    $entry = @(Get-HostJson $case.json) | Where-Object { $_.id -eq $Plugin } | Select-Object -First 1
    "{0}: entry is array={1}; parent={2}" -f $case.name, ($entry -is [array]), (Split-Path -Parent $entry.readFromFolder)
  } catch { "{0}: ERROR {1}" -f $case.name, $_.Exception.Message }
}
