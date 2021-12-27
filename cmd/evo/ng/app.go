package ng

import (
	"github.com/getevo/evo-ng/internal/file"
	"github.com/getevo/evo-ng/internal/proc"
)

func GetSkeleton(path string) Skeleton {
	var skeleton Skeleton
	if err := file.ParseJSON(path, &skeleton); err != nil {
		proc.Die(err)
	}
	skeleton.Packages = map[string]*Package{}
	return skeleton
}
