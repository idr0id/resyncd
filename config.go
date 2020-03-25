package main

type config struct {
	Syncs []configSync `toml:"sync"`
}

type configSync struct {
	Source  string
	Target  string
	Exclude []string
	Rsync   configRsync
}

type configRsync struct {
	Rsh   string `toml:"rsh"`
	ACLs  bool   `toml:"acls"`
	Perms bool   `toml:"perms"`
}
