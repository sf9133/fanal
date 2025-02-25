package debian

import (
	"bufio"
	"bytes"
	"os"

	"github.com/sf9133/fanal/analyzer"
	aos "github.com/sf9133/fanal/analyzer/os"
	"golang.org/x/xerrors"

	"github.com/sf9133/fanal/types"
	"github.com/sf9133/fanal/utils"
)

func init() {
	analyzer.RegisterAnalyzer(&debianOSAnalyzer{})
}

const version = 1

var requiredFiles = []string{"etc/debian_version"}

type debianOSAnalyzer struct{}

func (a debianOSAnalyzer) Analyze(target analyzer.AnalysisTarget) (*analyzer.AnalysisResult, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(target.Content))
	for scanner.Scan() {
		line := scanner.Text()
		return &analyzer.AnalysisResult{
			OS: &types.OS{Family: aos.Debian, Name: line},
		}, nil
	}
	return nil, xerrors.Errorf("debian: %w", aos.AnalyzeOSError)
}

func (a debianOSAnalyzer) Required(filePath string, _ os.FileInfo) bool {
	return utils.StringInSlice(filePath, requiredFiles)
}

func (a debianOSAnalyzer) Type() analyzer.Type {
	return analyzer.TypeDebian
}

func (a debianOSAnalyzer) Version() int {
	return version
}
