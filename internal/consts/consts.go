package consts

import "runtime"

const (
	ProjName = "md2img"
	Version  = "0.3.5"
)

var (
	GitTag      = ""
	GitCommit   = ""
	BuildTime   = ""
	Description = "Version: " + Version +
		"\nGo Version: " + runtime.Version() +
		"\nGit Tag: " + GitTag +
		"\nGit Commit: " + GitCommit +
		"\nBuild Time: " + BuildTime
)
