package alpine

import (
	"bufio"
	"bytes"
	"os"

	"golang.org/x/xerrors"

	"github.com/sf9133/fanal/analyzer"
	aos "github.com/sf9133/fanal/analyzer/os"
	"github.com/sf9133/fanal/types"
	"github.com/sf9133/fanal/utils"
)

func init() {
	analyzer.RegisterAnalyzer(&alpineOSAnalyzer{})
}

const version = 1

var requiredFiles = []string{"etc/alpine-release"}

type alpineOSAnalyzer struct{}

func (a alpineOSAnalyzer) Analyze(target analyzer.AnalysisTarget) (*analyzer.AnalysisResult, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(target.Content))
	for scanner.Scan() {
		line := scanner.Text()
		return &analyzer.AnalysisResult{
			OS: &types.OS{Family: aos.Alpine, Name: line},
		}, nil
	}
	return nil, xerrors.Errorf("alpine: %w", aos.AnalyzeOSError)
}

func (a alpineOSAnalyzer) Required(filePath string, _ os.FileInfo) bool {
	return utils.StringInSlice(filePath, requiredFiles)
}

func (a alpineOSAnalyzer) Type() analyzer.Type {
	return analyzer.TypeAlpine
}

func (a alpineOSAnalyzer) Version() int {
	return version
}
