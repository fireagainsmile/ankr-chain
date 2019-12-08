package version

var (
	APPVersion      = "1.0.0"
	CompilerVersion = "1.0.0"
	GitCommit  string
)

func init() {
	if GitCommit != "" {
		APPVersion += "-" + GitCommit
		CompilerVersion += "-" + GitCommit
	}
}