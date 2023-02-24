package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var (
	gitVersion = ""
	gitCommit  = "" // sha1 from git, output of $(git rev-parse HEAD)
	buildDate  = "" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	goVersion  = runtime.Version()
	platform   = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return the bookimporter version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bookimporter version info")
		fmt.Printf("  - Version: %v\n", gitVersion)
		fmt.Printf("  - Commit: %v\n", gitCommit)
		fmt.Printf("  - Build Date: %v\n", buildDate)
		fmt.Printf("  - Go Version: %v\n", goVersion)
		fmt.Printf("  - Platform: %s\n", platform)
	},
}
