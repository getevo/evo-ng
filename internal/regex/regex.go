package regex

import "regexp"

var Version = regexp.MustCompile(`(?m)v([0-9]\.{0,1})+`)
var VersionedPackage = regexp.MustCompile(`^(([a-z0-9\._\-]+\/)+)(v[0-9]+)$`)
