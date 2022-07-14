module github.com/idr0id/resyncd

go 1.23

require (
	github.com/BurntSushi/toml v1.4.0
	github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/fsnotify/fsevents v0.2.0
	github.com/kovetskiy/lorg v1.2.1-0.20240830111423-ba4fe8b6f7c4
	github.com/reconquest/cog v0.0.0-20240830113510-c7ba12d0beeb
	github.com/reconquest/karma-go v1.5.0
	github.com/stretchr/testify v1.9.0
	github.com/zloylos/grsync v1.7.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/lmittmann/tint v1.0.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/zazab/zhash v0.0.0-20221031090444-2b0d50417446 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/zloylos/grsync v0.0.0-20200204095520-71a00a7141be => github.com/idr0id/grsync v0.0.0-20200328163231-f3086de9b59d
