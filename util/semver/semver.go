package semver

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hokonco/gdk/util/cast"
)

// SemanticVersion ...
type SemanticVersion struct {
	Major uint8
	Minor uint8
	Patch uint8
	Raw   string
	Tag   string
}

// Parse ...
func Parse(s string) SemanticVersion {
	var semver SemanticVersion
	semver.Raw = s
	for i, _s := range strings.Split(s, ".") {
		u8, _ := cast.StringToUint8(regexp.MustCompile(`\D*`).ReplaceAllString(_s, ""))
		switch i {
		case 0:
			semver.Major = u8
		case 1:
			semver.Minor = u8
		case 2:
			semver.Patch = u8
		}
	}
	semver.Tag = fmt.Sprintf("v%d.%d.%d", semver.Major, semver.Minor, semver.Patch)
	return semver
}
