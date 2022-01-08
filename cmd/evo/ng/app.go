package ng

import (
	"github.com/getevo/evo-ng/lib/file"
	"github.com/getevo/evo/lib/log"
)

func GetSkeleton(path string) Skeleton {
	var skeleton Skeleton
	var err = file.ParseJSON(path, &skeleton)
	if err != nil {
		log.Error(err)
	}

	skeleton.Packages = map[string]*Package{}
	return skeleton
}
