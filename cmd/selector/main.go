package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/darthbanana13/artifact-selector/pkg/ghandledecorate"
	"github.com/darthbanana13/artifact-selector/pkg/github"
	"github.com/darthbanana13/artifact-selector/pkg/glogdecorate"
	"github.com/darthbanana13/artifact-selector/pkg/glogindecorate"
	"github.com/darthbanana13/artifact-selector/pkg/gretryclient"
	"github.com/darthbanana13/artifact-selector/pkg/log"

	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	archfilter "github.com/darthbanana13/artifact-selector/pkg/filter/arch"
	archhandleerror "github.com/darthbanana13/artifact-selector/pkg/filter/arch/decorator/handleerror"
	archlog "github.com/darthbanana13/artifact-selector/pkg/filter/arch/decorator/log"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur"
	// extfilter "github.com/darthbanana13/artifact-selector/pkg/filter/ext"
	osfilter "github.com/darthbanana13/artifact-selector/pkg/filter/os"

	"github.com/urfave/cli/v3"
)

func main() {
	// TODO: Handle different log-levels
	logger := log.InitLog("dev")
	// TODO: Default values should be based on current OS
	cmd := &cli.Command{
		Name:                  "Artifact finder",
		Usage:                 "Use this utility to find the best artifact according to your specifications",
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "github",
				Aliases: []string{"g"},
				Value:   "BurntSushi/ripgrep",
				Usage:   "Specify the 'user/project_name' to directly look up github projects artifacts",
				// Required: true,
			},
			&cli.StringFlag{
				Name:    "extension",
				Aliases: []string{"e"},
				Value:   "deb,,appimage,tar.zst,tbz,tar.gz,tar.xz",
				Usage:   "List the extension preference in a comma separated list. E.g. 'deb,appimage,LINUXBINARY'",
			},
			&cli.StringFlag{
				Name:    "arch",
				Aliases: []string{"a"},
				Value:   "x86_64",
				Usage:   "Specify the target architecture for the binary. E.g. amd64, arm64, x86",
			},
			&cli.StringFlag{
				Name:    "os",
				Aliases: []string{"o"},
				Value:   "ubuntu",
				Usage:   "Specify the target OS/Distro. E.g. ubuntu, linux, macos",
			},
			&cli.StringFlag{
				Name:    "token",
				Aliases: []string{"t"},
				Usage:   "Github token",
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

			input := make(chan github.Artifact)
			go func() {
				defer close(input)

				for _, artifact := range info.Artifacts {
					input <- artifact
				}
			}()

			// TODO: Add 4 more filters:
			//  xz (deb), for debs compressed with zst (lsd-rs/lsd), or a generic regex that can be applied multiple times
			//  musl/gnu (for musl vs gnu libc)
			//  size difference (for example, if a file is an order of magnitude smaller than other artifacts, it's probably a text file) (mikefarah/yq)
			//  common names (like checksum, checksums, hashes, man, only) (mikefarah/yq)
			newArchFilter := funcdecorator.DecorateFunction(archfilter.NewArchFilter,
				archhandleerror.HandleErrorConstructorDecorator(),
				archlog.LogConstructorDecorator(&logger),
			)
			archF, err := newArchFilter(cmd.String("arch"))
			if err != nil {
				return err
			}
			osF, err := osfilter.NewOSFilter(cmd.String("os"))
			if err != nil {
				return err
			}
			// extList := strings.Split(cmd.String("extension"), ",")
			// extF, err := extfilter.NewOSFilter(extList)
			// if err != nil {
			//   return err
			// }

			var archStrategy, osStrategy concur.FilterFunc
			archStrategy = archF.FilterArtifact
			osStrategy = osF.FilterArtifact

			output := osStrategy.Filter(archStrategy.Filter(input))

			artifacts := make([]github.Artifact, 0)
			for artifact := range output {
				artifacts = append(artifacts, artifact)
			}
			info.Artifacts = artifacts
			minfo, err := json.Marshal(info)

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
