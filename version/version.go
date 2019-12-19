package version

const (
	APPVersion = 1.0
)

var (
	NodeVersion     = "1.0.0"
	CompilerVersion = "1.0.0"
	CliVersion      = "1.0.0"
	LasVersion      = "1.0.0"
  	GitCommit       string
)

func init() {
	if GitCommit != "" {
		NodeVersion     += "-" + GitCommit
		CompilerVersion += "-" + GitCommit
		CliVersion      += "-" + GitCommit
		LasVersion      += "-" + GitCommit
	}
}
