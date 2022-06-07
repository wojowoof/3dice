module wojones.com/src/3dice

go 1.18

replace (
	wojones.com/src/dicegame => ./dicegame
	wojones.com/src/dicescore => ./dicescore
	wojones.com/src/diceturn => ./diceturn
)

require wojones.com/src/dicegame v0.0.0-00010101000000-000000000000

require (
	wojones.com/src/dicescore v0.0.0-00010101000000-000000000000
	wojones.com/src/diceturn v0.0.0-00010101000000-000000000000
)

require golang.org/x/exp v0.0.0-20220602145555-4a0574d9293f // indirect
