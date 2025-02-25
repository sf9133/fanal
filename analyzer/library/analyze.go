package library

import (
	"bytes"
	"io"

	"golang.org/x/xerrors"

	"github.com/sf9133/fanal/analyzer"
	"github.com/sf9133/fanal/types"
	godeptypes "github.com/sf9133/go-dep-parser/pkg/types"
)

type parser func(r io.Reader) ([]godeptypes.Library, error)

func Analyze(analyzerType, filePath string, content []byte, parse parser) (*analyzer.AnalysisResult, error) {
	r := bytes.NewReader(content)
	parsedLibs, err := parse(r)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse %s: %w", filePath, err)
	}

	if len(parsedLibs) == 0 {
		return nil, nil
	}

	return ToAnalysisResult(analyzerType, filePath, parsedLibs), nil
}

func ToAnalysisResult(analyzerType, filePath string, libs []godeptypes.Library) *analyzer.AnalysisResult {
	var libInfos []types.LibraryInfo
	for _, lib := range libs {
		libInfos = append(libInfos, types.LibraryInfo{
			Library: lib,
		})
	}
	apps := []types.Application{{
		Type:      analyzerType,
		FilePath:  filePath,
		Libraries: libInfos,
	}}

	return &analyzer.AnalysisResult{Applications: apps}
}
