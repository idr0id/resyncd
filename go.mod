module github.com/idr0id/resyncd

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/fsnotify/fsevents v0.1.1
	github.com/kovetskiy/lorg v0.0.0-20200107130803-9a7136a95634
	github.com/reconquest/cog v0.0.0-20191208202052-266c2467b936
	github.com/reconquest/karma-go v0.0.0-20200225150538-5e903a7b8526
	github.com/stretchr/testify v1.5.1 // indirect
	github.com/zazab/zhash v0.0.0-20170403032415-ad45b89afe7a // indirect
	github.com/zloylos/grsync v0.0.0-20200204095520-71a00a7141be
)

replace github.com/zloylos/grsync v0.0.0-20200204095520-71a00a7141be => github.com/idr0id/grsync v0.0.0-20200328163231-f3086de9b59d
