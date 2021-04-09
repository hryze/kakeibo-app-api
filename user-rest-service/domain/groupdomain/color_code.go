package groupdomain

import (
	"regexp"

	"golang.org/x/xerrors"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/userdomain"
)

type ColorCode string

const (
	colorCodeFormat = `^#([0-9A-F]{6})$`

	red                  = "#FF0000"
	cyan                 = "#00FFFF"
	chartreuseGreen      = "#80FF00"
	violet               = "#8000FF"
	orange               = "#FF8000"
	azure                = "#0080FF"
	emeraldGreen         = "#00FF80"
	rubyRed              = "#FF0080"
	yellow               = "#FFFF00"
	blue                 = "#0000FF"
	green                = "#00FF00"
	magenta              = "#FF00FF"
	goldenYellow         = "#FFBF00"
	cobaltBlue           = "#0040FF"
	brightYellowishGreen = "#BFFF00"
	hyacinth             = "#4000FF"
	cobaltGreen          = "#00FF40"
	reddishPurple        = "#FF00BF"
	leafGreen            = "#40FF00"
	purple               = "#BF00FF"
	vermilion            = "#FF4000"
	ceruleanBlue         = "#00BFFF"
	turquoiseGreen       = "#00FFBF"
	carmine              = "#FF0040"
)

var (
	colorCodeRegex = regexp.MustCompile(colorCodeFormat)

	colorCodeList = [24]string{
		red,
		cyan,
		chartreuseGreen,
		violet,
		orange,
		azure,
		emeraldGreen,
		rubyRed,
		yellow,
		blue,
		green,
		magenta,
		goldenYellow,
		cobaltBlue,
		brightYellowishGreen,
		hyacinth,
		cobaltGreen,
		reddishPurple,
		leafGreen,
		purple,
		vermilion,
		ceruleanBlue,
		turquoiseGreen,
		carmine,
	}
)

func NewColorCode(colorCode string) (ColorCode, error) {
	if ok := colorCodeRegex.MatchString(colorCode); !ok {
		return "", xerrors.Errorf("invalid colorCode format: %s", colorCode)
	}

	return ColorCode(colorCode), nil
}

func NewColorCodeToUser(approvedUserIDList []userdomain.UserID) (ColorCode, error) {
	if len(approvedUserIDList) == 0 {
		return "", xerrors.Errorf("approvedUserIDList must be more than 1 element: %d", len(approvedUserIDList))
	}

	idx := len(approvedUserIDList) % len(colorCodeList)
	colorCode := colorCodeList[idx]

	return ColorCode(colorCode), nil
}

func (c ColorCode) Value() string {
	return string(c)
}
