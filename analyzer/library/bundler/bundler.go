package bundler

import (
	"os"
	"path/filepath"

	"golang.org/x/xerrors"

	"github.com/sf9133/fanal/analyzer"
	"github.com/sf9133/fanal/analyzer/library"
	"github.com/sf9133/fanal/types"
	"github.com/sf9133/fanal/utils"
	"github.com/sf9133/go-dep-parser/pkg/bundler"
)

func init() {
	analyzer.RegisterAnalyzer(&bundlerLibraryAnalyzer{})
}

const version = 1

var (
	requiredFiles = []string{"Gemfile.lock"}
)

type bundlerLibraryAnalyzer struct{}

func (a bundlerLibraryAnalyzer) Analyze(target analyzer.AnalysisTarget) (*analyzer.AnalysisResult, error) {
	res, err := library.Analyze(types.Bundler, target.FilePath, target.Content, bundler.Parse)
	if err != nil {
		return nil, xerrors.Errorf("unable to parse Gemfile.lock: %w", err)
	}
	return res, nil
}

func (a bundlerLibraryAnalyzer) Required(filePath string, _ os.FileInfo) bool {
	fileName := filepath.Base(filePath)
	return utils.StringInSlice(fileName, requiredFiles)
}

func (a bundlerLibraryAnalyzer) Type() analyzer.Type {
	return analyzer.TypeBundler
}

func (a bundlerLibraryAnalyzer) Version() int {
	return version
}
