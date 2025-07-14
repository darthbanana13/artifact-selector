package main

import (
  "context"
	"encoding/json"
	"fmt"
	"strings"
	"os"

	"github.com/darthbanana13/artifact-selector/pkg/ghandledecorate"
	"github.com/darthbanana13/artifact-selector/pkg/github"
	"github.com/darthbanana13/artifact-selector/pkg/glogdecorate"
	"github.com/darthbanana13/artifact-selector/pkg/glogindecorate"
	"github.com/darthbanana13/artifact-selector/pkg/gretryclient"
	"github.com/darthbanana13/artifact-selector/pkg/log"

	archfilter "github.com/darthbanana13/artifact-selector/pkg/filter/arch"
	extfilter "github.com/darthbanana13/artifact-selector/pkg/filter/ext"
	osfilter "github.com/darthbanana13/artifact-selector/pkg/filter/os"

	"github.com/urfave/cli/v3"
)

func main() {
  // TODO: Handle different log-levels
  logger := log.InitLog("dev")
  // TODO: Default values should be based on current OS
  cmd := &cli.Command{
    Name: "Artifact finder",
    Usage: "Use this utility to find the best artifact according to your specifications",
    EnableShellCompletion: true,
    Flags: []cli.Flag{
      &cli.StringFlag{
        Name: "github",
        Aliases: []string{"g"},
        // Value: "neovim/neovim",
        Usage: "Specify the 'user/project_name' to directly look up github projects artifacts",
        Required: true,
      },
      &cli.StringFlag{
        Name: "extension",
        Aliases: []string{"e"},
        // Value: "deb,,appimage,tar.xz,tar.gz",
        Usage: "List the extension preference in a comma separated list. E.g. 'deb,appimage,LINUXBINARY'",
      },
      &cli.StringFlag{
        Name: "arch",
        Aliases: []string{"a"},
        // Value: "x86_64",
        Usage: "Specify the target architecture for the binary. E.g. amd64, arm64, x86",
      },
      &cli.StringFlag{
        Name: "os",
        Aliases: []string{"o"},
        // Value: "ubuntu",
        Usage: "Specify the target OS/Distro. E.g. ubuntu, linux, macos",
      },
      &cli.StringFlag{
        Name: "token",
        Aliases: []string{"t"},
        Usage: "Github token",
        //TODO: This should probably point to XDG_CONFIG_HOME 1st before trying $HOME
        Sources: cli.Files(os.Getenv("HOME") + "/.config/artifact-selector/github_token"),
      },
    },
    //TODO: Clean up this funcion
    Action: func(ctx context.Context, cmd *cli.Command) error {
      logAdapter := gretryclient.NewLeveledLoggerAdapter(&logger)
      // fetcher := github.NewDefaultHttpFetcher()
      fetcher := github.NewHttpFetcher(gretryclient.NewRetryClient(3, logAdapter))
      fetcherL, err := glogindecorate.NewLoginDecorator(fetcher, strings.TrimSpace(cmd.String("token")))
      if err != nil {
        return err
      }
      fetcherE := ghandledecorate.NewHandleErrorDecorator(fetcherL)
      fetcherD := glogdecorate.NewLogFetcherDecorator(&logger, fetcherE)
      info, err := fetcherD.FetchArtifacts(cmd.String("github"))
      if err != nil {
        return err
      }

      archF, err := archfilter.NewArchFilter(cmd.String("arch"))
      if err != nil {
        return err
      }
      osF, err := osfilter.NewOSFilter(cmd.String("os"))
      if err != nil {
        return err
      }
      extList := strings.Split(cmd.String("extension"), ",")
      extF, err := extfilter.NewOSFilter(extList)
      if err != nil {
        return err
      }
      osF.SetNext(extF)
      archF.SetNext(osF)
      filteredInfo := archF.Filter(info) 
      minfo, err := json.Marshal(filteredInfo)
      if err != nil {
        return err
      }
      fmt.Println(string(minfo))

      return err
    },
  }

  if err := cmd.Run(context.Background(), os.Args); err != nil {
    logger.Panic(err.Error())
  }
}
