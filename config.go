package main

type config struct {
	Syncs []configSync `toml:"sync"`
}

type path string

type configSync struct {
	Source  path
	Target  path
	Exclude []string
	Rsync   configRsync
}

type configRsync struct {
	Rsh   string `toml:"rsh"`
	ACLs  bool   `toml:"acls"`
	Perms bool   `toml:"perms"`
}

func (p path) String() string {
	return endsWithSlash(string(p))
}
