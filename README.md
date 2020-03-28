# resyncd

Synchronizes macOS local files and directorines with a remotes.
Inspired by [Lsyncd](https://github.com/axkibe/lsyncd/issues/587#issuecomment-598831069).

## Installation

```
$ go get github.com/idr0id/resyncd
$ resyncd example.toml
```

## Configuration

```
[[sync]]
source = "/Users/username/projects/example"
target = "root@example.com:/srv/http/example.com"
exclude = [
  "**/.idea",
  "**/.git",
  "some-file",
]
  [sync.rsync]
  rsh = "/usr/bin/ssh -i /Users/username/.ssh/id_rsa -o StrictHostKeyChecking=no"
  acls = true
  perms = true
```
