package linuxbindiff

import (
	"math"
	"slices"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/ext"
)

var CompressedExtensions = []string{
	"tar.gz",
	"zip",
	"tar.xz",
	"tar.bz2",
	"tbz",
	"tar.zst",
//The next extensions may not have compression
	"deb",
	"rpm",
	"msi",
	"dmg",
	"pkg",
}

//TODO: Refactor this package into something maintainable
//	This filter should apply the same logic for compressed, as it does for binaries
// 	Looking at cases like mikefarah/yq yq_man_page_only.tar.gz 
func Filter(artifacts <-chan filter.Artifact) <-chan filter.Artifact {
	arts := []filter.Artifact{}
	var maxBinSize, avgCompressed, numCompressed uint64 = 0, 0, 0

	for artifact := range artifacts {
		if artifact.Metadata["ext"] == ext.LINUXBINARY {
			if artifact.Source.Size > maxBinSize {
				maxBinSize = artifact.Source.Size
			}
		} else if slices.Contains(CompressedExtensions, artifact.Metadata["ext"].(string)) {
			avgCompressed += artifact.Source.Size
			numCompressed++
		}
		arts = append(arts, artifact)
	}
	output := make(chan filter.Artifact)
	if maxBinSize == 0 {
		go SliceToChannel(arts, output)
		return output
	}

	var compFunc func(uint64) bool
	if numCompressed > 0 {
		avgCompressed /= numCompressed
		compFunc = CompressedComparissonFunc(avgCompressed)
	} else {
		compFunc = AlwaysTrueFunc
	}

	go func() {
		for _, artifact := range arts {
			if artifact.Metadata["ext"] == ext.LINUXBINARY {
				if PercentDiff(artifact.Source.Size, maxBinSize) < 20 || compFunc(artifact.Source.Size) {
					output <- artifact
				}
			} else {
					output <- artifact
			}
		}
		close(output)
	}()
	return output
}

func SliceToChannel(artifacts []filter.Artifact, output chan filter.Artifact) {
	for _, artifact := range artifacts {
		output <- artifact
	}
	close(output)
}

func CompressedComparissonFunc(avgCompressed uint64) func(uint64) bool {
	return func(artifactSize uint64) bool {
		return artifactSize < avgCompressed && PercentDiff(artifactSize, avgCompressed) < 20
	}
}

func AlwaysTrueFunc(artifactSize uint64) bool {
	return true
}

func PercentDiff(a, b uint64) float64 {
	if a == 0 && b == 0 {
		return 0
	}

	fa, fb := float64(a), float64(b)

	num := math.Abs(fa - fb)
	den := (fa + fb) / 2

	return (num / den) * 100
}
