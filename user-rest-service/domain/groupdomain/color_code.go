package groupdomain

import (
	"regexp"

	"golang.org/x/xerrors"
)

type ColorCode string

const colorCodeFormat = `^#([0-9A-F]{6})$`

var colorCodeRegex = regexp.MustCompile(colorCodeFormat)

func NewColorCode(colorCode string) (ColorCode, error) {
	if ok := colorCodeRegex.MatchString(colorCode); !ok {
		return "", xerrors.Errorf("invalid colorCode format: %s", colorCode)
	}

	return ColorCode(colorCode), nil
}

func (c ColorCode) Value() string {
	return string(c)
}
