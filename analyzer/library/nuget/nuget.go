package nuget

import (
	"os"
	"path/filepath"

	"golang.org/x/xerrors"

	"github.com/sf9133/fanal/analyzer"
	"github.com/sf9133/fanal/analyzer/library"
	"github.com/sf9133/fanal/types"
	"github.com/sf9133/fanal/utils"
	"github.com/sf9133/go-dep-parser/pkg/nuget"
)

func init() {
	analyzer.RegisterAnalyzer(&nugetLibraryAnalyzer{})
}

const version = 1

var requiredFiles = []string{"packages.lock.json"}

type nugetLibraryAnalyzer struct{}

func (a nugetLibraryAnalyzer) Analyze(target analyzer.AnalysisTarget) (*analyzer.AnalysisResult, error) {
	res, err := library.Analyze(types.NuGet, target.FilePath, target.Content, nuget.Parse)
	if err != nil {
		return nil, xerrors.Errorf("unable to parse packages.lock.json: %w", err)
	}
	return res, nil
}

func (a nugetLibraryAnalyzer) Required(filePath string, _ os.FileInfo) bool {
	fileName := filepath.Base(filePath)
	return utils.StringInSlice(fileName, requiredFiles)
}

func (a nugetLibraryAnalyzer) Type() analyzer.Type {
	return analyzer.TypeNuget
}

func (a nugetLibraryAnalyzer) Version() int {
	return version
}
