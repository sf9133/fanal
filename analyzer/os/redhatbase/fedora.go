package redhatbase

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	"github.com/sf9133/fanal/analyzer"

	aos "github.com/sf9133/fanal/analyzer/os"
	"github.com/sf9133/fanal/types"
	"github.com/sf9133/fanal/utils"
	"golang.org/x/xerrors"
)

const fedoraAnalyzerVersion = 1

func init() {
	analyzer.RegisterAnalyzer(&fedoraOSAnalyzer{})
}

type fedoraOSAnalyzer struct{}

func (a fedoraOSAnalyzer) Analyze(target analyzer.AnalysisTarget) (*analyzer.AnalysisResult, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(target.Content))
	for scanner.Scan() {
		line := scanner.Text()
		result := redhatRe.FindStringSubmatch(strings.TrimSpace(line))
		if len(result) != 3 {
			return nil, xerrors.New("cent: Invalid fedora-release")
		}

		switch strings.ToLower(result[1]) {
		case "fedora", "fedora linux":
			return &analyzer.AnalysisResult{
				OS: &types.OS{Family: aos.Fedora, Name: result[2]},
			}, nil
		}
	}
	return nil, xerrors.Errorf("fedora: %w", aos.AnalyzeOSError)
}

func (a fedoraOSAnalyzer) Required(filePath string, _ os.FileInfo) bool {
	return utils.StringInSlice(filePath, a.requiredFiles())
}

func (a fedoraOSAnalyzer) requiredFiles() []string {
	return []string{
		"etc/fedora-release",
		"usr/lib/fedora-release",
	}
}

func (a fedoraOSAnalyzer) Type() analyzer.Type {
	return analyzer.TypeFedora
}

func (a fedoraOSAnalyzer) Version() int {
	return fedoraAnalyzerVersion
}
