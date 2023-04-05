module wojones.com/src/3dice

go 1.18

replace (
	wojones.com/src/dicegame => ./dicegame
	wojones.com/src/dicescore => ./dicescore
	wojones.com/src/diceturn => ./diceturn
)

require wojones.com/src/dicegame v0.0.0-00010101000000-000000000000

require (
	golang.org/x/exp v0.0.0-20220602145555-4a0574d9293f
	wojones.com/src/dicescore v0.0.0-00010101000000-000000000000
)

require wojones.com/src/diceturn v0.0.0-00010101000000-000000000000

require (
	github.com/cosiner/argv v0.1.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/derekparker/trie v0.0.0-20200317170641-1fdf38b7b0e9 // indirect
	github.com/flosch/pongo2 v0.0.0-20200913210552-0d938eb266f3 // indirect
	github.com/go-chi/chi v1.5.4
	github.com/go-chi/chi/v5 v5.0.8
	github.com/go-delve/delve v1.20.1 // indirect
	github.com/go-delve/liner v1.2.3-0.20220127212407-d32d89dd2a5d // indirect
	github.com/google/go-dap v0.6.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/mattn/go-colorable v0.0.9 // indirect
	github.com/mattn/go-isatty v0.0.3 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/spf13/cobra v1.1.3 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.starlark.net v0.0.0-20220816155156-cfacd8902214 // indirect
	golang.org/x/arch v0.0.0-20190927153633-4e8777c89be4 // indirect
	golang.org/x/sys v0.0.0-20220908164124-27713097b956 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
