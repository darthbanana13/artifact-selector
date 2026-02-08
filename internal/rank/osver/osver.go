package osver

import (
	"strconv"
	"strings"

	"github.com/darthbanana13/artifact-selector/internal/filter"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/os/metadata/osver"
)

const (
	//Major version mismatch
	MaxRankMajorDiff = 4
	MinRankMajorDiff = 1

	//Minor version mismatch
	MaxRankMinorDiff = 7
	MinRankMinorDiff = 5

	//Patch or lower version mismatch
	PatchRankDiff = 8
)

// NOTE: This currently has the same limitation as the OSVer filter regarding named vs numerical versions
type OSVer struct {
	Ver string
}

func NewOSVer(version string) *OSVer {
	return &OSVer{
		Ver: version,
	}
}

func (o *OSVer) RankArtifact(artifact filter.Artifact) uint {
	foundVer := artifact.Metadata[osver.MetadataKey].(string)
	if o.Ver == foundVer {
		return 9
	}

	hasPeriod, foundHasPeriod := strings.Contains(o.Ver, "."), strings.Contains(foundVer, ".")
	if !foundHasPeriod && hasPeriod {
		targetNoPeriod := strings.ReplaceAll(o.Ver, ".", "")
		if len(targetNoPeriod) == len(foundVer) {
			if targetNoPeriod == foundVer {
				return 9
			}
		}
	}
	return CompareWithPeriod(o.Ver, foundVer)
}

func SliceVersion(versionNumber string) ([]int, error) {
	stringSlice := strings.Split(versionNumber, ".")
	intSlice := make([]int, len(stringSlice))
	var err error
	for i, val := range stringSlice {
		intSlice[i], err = strconv.Atoi(val)
		if err != nil {
			return []int{}, err
		}
	}
	return intSlice, err
}

func CompareWithPeriod(target, found string) uint {
	targetParts, err := SliceVersion(target)
	if err != nil {
		return 1
	}
	foundParts, err := SliceVersion(found)
	if err != nil {
		return 1
	}

	return CompareSlices(targetParts, foundParts)
}

func CompareSlices(target, found []int) uint {
	var i int
	for i = 0; i < len(target) && i < len(found); i++ {
		if target[i] == found[i] {
			continue
		}
		return RankVersionDistance(i, AbsDiff(target[i], found[i]))
	}
	if i >= len(target) {
		return 9
	}
	return RankVersionDistance(i, AbsDiff(target[i], found[i]))
}

func RankVersionDistance(idx, diff int) uint {
	var maxRank, minRank int

	switch idx {
	case 0:
		maxRank, minRank = MaxRankMajorDiff, MinRankMajorDiff
	case 1:
		maxRank, minRank = MaxRankMinorDiff, MinRankMinorDiff
	default:
		return PatchRankDiff
	}

	rank := maxRank - (diff - 1)
	if rank < minRank {
		return uint(minRank)
	}
	return uint(rank)
}

func AbsDiff(min, sub int) int {
	diff := min - sub
	if diff < 0 {
		return diff * -1
	}
	return diff
}
