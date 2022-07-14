package main

const (
	rsyncDefaultTimeoutSeconds = 10
)

type config struct {
	Syncs []configSync `toml:"sync"`
}

type configSync struct {
	Source  string      `toml:"source"`
	Target  string      `toml:"target"`
	Exclude []string    `toml:"exclude"`
	Rsync   configRsync `toml:"rsync"`
}

type configRsync struct {
	Rsh                   string `toml:"rsh"`
	ACLs                  bool   `toml:"acls"`
	Perms                 bool   `toml:"perms"`
	TimeoutSeconds        int    `toml:"timeout_seconds"`
	ConnectTimeoutSeconds int    `toml:"connect_timeout_seconds"`
}

func (c configRsync) getTimeoutSeconds() int {
	if c.TimeoutSeconds > 0 {
		return c.TimeoutSeconds
	}
	return rsyncDefaultTimeoutSeconds
}

func (c configRsync) getConnectTimeoutSeconds() int {
	if c.ConnectTimeoutSeconds > 0 {
		return c.ConnectTimeoutSeconds
	}
	return rsyncDefaultTimeoutSeconds
}
