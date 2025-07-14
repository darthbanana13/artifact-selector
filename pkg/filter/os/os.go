package os

import (
	"strings"

	"github.com/darthbanana13/artifact-selector/pkg/filter/handler"
	"github.com/darthbanana13/artifact-selector/pkg/github"
)

var osMap = map[string][]string{
  "linux":    {"linux64", "linux"},
  "android":  {"android"},
  "windows":  {"windows", "win64", "win32", "win"},
  "macos":    {"macos", "mac", "darwin", "osx"},
}

//TODO: Should we also distinguish between versions of win, mac, or android?
//	Should we separate distros into a separate filter?
var distroMap = map[string][]string{
	"debian":		{"debian"},
	"ubuntu":		{"ubuntu", "debian"},
	"fedora":		{"fedora", "rhel"},
	"redhat":		{"redhat", "rhel"},
}

type OSFilter struct {
	Next			handler.IFilterHandler
	TargetOS	string
}

func (of *OSFilter) SetNext(next handler.IFilterHandler) {
	of.Next = next
}

//TODO: Error handling & logging
func NewOSFilter(targetOS string) (*OSFilter, error) {
	return &OSFilter{TargetOS: targetOS}, nil
}

//TODO: Refactor to a smaller version
//TODO: Should this filter also sort?
func (of *OSFilter) Filter(releases github.ReleasesInfo) github.ReleasesInfo {
	var filteredArtifacts []github.Artifact
	
	targetOSes := []string{}
	nonTargetOSes := []string{}
	targetOS := strings.ToLower(of.TargetOS)
	if of.IsOSNameADistro(targetOS) {
		targetOSes = append(targetOSes, distroMap[targetOS]...)
		nonTargetOSes = append(nonTargetOSes, of.ComputeNonTargetDistros(targetOS)...)
		targetOS = "linux"
	} 
	targetOSes = append(targetOSes, osMap[targetOS]...)
	nonTargetOSes = append(nonTargetOSes, of.ComputeNonTargetOSes(targetOS)...)
	matched := false

	for _, artifact := range releases.Artifacts {
		matched = false
		for _, osName := range targetOSes {
			if of.MatchesOS(artifact.FileName, osName) {
				filteredArtifacts = append(filteredArtifacts, artifact)
				matched = true
			}
		}
		if !matched && of.DoesntMatchOSes(artifact.FileName, nonTargetOSes) {
			matched = true
		}
	}

  releases.Artifacts = filteredArtifacts
  if len(filteredArtifacts) > 0 && of.Next != nil {
    return of.Next.Filter(releases)
  }
  return releases
}

func (of *OSFilter) MatchesOS(fileName, osName string) bool {
	if strings.Contains(strings.ToLower(fileName), osName) {
		return true
	}
	return false
}

func (OSFilter) DoesntMatchOSes(fileName string, oses []string) bool {
	for _, osName := range oses {
		if strings.Contains(strings.ToLower(fileName), osName) {
			return false
		}
	}
	return true
}

func (OSFilter) ComputeNonTargetOSes(targetOS string) []string {
	nonTargetOSes := []string{}
	for osType, osStrings := range osMap {
		if osType == targetOS {
			continue
		}
		nonTargetOSes = append(nonTargetOSes, osStrings...)
	}
	return nonTargetOSes
}

func (OSFilter) IsOSNameADistro(osName string) bool {
	if _, ok := distroMap[osName]; ok {
		return true
	}
	return false
}

func (OSFilter) ComputeNonTargetDistros(distroName string) []string {
	nonTargetDistros := []string{}
	for distro, distroStrings := range distroMap {
		if distroName == distro {
			continue
		}
		nonTargetDistros = append(nonTargetDistros, distroStrings...)
	}
	return nonTargetDistros
}
