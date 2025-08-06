package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/darthbanana13/artifact-selector/pkg/github/builder"
	fetcherconcur "github.com/darthbanana13/artifact-selector/pkg/github/concur"
	"github.com/darthbanana13/artifact-selector/pkg/log"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/transmute"
	archfilter "github.com/darthbanana13/artifact-selector/pkg/filter/concur/arch"
	archhandleerr "github.com/darthbanana13/artifact-selector/pkg/filter/concur/arch/decorator/handleerr"
	archlog "github.com/darthbanana13/artifact-selector/pkg/filter/concur/arch/decorator/log"
	extfilter "github.com/darthbanana13/artifact-selector/pkg/filter/concur/ext"
	exthandleerr "github.com/darthbanana13/artifact-selector/pkg/filter/concur/ext/decorator/handleerr"
	extlog "github.com/darthbanana13/artifact-selector/pkg/filter/concur/ext/decorator/log"
	osfilter "github.com/darthbanana13/artifact-selector/pkg/filter/concur/os"
	oshandleerr "github.com/darthbanana13/artifact-selector/pkg/filter/concur/os/decorator/handleerr"
	oslog "github.com/darthbanana13/artifact-selector/pkg/filter/concur/os/decorator/log"

	"github.com/darthbanana13/artifact-selector/pkg/filter/extractor"
	extext "github.com/darthbanana13/artifact-selector/pkg/filter/extractor/ext"
	"github.com/darthbanana13/artifact-selector/pkg/filter/extractor/withinsize"
	"github.com/darthbanana13/artifact-selector/pkg/filter/extractor/max"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"

	// "github.com/darthbanana13/artifact-selector/pkg/filter/linuxbindiff"

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
				Value:   "neovim/neovim",
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
			fetcher, err := builder.NewGihubFetcher().
				WithLogger(&logger).
				WithRetry(3).
				WithLogin(strings.TrimSpace(cmd.String("token"))).
				Build()

			if err != nil {
				return err
			}

			artifacts, info, err := fetcherconcur.FetchArtifacts(fetcher, cmd.String("github"))
			if err != nil {
				return err
			}

			input := transmute.ToFilter(artifacts)

			// TODO: Add 4 more filters:
			//  xz (deb), for debs compressed with zst (lsd-rs/lsd), or a generic regex that can be applied multiple times
			//  musl/gnu (for musl vs gnu libc)
			//  common names (like checksum, checksums, hashes, man, only) (mikefarah/yq)
			newArchFilter := funcdecorator.DecorateFunction(archfilter.NewArchFilter,
				archhandleerr.HandleErrConstructorDecorator(),
				archlog.LogConstructorDecorator(&logger),
			)
			archF, err := newArchFilter(cmd.String("arch"))
			if err != nil {
				return err
			}
			newOSFilter := funcdecorator.DecorateFunction(osfilter.NewOSFilter,
				oshandleerr.HandleErrConstructorDecorator(),
				oslog.LogConstructorDecorator(&logger),
			)
			osF, err := newOSFilter(cmd.String("os"))
			if err != nil {
				return err
			}
			newExtFilter := funcdecorator.DecorateFunction(extfilter.NewExtFilter,
				exthandleerr.HandleErrConstructorDecorator(),
				extlog.LogConstructorDecorator(&logger),
			)
			extList := strings.Split(cmd.String("extension"), ",")
			extF, err := newExtFilter(extList)
			if err != nil {
				return err
			}

			// binsize := make(chan uint64, len(info.Artifacts))
			binsize := make(chan uint64)
			//TODO: Test more if including appimage is a good idea
			extr, err := extext.NewExt([]string{extfilter.LINUXBINARY, "appimage"}, binsize)
			if err != nil {
				return err
			}

			var (
				archStrategy				concur.FilterFunc = archF.FilterArtifact
				osStrategy 					concur.FilterFunc = osF.FilterArtifact
				extStrategy 				concur.FilterFunc = extF.FilterArtifact
			)
			extractorStrategy, err := extractor.NewExtractor(extr)

			var withinSizeOnce sync.Once
			var withinSizeF *withinsize.WithinSize
			var withinSizeStrategy concur.FilterFunc = func(a filter.Artifact) (filter.Artifact, bool) {
				withinSizeOnce.Do(func() {
					withinSizeF, err = withinsize.NewWithinSize(max.Find(binsize), 20, []string{extfilter.LINUXBINARY})
				})
				if err != nil {
					//TODO: How do we return an error here? I guess we bypass the filter?
					return a, true
				}
				return withinSizeF.FilterArtifact(a)
			}

			// output := linuxbindiff.Filter(extStrategy.Filter(osStrategy.Filter(archStrategy.Filter(input))))
			output := withinSizeStrategy.Filter(osStrategy.Filter(archStrategy.Filter(extractorStrategy.Extract(extStrategy.Filter(input)))))

			artifactss := make([]filter.Artifact, 0)
			for artifact := range output {
				artifactss = append(artifactss, artifact)
			}
			releases := filter.ReleasesInfo{
				Version:    info.Version,
				PreRelease: info.PreRelease,
				Draft:      info.Draft,
				Artifacts:  artifactss,
			}
			minfo, err := json.Marshal(releases)

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
