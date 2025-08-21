package os

import (
	"slices"
	"strings"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
	"github.com/darthbanana13/artifact-selector/pkg/filter/separator"
)

var OSMap = map[string][]string{
	"linux":   {"linux64", "linux"},
	"android": {"android"},
	"windows": {"windows", "win64", "win32", "win11", "win10", "win"},
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
	targetOS            string
	targetDistroAliases []string
	targetOSAliases     []string
	excludedAliases     []string
}

func (o *OS) SetTargetOS(targetOS string) error {
	o.targetOS = targetOS
	o.targetDistroAliases, o.targetOSAliases, o.excludedAliases = PartitionOSAliases(o.targetOS)
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
	if IsInAliases(o.targetDistroAliases, artifact.FileName) {
		artifact.Metadata["os"] = "distro"
		return artifact, true
	} else if IsInAliases(o.targetOSAliases, artifact.FileName) {
		artifact.Metadata["os"] = "os"
		return artifact, true
	} else if DoesntMatchAliases(o.excludedAliases, artifact.FileName) {
		artifact.Metadata["os"] = "missing"
		return artifact, true
	}
	return artifact, false
}

func PartitionOSAliases(targetOS string) (targetDistroAliases, targetOSAliases, excludedAliases []string) {
	if IsOSNameADistro(targetOS) {
		targetDistroAliases = append(targetDistroAliases, DistroMap[targetOS]...)
		excludedAliases = append(excludedAliases, GetExcludedDistros(targetOS)...)
		targetOS = "linux"
	}
	targetOSAliases = OSMap[targetOS]
	excludedAliases = append(excludedAliases, GetExcludedOSes(targetOS)...)
	return targetDistroAliases, targetOSAliases, excludedAliases
}

func MatchesAlias(alias, s string) bool {
	r := separator.MakeAliasRegex(alias)
	return r.MatchString(strings.ToLower(s))
}

func IsInAliases(aliases []string, s string) bool {
	for _, alias := range aliases {
		if MatchesAlias(alias, s) {
			return true
		}
	}
	return false
}

func DoesntMatchAliases(aliases []string, s string) bool {
	for _, alias := range aliases {
		r := separator.MakeAliasRegex(alias)
		if r.MatchString(strings.ToLower(s)) {
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
