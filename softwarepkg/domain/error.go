package domain

import "errors"

const (
	codeSoftwarePkgNotImporter = "software_pkg_not_importer"
	codeSoftwarePkgCIIsRunning = "software_pkg_ci_is_running"
)

var (
	errorCIIsRunning    = errors.New("ci is running")
	errorNotTheImporter = errors.New("not the importer")
)

func ParseErrorCode(err error) string {
	if errors.Is(err, errorNotTheImporter) {
		return codeSoftwarePkgNotImporter
	}

	if errors.Is(err, errorCIIsRunning) {
		return codeSoftwarePkgCIIsRunning
	}

	return ""
}
