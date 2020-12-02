module github.com/vedranvuk/templatex

go 1.14

require (
	github.com/vedranvuk/errorex v0.3.2
	github.com/vedranvuk/fs v0.0.0-20201117213700-930f9ecd25d2
	github.com/vedranvuk/fsex v0.0.2
)

replace (
	github.com/vedranvuk/errorex => ../errorex
	github.com/vedranvuk/fs => ../fs
)
