package main

import (
  "encoding/json"
  "fmt"
  "os"

  "github.com/darthbanana13/artifact-selector/pkg/log"
  "github.com/darthbanana13/artifact-selector/pkg/github"

  "github.com/urfave/cli/v2"
)

func main() {
  // TODO: Handle different log-level
  logger := log.InitLog("dev")
  app := &cli.App{
    Name: "Artifact finder",
    Usage: "Use this utility to find the best artifact according to your specifications",
    Flags: []cli.Flag{
      &cli.StringFlag{
        Name: "github",
        Aliases: []string{"g"},
        Usage: "Specify the 'user/project_name' to directly look up github projects artifacts",
      },
      &cli.StringFlag{
        Name: "extension",
        Aliases: []string{"e"},
        Usage: "List the extension preference in a comma separated list. E.g. 'deb,appimage,LINUXBINARY'",
      },
      &cli.StringFlag{
        Name: "arch",
        Aliases: []string{"a"},
        Usage: "Specify the taget architecture for the binary. E.g. amd64, arm64, x86",
      },
    },
    Action: func(ctx *cli.Context) error {
      info, err := github.FetchArtifacts(ctx.String("github"))
      minfo, _ := json.Marshal(info)
      fmt.Println(string(minfo))
      return err
    },
  }

  if err := app.Run(os.Args); err != nil {
    logger.Fatal(err.Error())
  }
}
