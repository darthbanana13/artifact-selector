package arch

import (
	"regexp"
	"strings"

	"github.com/darthbanana13/artifact-selector/internal/filter"
	"github.com/darthbanana13/artifact-selector/internal/filter/separator"
)

var archRegexCache = make(map[string][]*regexp.Regexp)

func init() {
	for arch, aliases := range ArchMap {
		regexes := make([]*regexp.Regexp, 0, len(aliases))
		for _, alias := range aliases {
			regexes = append(regexes, separator.MakeAliasRegex(alias))
		}
		archRegexCache[arch] = regexes
	}
}

const (
	Exact       = "exact"
	Missing     = "missing"
	MetadataKey = "arch"
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
	"m68k":			 {"m68k"},
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
	fileNameLower := strings.ToLower(artifact.FileName)
	if MatchesArch(fileNameLower, a.TargetArch) {
		artifact.Metadata, _ = filter.AddMetadata(artifact.Metadata, MetadataKey, Exact)
		return artifact, true
	} else if a.TargetArch == "x86_64" && !MatchesOtherArch(fileNameLower, a.TargetArch) {
		artifact.Metadata, _ = filter.AddMetadata(artifact.Metadata, MetadataKey, Missing)
		return artifact, true
	}
	return artifact, false
}

func MatchesArch(fileNameLower string, targetArch string) bool {
	for _, r := range archRegexCache[targetArch] {
		if r.MatchString(fileNameLower) {
			return true
		}
	}
	return false
}

func MatchesOtherArch(fileNameLower string, besidesArch string) bool {
	for arch, regexes := range archRegexCache {
		if arch == besidesArch {
			continue
		}
		for _, r := range regexes {
			if r.MatchString(fileNameLower) {
				return true
			}
		}
	}
	return false
}
