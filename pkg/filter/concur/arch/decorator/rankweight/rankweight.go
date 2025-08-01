package rankweight

import (
	"github.com/darthbanana13/artifact-selector/pkg/filter"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/arch"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/arch/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
)

type RankWeightDecorator struct {
	Arch   arch.IArch
	Weight uint
}

func RankWeightDecoratorConstructor(weight uint8) funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(afc decorator.Constructor) decorator.Constructor {
		return func(targetArch string) (arch.IArch, error) {
			af, err := afc(targetArch)
			if err != nil {
				return af, err
			}
			return NewRankWeightDecorator(af, weight), nil
		}
	}
}

func NewRankWeightDecorator(arch arch.IArch, weight uint8) arch.IArch {
	return &RankWeightDecorator{
		Arch:   arch,
		Weight: uint(weight),
	}
}

func (rwd *RankWeightDecorator) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	filterArtifact, keep := rwd.Arch.FilterArtifact(artifact)
	if keep && filterArtifact.Metadata["arch"] == "exact" {
		filterArtifact.Rank += 100 * rwd.Weight
	}
	return filterArtifact, keep
}

func (rwd *RankWeightDecorator) SetTargetArch(targetArch string) error {
	return rwd.Arch.SetTargetArch(targetArch)
}
