package cargo

import (
	"os"
	"path/filepath"

	"golang.org/x/xerrors"

	"github.com/sf9133/fanal/analyzer"
	"github.com/sf9133/fanal/analyzer/library"
	"github.com/sf9133/fanal/types"
	"github.com/sf9133/fanal/utils"
	"github.com/sf9133/go-dep-parser/pkg/cargo"
)

func init() {
	analyzer.RegisterAnalyzer(&cargoLibraryAnalyzer{})
}

const version = 1

var requiredFiles = []string{"Cargo.lock"}

type cargoLibraryAnalyzer struct{}

func (a cargoLibraryAnalyzer) Analyze(target analyzer.AnalysisTarget) (*analyzer.AnalysisResult, error) {
	res, err := library.Analyze(types.Cargo, target.FilePath, target.Content, cargo.Parse)
	if err != nil {
		return nil, xerrors.Errorf("error with Cargo.lock: %w", err)
	}
	return res, nil
}

func (a cargoLibraryAnalyzer) Required(filePath string, _ os.FileInfo) bool {
	fileName := filepath.Base(filePath)
	return utils.StringInSlice(fileName, requiredFiles)
}

func (a cargoLibraryAnalyzer) Type() analyzer.Type {
	return analyzer.TypeCargo
}

func (a cargoLibraryAnalyzer) Version() int {
	return version
}
