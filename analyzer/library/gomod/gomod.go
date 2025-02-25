package gomod

import (
	"os"
	"path/filepath"

	"golang.org/x/xerrors"

	"github.com/sf9133/fanal/analyzer"
	"github.com/sf9133/fanal/analyzer/library"
	"github.com/sf9133/fanal/types"
	"github.com/sf9133/fanal/utils"
	"github.com/sf9133/go-dep-parser/pkg/gomod"
)

func init() {
	analyzer.RegisterAnalyzer(&gomodAnalyzer{})
}

const version = 1

var requiredFiles = []string{"go.sum"}

type gomodAnalyzer struct{}

func (a gomodAnalyzer) Analyze(target analyzer.AnalysisTarget) (*analyzer.AnalysisResult, error) {
	res, err := library.Analyze(types.GoMod, target.FilePath, target.Content, gomod.Parse)
	if err != nil {
		return nil, xerrors.Errorf("failed to analyze %s: %w", target.FilePath, err)
	}
	return res, nil
}

func (a gomodAnalyzer) Required(filePath string, _ os.FileInfo) bool {
	fileName := filepath.Base(filePath)
	return utils.StringInSlice(fileName, requiredFiles)
}

func (a gomodAnalyzer) Type() analyzer.Type {
	return analyzer.TypeGoMod
}

func (a gomodAnalyzer) Version() int {
	return version
}
