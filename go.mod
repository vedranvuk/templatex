module github.com/vedranvuk/templatex

go 1.14

require (
	github.com/vedranvuk/errorex v0.3.2
	github.com/vedranvuk/fs v0.0.0-20201017093503-a40e3f175c5e
)

replace (
	github.com/vedranvuk/errorex => ../errorex
	github.com/vedranvuk/fs => ../fs
)
