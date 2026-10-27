// Private connection storage checks POSIX ownership/modes or Windows ACLs, and rejects links on either platform.
import { execFile } from "node:child_process";
import { chmod, lstat, mkdir, readdir } from "node:fs/promises";
import { promisify } from "node:util";

// The folder cannot be made or proven private; the connection reports reason private_folder (D02 R-02).
export class PrivateFolderError extends Error {
  constructor(readonly path: string) {
    super("private folder required");
  }
}
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
      # Writing even the same owner needs WRITE_OWNER, which non-system drives (D:) do not grant; write it only to change it.
      if ([System.IO.Directory]::GetAccessControl($path).GetOwner([System.Security.Principal.SecurityIdentifier]).Value -ne $sid.Value) { $acl.SetOwner($sid) }
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
  ).catch(() => {
    throw new PrivateFolderError(path);
  });
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
    throw new PrivateFolderError(path);
  if (process.platform === "win32") await windowsACL(path, false);
  else if ((st.mode & 0o077) !== 0 || st.uid !== process.getuid?.())
    throw new PrivateFolderError(path);
}

export async function privateDirectory(path: string): Promise<void> {
  try {
    await mkdir(path, { mode: 0o700 });
  } catch (e) {
    // An empty existing folder is used once it is made private (R-01); a folder with any content fails, so a working key is never overwritten.
    if (
      (e as NodeJS.ErrnoException).code !== "EEXIST" ||
      (await lstat(path)).isSymbolicLink() ||
      (await readdir(path)).length > 0
    )
      throw e;
    if (process.platform !== "win32")
      await chmod(path, 0o700).catch(() => {
        throw new PrivateFolderError(path);
      });
  }
  if (process.platform === "win32") await windowsACL(path, true);
  await privatePath(path, true);
}
