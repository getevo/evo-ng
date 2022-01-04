package evo

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/getevo/evo-ng/internal/shared"
)

func GetView(env string) (*jet.Set, error) {
	if e, ok := shared.Views[env]; !ok {
		return nil, fmt.Errorf("invalid view environment")
	} else {
		return e, nil
	}

}

func RegisterView(prefix, path string) *jet.Set {
	//TODO: Cache and
	var set = jet.NewSet(
		jet.NewOSFileSystemLoader(path),
		jet.InDevelopmentMode(),
	)

	shared.Views[prefix] = set
	return set
}
