package domain

import "errors"

const (
	errorSoftwarePkgNotImporter = "software_pkg_not_importer"
)

var (
	errorNotTheImporter = errors.New("not the importer")
)

func ParseErrorCode(err error) string {
	if errors.Is(err, errorNotTheImporter) {
		return errorSoftwarePkgNotImporter
	}

	return ""
}
