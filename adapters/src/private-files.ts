// Private connection storage checks POSIX ownership/modes or Windows ACLs, and rejects links on either platform.
import { execFile } from "node:child_process";
import { lstat, mkdir } from "node:fs/promises";
import { promisify } from "node:util";

const execute = promisify(execFile);
async function windowsACL(path: string, create: boolean): Promise<void> {
  await execute(
    "powershell.exe",
    [
      "-NoProfile",
      "-NonInteractive",
      "-Command",
      `$ErrorActionPreference = 'Stop'
      $path = $env:KNOWSLINK_PRIVATE_PATH
      $sid = [System.Security.Principal.WindowsIdentity]::GetCurrent().User
      ${
        create
          ? `$acl = New-Object System.Security.AccessControl.DirectorySecurity
      $acl.SetOwner($sid)
      $acl.SetAccessRuleProtection($true, $false)
      $rule = New-Object System.Security.AccessControl.FileSystemAccessRule($sid, 'FullControl', 'ContainerInherit,ObjectInherit', 'None', 'Allow')
      $acl.AddAccessRule($rule)
      [System.IO.Directory]::SetAccessControl($path, $acl)`
          : ""
      }
      $acl = if (([System.IO.File]::GetAttributes($path) -band [System.IO.FileAttributes]::Directory) -ne 0) { [System.IO.Directory]::GetAccessControl($path) } else { [System.IO.File]::GetAccessControl($path) }
      if ($acl.GetOwner([System.Security.Principal.SecurityIdentifier]).Value -ne $sid.Value) { throw 'private owner required' }
      foreach ($rule in $acl.GetAccessRules($true, $true, [System.Security.Principal.SecurityIdentifier])) {
        if ($rule.AccessControlType -eq 'Allow' -and $rule.IdentityReference.Value -notin @($sid.Value, 'S-1-5-18', 'S-1-5-32-544')) { throw 'private ACL required' }
      }`,
    ],
    { env: { ...process.env, KNOWSLINK_PRIVATE_PATH: path }, timeout: 10000 },
  );
}

export async function privatePath(
  path: string,
  directory = false,
): Promise<void> {
  const st = await lstat(path);
  if (
    st.isSymbolicLink() ||
    (directory ? !st.isDirectory() : !st.isFile() || st.size > 8192)
  )
    throw new Error("private path required");
  if (process.platform === "win32") await windowsACL(path, false);
  else if ((st.mode & 0o077) !== 0 || st.uid !== process.getuid?.())
    throw new Error("private permissions required");
}

export async function privateDirectory(path: string): Promise<void> {
  await mkdir(path, { mode: 0o700 }); // Existing directories fail; never overwrite a working key.
  if (process.platform === "win32") await windowsACL(path, true);
  await privatePath(path, true);
}
