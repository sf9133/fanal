package npm

import (
	"os"
	"path/filepath"

	"golang.org/x/xerrors"

	"github.com/sf9133/fanal/analyzer"
	"github.com/sf9133/fanal/analyzer/library"
	"github.com/sf9133/fanal/types"
	"github.com/sf9133/fanal/utils"
	"github.com/sf9133/go-dep-parser/pkg/npm"
)

func init() {
	analyzer.RegisterAnalyzer(&npmLibraryAnalyzer{})
}

const version = 1

var requiredFiles = []string{"package-lock.json"}

type npmLibraryAnalyzer struct{}

func (a npmLibraryAnalyzer) Analyze(target analyzer.AnalysisTarget) (*analyzer.AnalysisResult, error) {
	res, err := library.Analyze(types.Npm, target.FilePath, target.Content, npm.Parse)
	if err != nil {
		return nil, xerrors.Errorf("unable to parse package-lock.json: %w", err)
	}
	return res, nil
}

func (a npmLibraryAnalyzer) Required(filePath string, _ os.FileInfo) bool {
	fileName := filepath.Base(filePath)
	return utils.StringInSlice(fileName, requiredFiles)
}

func (a npmLibraryAnalyzer) Type() analyzer.Type {
	return analyzer.TypeNpm
}

func (a npmLibraryAnalyzer) Version() int {
	return version
}
