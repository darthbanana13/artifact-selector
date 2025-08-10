package arch

import (
	"strings"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
	"github.com/darthbanana13/artifact-selector/pkg/filter/separator"
)

var ArchMap = map[string][]string{
	"x86_64":    {"x86_64", "amd64", "x64", "win64", "linux64"},
	"x86":       {"x86", "i386", "386", "i486", "i586", "i686", "i786"},
	"arm64":     {"aarch64", "arm64v8l", "arm64v8", "arm64"},
	"arm32":     {"armv5l", "armv5", "armv6l", "armv6", "armv7l", "armv7", "armv8l", "armv8", "armhf", "armel", "arm"},
	"riscv64":   {"riscv64"},
	"loongarch": {"loongarch64le", "loongarch64be", "loongarch64", "loongarch"},
	"s390":      {"s390x", "s390"},
	"powerpc":   {"powerpc64", "powerpc", "ppc64le", "ppc64el", "ppc64"},
	"mips":      {"mipsel", "mipsr6el", "mipsr6le", "mipsr6", "mips32", "mips64le", "mips64r6le", "mips64r6", "mips64", "mipsle", "mips"},
	"sparc":     {"sparc64", "sparc"},
	"ia64":      {"ia64"},
}

type Arch struct {
	TargetArch string
}

func (a *Arch) SetTargetArch(targetArch string) error {
	a.TargetArch = targetArch
	return nil
}

func NewArch(targetArch string) (IArch, error) {
	a := &Arch{}
	return a, a.SetTargetArch(targetArch)
}

func (a *Arch) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	if MatchesArch(artifact.Source.FileName, a.TargetArch) {
		artifact.Metadata["arch"] = "exact"
		return artifact, true
	} else if a.TargetArch == "x86_64" && !MatchesOtherArch(artifact.Source.FileName, a.TargetArch) {
		artifact.Metadata["arch"] = "missing"
		return artifact, true
	}
	return artifact, false
}

// TODO: Cache all the regexes for better performance
func MatchesArch(fileName string, targetArch string) bool {
	for _, alias := range ArchMap[targetArch] {
		r := separator.MakeAliasRegex(alias)
		if r.MatchString(strings.ToLower(fileName)) {
			return true
		}
	}
	return false
}

func MatchesOtherArch(fileName string, besidesArch string) bool {
	for arch, aliases := range ArchMap {
		if arch == besidesArch {
			continue
		}
		for _, alias := range aliases {
			r := separator.MakeAliasRegex(alias)
			if r.MatchString(strings.ToLower(fileName)) {
				return true
			}
		}
	}
	return false
}
