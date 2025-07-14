package arch

import (
	"strings"

	"github.com/darthbanana13/artifact-selector/pkg/filter/handler"
	"github.com/darthbanana13/artifact-selector/pkg/github"
)

var archMap = map[string][]string{
  "x86_64": {"x86_64", "amd64", "x64", "win64", "linux64"},
  "x86": {"x86", "i386", "386", "i486", "i586", "i686", "i786"},
  "arm64": {"arm64", "aarch64", "arm64v8l", "arm64v8", "aarch64"},
  "arm32": {"arm", "armv5", "armv5l", "armv6", "armv6l", "armv7", "armv7l", "armv8", "armv8l", "armhf", "armel"},
  "riscv64": {"riscv64"},
  "s390": {"s390x", "s390"},
  "powerpc": {"powerpc", "powerpc64", "ppc64le", "ppc64"},
  "mips": {"mipsel", "mipsr6el", "mipsr6le", "mipsr6", "mips32", "mips64le", "mips64", "mipsle", "mips"},
  "sparc": {"sparc", "sparc64"},
  "ia64": {"ia64"},
}

type ArchFilter struct {
  Next        handler.IFilterHandler
  TargetArch  string
}

func (af *ArchFilter) SetNext(next handler.IFilterHandler) {
  af.Next = next
}

func (af *ArchFilter) SetTargetArch(targetArch string) error {
  if _, ok := archMap[targetArch]; !ok {
    return UnsupportedArchErr{Arch: targetArch}
  }
  af.TargetArch = targetArch
  return nil
}

func NewArchFilter(targetArch string) (*ArchFilter, error) {
  af := &ArchFilter{}
  return af, af.SetTargetArch(targetArch)
}

//TODO: Refactor to a smaller version
func (af *ArchFilter) Filter(releases github.ReleasesInfo) github.ReleasesInfo {
  var filteredArtifacts []github.Artifact

  for _, artifact := range releases.Artifacts {
    // artifact.Size <= 102400
    if af.FilterExactMatch(artifact) {
      filteredArtifacts = append(filteredArtifacts, artifact)
    } else if(af.TargetArch == "x86_64" && !af.DoesMatchOtherArch(artifact, af.TargetArch)) { 
      filteredArtifacts = append(filteredArtifacts, artifact)
    }
  }
  
  releases.Artifacts = filteredArtifacts
  if len(filteredArtifacts) > 0 && af.Next != nil {
    return af.Next.Filter(releases)
  }
  return releases
}

func (af *ArchFilter) FilterExactMatch(artifact github.Artifact) bool {
  for _, alias := range archMap[af.TargetArch] {
    if strings.Contains(strings.ToLower(artifact.FileName), alias) {
      return true
    }
  }
  return false
}

func (af *ArchFilter) DoesMatchOtherArch(artifact github.Artifact, besidesArch string) bool {
  for arch, aliases := range archMap {
    if arch == besidesArch {
      continue
    }
    for _, alias := range aliases {
      if strings.Contains(strings.ToLower(artifact.FileName), alias) {
        return true
      }
    }
  }
  return false
}
