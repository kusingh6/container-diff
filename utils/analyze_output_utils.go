package utils

type Result interface {
	OutputStruct() interface{}
	OutputText(resultType string) error
}

type AnalyzeResult struct {
	Image       string
	AnalyzeType string
	Analysis    interface{}
}

type ListAnalyzeResult AnalyzeResult

func (r ListAnalyzeResult) OutputStruct() interface{} {
	return r
}

func (r ListAnalyzeResult) OutputText(resultType string) error {
	r.Analysis = r.Analysis.([]string)
	return TemplateOutput(r, "ListAnalyze")
}

type MultiVersionPackageAnalyzeResult AnalyzeResult

func (r MultiVersionPackageAnalyzeResult) OutputStruct() interface{} {
	analysis := r.Analysis.(map[string]map[string]PackageInfo)
	analysisOutput := getMultiVersionPackageOutput(analysis)
	output := struct {
		Image       string
		AnalyzeType string
		Analysis    []PackageOutput
	}{
		Image:       r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis:    analysisOutput,
	}
	return output
}

func (r MultiVersionPackageAnalyzeResult) OutputText(resultType string) error {
	analysis := r.Analysis.(map[string]map[string]PackageInfo)
	analysisOutput := getMultiVersionPackageOutput(analysis)

	strAnalysis := stringifyPackages(analysisOutput)
	strResult := struct {
		Image       string
		AnalyzeType string
		Analysis    []StrPackageOutput
	}{
		Image:       r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis:    strAnalysis,
	}
	return TemplateOutput(strResult, "MultiVersionPackageAnalyze")
}

type SingleVersionPackageAnalyzeResult AnalyzeResult

func (r SingleVersionPackageAnalyzeResult) OutputStruct() interface{} {
	analysis := r.Analysis.(map[string]PackageInfo)
	analysisOutput := getSingleVersionPackageOutput(analysis)
	output := struct {
		Image       string
		AnalyzeType string
		Analysis    []PackageOutput
	}{
		Image:       r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis:    analysisOutput,
	}
	return output
}

func (r SingleVersionPackageAnalyzeResult) OutputText(diffType string) error {
	analysis := r.Analysis.(map[string]PackageInfo)
	analysisOutput := getSingleVersionPackageOutput(analysis)

	strAnalysis := stringifyPackages(analysisOutput)
	strResult := struct {
		Image       string
		AnalyzeType string
		Analysis    []StrPackageOutput
	}{
		Image:       r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis:    strAnalysis,
	}
	return TemplateOutput(strResult, "SingleVersionPackageAnalyze")
}

type PackageOutput struct {
	Name    string
	Path    string `json:",omitempty"`
	Version string
	Size    int64
}

func getSingleVersionPackageOutput(packageMap map[string]PackageInfo) []PackageOutput {
	packages := []PackageOutput{}
	for name, info := range packageMap {
		packages = append(packages, PackageOutput{Name: name, Version: info.Version, Size: info.Size})
	}

	if SortSize {
		packageBy(packageSizeSort).Sort(packages)
	} else {
		packageBy(packageNameSort).Sort(packages)
	}
	return packages
}

func getMultiVersionPackageOutput(packageMap map[string]map[string]PackageInfo) []PackageOutput {
	packages := []PackageOutput{}
	for name, versionMap := range packageMap {
		for path, info := range versionMap {
			packages = append(packages, PackageOutput{Name: name, Path: path, Version: info.Version, Size: info.Size})
		}
	}

	if SortSize {
		packageBy(packageSizeSort).Sort(packages)
	} else {
		packageBy(packageNameSort).Sort(packages)
	}
	return packages
}

type FileAnalyzeResult AnalyzeResult

func (r FileAnalyzeResult) OutputStruct() interface{} {
	analysis := r.Analysis.([]DirectoryEntry)
	if SortSize {
		directoryBy(directorySizeSort).Sort(analysis)
	} else {
		directoryBy(directoryNameSort).Sort(analysis)
	}
	r.Analysis = analysis
	return r
}

func (r FileAnalyzeResult) OutputText(analyzeType string) error {
	analysis := r.Analysis.([]DirectoryEntry)
	if SortSize {
		directoryBy(directorySizeSort).Sort(analysis)
	} else {
		directoryBy(directoryNameSort).Sort(analysis)
	}
	strAnalysis := stringifyDirectoryEntries(analysis)

	strResult := struct {
		Image       string
		AnalyzeType string
		Analysis    []StrDirectoryEntry
	}{
		Image:       r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis:    strAnalysis,
	}
	return TemplateOutput(strResult, "FileAnalyze")
}
