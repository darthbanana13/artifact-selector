package os

import (
	"slices"
	"strings"

	"github.com/darthbanana13/artifact-selector/pkg/github"
)

var OSMap = map[string][]string{
	"linux":   {"linux64", "linux"},
	"android": {"android"},
	"windows": {"windows", "win64", "win32", "win"},
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
	excludedAliases []string
}

func (o *OS) SetTargetOS(targetOS string) error {
	o.targetOS = targetOS
	o.targetAliases, o.excludedAliases = PartitionOSAliases(o.targetOS)
	return nil
}

func (o *OS) TargetOS() string {
	return o.targetOS
}

func NewOSFilter(targetOS string) (IOS, error) {
	o := &OS{}
	err := o.SetTargetOS(targetOS)
	return o, err
}

// TODO: Should this filter give this hint about sorting, 1st distro then os?
func (o *OS) FilterArtifact(artifact github.Artifact) (github.Artifact, bool) {
	for _, osName := range o.targetAliases {
		if MatchesAlias(osName, artifact.FileName) {
			return artifact, true
		}
	}
	if DoesntMatchAliases(o.excludedAliases, artifact.FileName) {
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

func MatchesAlias(s, osName string) bool {
	return strings.Contains(strings.ToLower(s), osName)
}

func DoesntMatchAliases(oses []string, s string) bool {
	for _, osName := range oses {
		if strings.Contains(strings.ToLower(s), osName) {
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
