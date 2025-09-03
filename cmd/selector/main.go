package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/darthbanana13/artifact-selector/internal/fetcher"
	fetcherconcur "github.com/darthbanana13/artifact-selector/internal/fetcher/concur"
	"github.com/darthbanana13/artifact-selector/internal/fetcher/github/builder"
	"github.com/darthbanana13/artifact-selector/internal/log"

	regexcli "github.com/darthbanana13/artifact-selector/internal/cli/regex"
	"github.com/darthbanana13/artifact-selector/internal/filter"
	archbuilder "github.com/darthbanana13/artifact-selector/internal/filter/concur/arch/builder"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/convert"
	extfilter "github.com/darthbanana13/artifact-selector/internal/filter/concur/ext"
	extbuilder "github.com/darthbanana13/artifact-selector/internal/filter/concur/ext/builder"
	contenttypebuilder "github.com/darthbanana13/artifact-selector/internal/filter/concur/ext/metadata/contenttype/builder"
	extmetadatafilter "github.com/darthbanana13/artifact-selector/internal/filter/concur/ext/metadata/ext"
	withinsizebuilder "github.com/darthbanana13/artifact-selector/internal/filter/concur/extswithinsize/builder"
	muslbuilder "github.com/darthbanana13/artifact-selector/internal/filter/concur/musl/builder"
	osbuilder "github.com/darthbanana13/artifact-selector/internal/filter/concur/os/builder"
	osverbuilder "github.com/darthbanana13/artifact-selector/internal/filter/concur/osver/builder"
	"github.com/darthbanana13/artifact-selector/internal/filter/pipeline"
	"github.com/darthbanana13/artifact-selector/internal/filter/tee"

	altsrc "github.com/urfave/cli-altsrc/v3"
	jsonconfig "github.com/urfave/cli-altsrc/v3/json"
	"github.com/urfave/cli/v3"
)

