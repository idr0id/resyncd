module github.com/idr0id/resyncd

go 1.13

require (
	github.com/BurntSushi/toml v1.1.0
	github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/fsnotify/fsevents v0.1.1
	github.com/kovetskiy/lorg v0.0.0-20200107130803-9a7136a95634
	github.com/reconquest/cog v0.0.0-20210820140837-c5c4e8f49c65
	github.com/reconquest/karma-go v0.0.0-20211029072727-6027c6225ce4
	github.com/stretchr/testify v1.7.0
	github.com/zloylos/grsync v1.6.1
)

replace github.com/zloylos/grsync v0.0.0-20200204095520-71a00a7141be => github.com/idr0id/grsync v0.0.0-20200328163231-f3086de9b59d
