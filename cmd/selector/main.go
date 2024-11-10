package main

import (
	"fmt"
  "log"
  "os"

  "github.com/urfave/cli/v2"
)

func main() {
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
    Action: func(*cli.Context) error {
      fmt.Println("Sorting artifacts...")
      return nil
    },
  }

  if err := app.Run(os.Args); err != nil {
    log.Fatal(err)
  }
}