func main() {
	var verbosityCount int
	var logger log.ILogger
	// TODO: Default values should be based on current OS
	cmd := &cli.Command{
		Name:                   "Artifact finder",
		Usage:                  "Use this utility to find the best artifact according to your specifications",
		EnableShellCompletion:  true,
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "github",
				Aliases: []string{"g"},
				Value:   "neovim/neovim",
				Usage:   "Specify the 'user/project_name' to directly look up github projects artifacts",
				// Required: true,
			},
			&cli.StringFlag{
				Name:    "release",
				Aliases: []string{"r"},
				Value:   "latest",
				Usage:   "What release version of the application artifact to fetch. E.g. latest, v1.2",
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
				Name:    "os-version",
				Aliases: []string{"O"},
				Value:   "24.04",
				Usage:   "The version of the distro/os that you are targeting",
			},
			&cli.BoolFlag{
				Name:    "musl",
				Aliases: []string{"m"},
				Usage:   "Exclude musl artifacts",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Specify the verbosity level, default is none/only errors",
				Config: cli.BoolConfig{
					Count: &verbosityCount,
				},
			},
			&cli.StringSliceFlag{
				Name:    "regex",
				Aliases: []string{"X"},
				Usage:   "Additional regex(es) filter to apply",
			},
			&cli.StringSliceFlag{
				Name:    "regex-meta",
				Aliases: []string{"M"},
				Usage:   "Give a name to the regex match metadata key",
			},
			&cli.StringSliceFlag{
				Name:    "regex-lower",
				Aliases: []string{"L"},
				Usage: `Should the regex(es) apply to a string that has been lowercased, possible values "yes", "no", "y", "n".
Default: "no"`,
			},
			&cli.StringSliceFlag{
				Name:    "regex-filter",
				Aliases: []string{"F"},
				Usage: `Should the regex(es) exclude the matched values.
Possible values:
	"yes", "y" - filters the values that don't match the regex
	"no", "n" - only adds metadata if there is a match
	"exclude", "e" - excludes the values that match the regex
Default: "no"`,
			},
			//TODO: The token should be optional
			&cli.StringFlag{
				Name:    "token",
				Aliases: []string{"t"},
				Usage:   "Github token, either classic or fine-grained with 'repo' scope",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("GITHUB_TOKEN"),
					jsonconfig.JSON("github_token", altsrc.StringSourcer(os.Getenv("XDG_CONFIG_HOME")+"/artifact-selector/config.json")),
					jsonconfig.JSON("github_token", altsrc.StringSourcer(os.Getenv("HOME")+"/.config/artifact-selector/config.json")),
				),
			},
		},
		//TODO: Clean up this function
		Action: func(ctx context.Context, cmd *cli.Command) error {
			var logLevel string
			switch verbosityCount {
			case 0:
				logLevel = log.Prod
			case 1:
				logLevel = log.Verbose
			default:
				logLevel = log.VeryVerbose
			}
			logger = log.InitLog(logLevel)

			fetcherStrategy, err := builder.NewGithubFetcher().
				WithLogger(logger).
				WithRetry(3).
				WithLogin(strings.TrimSpace(cmd.String("token"))).
				Build()

			if err != nil {
				return err
			}

			f := fetcher.NewFetcher(fetcherStrategy)

			artifacts, info, err := fetcherconcur.FetchArtifacts(f, cmd.String("github"), cmd.String("release"))
			if err != nil {
				return err
			}

			input := convert.ToFilter(artifacts)

			// TODO: Maybe add 1 more filter: common names (like checksum, checksums, hashes, man, only) (mikefarah/yq)
			//	Could be useful for non-github sources where maybe we don't know the size and/or content type
			extBuilder := extbuilder.
				NewExtFilterBuilder().
				WithLogger(logger)

			extStrategy, err := extBuilder.
				WithExts(strings.Split(cmd.String("extension"), ",")).
				WithLoggerName("Extension Filter").
				WithConstructor(extfilter.NewExt).
				Build()
			if err != nil {
				return err
			}

			compressedExtensions := []string{"deb", "tar.gz", "zip", "tar.xz", "tar.bz2", "tbz", "tar.zst", "rpm"}
			compressedStrategy, err := extBuilder.
				WithExts(compressedExtensions).
				WithLoggerName("Compressed Extractor").
				WithConstructor(extmetadatafilter.NewExt).
				Build()
			if err != nil {
				return err
			}

			//TODO: Test more if including appimage is a good idea
			binaryExtensions := []string{extfilter.LinuxBinary, "appimage"}
			//NOTE: In case we're not sure if the filtered artifacts are actually a binary, and there is no actual binary in the
			//	artifacts list, we can add compressed extensions for calculating the max. This way, if there is something silly
			//	like a txt file renamed to a random extension, because it's a lot smaller than the artifacts we know pretty sure
			//	are compressed, then it's most likely not the binary we're looking for
			binaryExtensions = append(binaryExtensions, compressedExtensions...)
			binaryStrategy, err := extBuilder.
				WithExts(binaryExtensions).
				WithLoggerName("Binary Extractor").
				WithConstructor(extmetadatafilter.NewExt).
				Build()
			if err != nil {
				return err
			}

			pipe := pipeline.Process(input, extStrategy)
			pipe, extractor := tee.Tee(pipe)
			binExtractor, compressedExtractor := tee.Tee(extractor)
			compressedExtractor = pipeline.Process(compressedExtractor, compressedStrategy)
			binExtractor = pipeline.Process(binExtractor, binaryStrategy)

			archStrategy, err := archbuilder.NewArchBuilder().
				WithLogger(logger).
				WithArch(cmd.String("arch")).
				Build()
			if err != nil {
				return err
			}
			pipe = pipeline.Process(pipe, archStrategy)

			osStrategy, err := osbuilder.
				NewOSBuilder().
				WithLogger(logger).
				WithOS(cmd.String("os")).
				Build()
			if err != nil {
				return err
			}
			pipe = pipeline.Process(pipe, osStrategy)

			contentTypeStrategy, err := contenttypebuilder.
				NewContentTypeFilterBuilder().
				WithLogger(logger).
				Build()
			if err != nil {
				return err
			}
			pipe = pipeline.Process(pipe, contentTypeStrategy)

			muslStrategy, err := muslbuilder.
				NewMuslFilterBuilder().
				WithLogger(logger).
				WithFilter(cmd.Bool("musl")).
				Build()
			if err != nil {
				return err
			}
			pipe = pipeline.Process(pipe, muslStrategy)

			osVerStrategy, err := osverbuilder.
				NewOSVerBuilder().
				WithLogger(logger).
				Build()
			if err != nil {
				return err
			}
			pipe = pipeline.Process(pipe, osVerStrategy)

			regexFilters := make([]filter.IFilter, len(cmd.StringSlice("regex")))
			regexStrategies, err := regexcli.ProcessRegexParams(
				cmd.StringSlice("regex"),
				cmd.StringSlice("regex-lower"),
				cmd.StringSlice("regex-filter"),
				cmd.StringSlice("regex-meta"),
				logger,
			)
			if err != nil {
				return err
			}
			for i, regexStrategy := range regexStrategies {
				regexFilters[i] = regexStrategy
			}
			pipe = pipeline.Process(pipe, regexFilters...)

			withinSizeBuilder := withinsizebuilder.NewWithinSizeFilterBuilder().
				WithLogger(logger).
				WithPercentage(20)

			compressedWithinSizeStrategy, err := withinSizeBuilder.
				WithExts(compressedExtensions).
				WithLoggerName("Compressed Filter").
				WithChannelMax(compressedExtractor).
				Build()
			if err != nil {
				return err
			}
			pipe = pipeline.Process(pipe, compressedWithinSizeStrategy)

			binWithinSizeStrategy, err := withinSizeBuilder.
				WithExts([]string{extfilter.LinuxBinary}).
				WithLoggerName("Binary Filter").
				WithChannelMax(binExtractor).
				Build()
			if err != nil {
				return err
			}
			pipe = pipeline.Process(pipe, binWithinSizeStrategy)

			artifactSlice := make([]filter.Artifact, 0)
			for artifact := range pipe {
				artifactSlice = append(artifactSlice, artifact)
			}
			releases := filter.ReleasesInfo{
				Version:   info.Version,
				Artifacts: artifactSlice,
			}
			mInfo, err := json.Marshal(releases)

			if err != nil {
				return err
			}
			fmt.Println(string(mInfo))

			return err
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		if logger != nil {
			logger.Panic(err.Error())
		}
		panic(err.Error())
	}
}
