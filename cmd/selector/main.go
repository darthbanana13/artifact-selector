package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/darthbanana13/artifact-selector/pkg/github/builder"
	fetcherconcur "github.com/darthbanana13/artifact-selector/pkg/github/concur"
	"github.com/darthbanana13/artifact-selector/pkg/log"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
	archbuilder "github.com/darthbanana13/artifact-selector/pkg/filter/concur/arch/builder"
	extfilter "github.com/darthbanana13/artifact-selector/pkg/filter/concur/ext"
	extbuilder "github.com/darthbanana13/artifact-selector/pkg/filter/concur/ext/builder"
	extmetadatafilter "github.com/darthbanana13/artifact-selector/pkg/filter/concur/ext/metadata"
	osbuilder "github.com/darthbanana13/artifact-selector/pkg/filter/concur/os/builder"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/transmute"
	"github.com/darthbanana13/artifact-selector/pkg/filter/pipeline"
	"github.com/darthbanana13/artifact-selector/pkg/filter/tee"

	withinsizeBuilder "github.com/darthbanana13/artifact-selector/pkg/filter/concur/extswithinsize/builder"

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
			archStrategy, err := archbuilder.NewArchBuilder().
				WithLogger(&logger).
				WithArch(cmd.String("arch")).
				Build()
			if err != nil {
				return err
			}

			osStrategy, err := osbuilder.
				NewOSBuilder().
				WithLogger(&logger).
				WithOS(cmd.String("os")).
				Build()
			if err != nil {
				return err
			}

			extBuilder := extbuilder.
				NewExtFilterBuilder().
				WithLogger(&logger)

			extStrategy, err := extBuilder.
				WithExts(strings.Split(cmd.String("extension"), ",")).
				WithLoggerName("Extension Filter").
				WithConstructor(extfilter.NewExt).
				Build()
			if err != nil {
				return err
			}

			binaryStrategy, err := extBuilder. //TODO: Test more if including appimage is a good idea
								WithExts([]string{extfilter.LINUXBINARY, "appimage"}).
								WithLoggerName("Binary Extractor").
								WithConstructor(extmetadatafilter.NewExt).
								Build()
			if err != nil {
				return err
			}

			pipe := pipeline.Process(input, extStrategy)
			pipe, extractor := tee.Tee(pipe)
			extractor = pipeline.Process(extractor, binaryStrategy)

			withinSizeStrategy, err := withinsizeBuilder.NewWithinSizeFilterBuilder().
				WithLogger(&logger).
				WithExts([]string{extfilter.LINUXBINARY}).
				WithPercentage(20).
				WithChannelMax(extractor).
				Build()

			pipe = pipeline.Process(pipe, archStrategy, osStrategy, withinSizeStrategy)

			artifactSlice := make([]filter.Artifact, 0)
			for artifact := range pipe {
				artifactSlice = append(artifactSlice, artifact)
			}
			releases := filter.ReleasesInfo{
				Version:    info.Version,
				PreRelease: info.PreRelease,
				Draft:      info.Draft,
				Artifacts:  artifactSlice,
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
