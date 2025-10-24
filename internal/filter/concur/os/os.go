package os

import (
	"regexp"
	"slices"
	"strings"

	"github.com/darthbanana13/artifact-selector/internal/filter"
	"github.com/darthbanana13/artifact-selector/internal/filter/separator"
)

const (
	Missing = "missing"
)

var OSMap = map[string][]string{
	"linux":   {"linux64", "linux"},
	"android": {"android"},
	"windows": {"win11", "win10", "win64", "windows", "win32", "win"},
	"macos":   {"macos", "mac", "darwin", "osx", "apple"},
	"freebsd": {"freebsd", "bsd"},
	"openbsd": {"openbsd", "bsd"},
	"netbsd":  {"netbsd", "bsd"},
}

// TODO: Should we also distinguish between versions of win, mac, or android?
var DistroMap = map[string][]string{
	"debian": {"debian"},
	"ubuntu": {"ubuntu", "debian"},
	"fedora": {"fedora", "rhel"},
	"redhat": {"redhat", "rhel"},
}

type OS struct {
	targetOS        string
	targetAliases   []string
	targetRegexes   []*regexp.Regexp
	excludedAliases []string
	excludedRegexes []*regexp.Regexp
}

func (o *OS) SetTargetOS(targetOS string) error {
	o.targetOS = targetOS
	o.targetAliases, o.excludedAliases = PartitionOSAliases(o.targetOS)
	o.targetRegexes, o.excludedRegexes = MakeOSRegexes(o.targetAliases), MakeOSRegexes(o.excludedAliases)
	return nil
}

func (o *OS) TargetOS() string {
	return o.targetOS
}

func NewOS(targetOS string) (IOS, error) {
	o := &OS{}
	err := o.SetTargetOS(targetOS)
	return o, err
}

func (o *OS) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	i := IndexInAliases(o.targetRegexes, artifact.FileName)
	if i >= 0 {
		artifact.Metadata, _ = filter.AddMetadata(artifact.Metadata, "os", o.targetAliases[i], "os-index", len(o.targetAliases) - i)
		return artifact, true
	} else if DoesntMatchAliases(o.excludedRegexes, artifact.FileName) {
		artifact.Metadata, _ = filter.AddMetadata(artifact.Metadata, "os", Missing, "os-index", i)
		return artifact, true
	}
	return artifact, false
}

func PartitionOSAliases(targetOS string) (targetAliases, excludedAliases []string) {
	if IsOSNameADistro(targetOS) {
		targetAliases = append(targetAliases, DistroMap[targetOS]...)
		excludedAliases = append(excludedAliases, GetExcludedDistros(targetOS)...)
		targetOS = "linux"
	}
	targetAliases = append(targetAliases, OSMap[targetOS]...)
	excludedAliases = append(excludedAliases, GetExcludedOSes(targetOS)...)
	return targetAliases, excludedAliases
}

func MakeOSRegexes(aliases []string) []*regexp.Regexp {
	out := make([]*regexp.Regexp, len(aliases))
	for aliasIndex, alias := range aliases {
		out[aliasIndex] = separator.MakeAliasRegex(alias)
	}
	return out
}

func IndexInAliases(regexes []*regexp.Regexp, s string) int {
	for index, regex := range regexes {
		if regex.MatchString(strings.ToLower(s)) {
			return index
		}
	}
	return -1
}

func DoesntMatchAliases(regexes []*regexp.Regexp, s string) bool {
	for _, regex := range regexes {
		if regex.MatchString(strings.ToLower(s)) {
			return false
		}
	}
	return true
}

func ValuesNotInKey(key string, hashmap map[string][]string) []string {
	difference := []string{}
	for k, values := range hashmap {
		if k == key {
			continue
		}
		for _, value := range values {
			if !slices.Contains(hashmap[key], value) {
				difference = append(difference, value)
			}
		}
	}
	return difference
}

func GetExcludedOSes(targetOS string) []string {
	return ValuesNotInKey(targetOS, OSMap)
}

func GetExcludedDistros(distroName string) []string {
	return ValuesNotInKey(distroName, DistroMap)
}

func IsOSNameADistro(osName string) bool {
	if _, ok := DistroMap[osName]; ok {
		return true
	}
	return false
}
