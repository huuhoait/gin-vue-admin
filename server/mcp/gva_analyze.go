package mcpTool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/mark3labs/mcp-go/mcp"
)

// register tool
func init() {
	RegisterTool(&GVAAnalyzer{})
}

// GVAAnalyzer GVAPartAnalyzeDevice - Used forPartAnalyzeCurrentFunctionYesNoNeedCreateIndependentofpackageAndmodule
type GVAAnalyzer struct{}

// AnalyzeRequest PartAnalyzeRequestStructureBody
type AnalyzeRequest struct {
	Requirement string `json:"requirement" binding:"required"` // user requirementdescription
}

// AnalyzeResponse PartAnalyzeResponseStructureBody
type AnalyzeResponse struct {
	ExistingPackages   []PackageInfo           `json:"existingPackages"`   // AppearHavepackage info
PredesignedModules []PredesignedModuleInfo `json:"predesignedModules"` // predesigned module info
	Dictionaries       []DictionaryPre         `json:"dictionaries"`       // DictionaryInformation
	CleanupInfo        *CleanupInfo            `json:"cleanupInfo"`        // CleanInformation(IfHave)
}

// ModuleInfo ModuleInformation
type ModuleInfo struct {
	ModuleName  string   `json:"moduleName"`  // Modulename
	PackageName string   `json:"packageName"` // package name
	Template    string   `json:"template"`    // Templatetype
	StructName  string   `json:"structName"`  // struct name
	TableName   string   `json:"tableName"`   // table name
	Description string   `json:"description"` // description
FilePaths   []string `json:"filePaths"`   // RelatedFilepath
}

// PackageInfo package info
type PackageInfo struct {
	PackageName string `json:"packageName"` // package name
Template    string `json:"template"`    // Templatetype
	Label       string `json:"label"`       // Tag
	Desc        string `json:"desc"`        // description
	Module      string `json:"module"`      // Module
	IsEmpty     bool   `json:"isEmpty"`     // YesNoEmptyPackage
}

// PredesignedModuleInfo predesigned module info
type PredesignedModuleInfo struct {
	ModuleName  string   `json:"moduleName"`  // Modulename
	PackageName string   `json:"packageName"` // package name
Template    string   `json:"template"`    // Templatetype
	FilePaths   []string `json:"filePaths"`   // FilepathList
	Description string   `json:"description"` // description
}

// CleanupInfo CleanInformation
type CleanupInfo struct {
	DeletedPackages []string `json:"deletedPackages"` // AlreadydeleteofPackage
	DeletedModules  []string `json:"deletedModules"`  // Alreadydelete module
	CleanupMessage  string   `json:"cleanupMessage"`  // CleanMessage
}

// New CreateGVAPartAnalyzeDeviceUtility
func (g *GVAAnalyzer) New() mcp.Tool {
	return mcp.NewTool("gva_analyze",
		mcp.WithDescription("Returns valid packages and modules, analyzes whether new packages/modules/dictionaries are needed, checks and cleans empty packages."),
		mcp.WithString("requirement",
			mcp.Description("user requirementdescription, Used forPartAnalyzeYesNoNeedCreateNewofPackageAndModule"),
			mcp.Required(),
		),
	)
}

// Handle HandlePartAnalyzeRequest
func (g *GVAAnalyzer) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parserequest parameters
	requirementStr, ok := request.GetArguments()["requirement"].(string)
	if !ok || requirementStr == "" {
		return nil, errors.New("Invalid parameters:requirement MustYesnon-empty string")
	}

	// CreatePartAnalyzeRequest
	analyzeReq := AnalyzeRequest{
		Requirement: requirementStr,
	}

	// run analysis logic
	response, err := g.performAnalysis(ctx, analyzeReq)
	if err != nil {
		return nil, fmt.Errorf("PartAnalyzefailed: %v", err)
	}

	// serialize response
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("serialize responsefailed: %v", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(responseJSON)),
		},
	}, nil
}

// performAnalysis run analysis logic
func (g *GVAAnalyzer) performAnalysis(ctx context.Context, req AnalyzeRequest) (*AnalyzeResponse, error) {
	_ = req

	packages, err := fetchAutoCodePackages(ctx)
	if err != nil {
		return nil, fmt.Errorf("getpackage infofailed: %v", err)
	}

	histories, err := fetchAutoCodeHistories(ctx)
	if err != nil {
		return nil, fmt.Errorf("getHistoryRecordfailed: %v", err)
	}

	cleanupInfo := &CleanupInfo{
		DeletedPackages: []string{},
		DeletedModules:  []string{},
	}

	validPackages := make([]PackageInfo, 0, len(packages))
	var emptyPackageHistoryIDs []uint

	for _, pkg := range packages {
		isEmpty, err := g.isPackageFolderEmpty(pkg.PackageName, pkg.Template)
		if err != nil {
global.GVA_LOG.Warn(fmt.Sprintf("error while checking whether package %s is empty: %v", pkg.PackageName, err))
			continue
		}

		if isEmpty {
			if err := g.removeEmptyPackageFolder(pkg.PackageName, pkg.Template); err != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("deleteempty package folder %s failed: %v", pkg.PackageName, err))
			} else {
				cleanupInfo.DeletedPackages = append(cleanupInfo.DeletedPackages, pkg.PackageName)
			}

			if err := deleteAutoCodePackage(ctx, pkg.ID); err != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("deletePackageDatabaseRecord %s failed: %v", pkg.PackageName, err))
			}

			for _, history := range histories {
				if history.Package == pkg.PackageName {
					emptyPackageHistoryIDs = append(emptyPackageHistoryIDs, history.ID)
					cleanupInfo.DeletedModules = append(cleanupInfo.DeletedModules, history.StructName)
				}
			}
			continue
		}

		validPackages = append(validPackages, PackageInfo{
			PackageName: pkg.PackageName,
			Template:    pkg.Template,
			Label:       pkg.Label,
			Desc:        pkg.Desc,
			Module:      pkg.Module,
			IsEmpty:     false,
		})
	}

	var dirtyHistoryIDs []uint
	for _, history := range histories {
		for _, emptyID := range emptyPackageHistoryIDs {
			if history.ID == emptyID {
				dirtyHistoryIDs = append(dirtyHistoryIDs, history.ID)
				break
			}
		}
	}

	if len(dirtyHistoryIDs) > 0 {
		deletedCount := 0
		for _, historyID := range dirtyHistoryIDs {
			if err := deleteAutoCodeHistory(ctx, historyID); err != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("deletedirtyHistoryRecordfailed: %v", err))
				continue
			}
			deletedCount++
		}
		if deletedCount > 0 {
			global.GVA_LOG.Info(fmt.Sprintf("succeededdelete %d RowdirtyHistoryRecord", deletedCount))
		}

		if err := g.cleanupRelatedApiAndMenus(dirtyHistoryIDs); err != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("CleanRelatedAPIAndMenuRecordfailed: %v", err))
		}
	}

	predesignedModules, err := g.scanPredesignedModules()
	if err != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("scan predesigned modulesfailed: %v", err))
		predesignedModules = []PredesignedModuleInfo{}
	}

	filteredModules := []PredesignedModuleInfo{}
	for _, module := range predesignedModules {
		isDeleted := false
		for _, deletedPkg := range cleanupInfo.DeletedPackages {
			if module.PackageName == deletedPkg {
				isDeleted = true
				break
			}
		}
		if !isDeleted {
			filteredModules = append(filteredModules, module)
		}
	}

	dictionaries := []DictionaryPre{}
	dictEntities, err := fetchDictionaryList(ctx, "")
	if err != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("get dictionaryInformationfailed: %v", err))
	} else {
		for _, dictionary := range dictEntities {
			dictionaries = append(dictionaries, DictionaryPre{
				Type: dictionary.Type,
				Desc: dictionary.Desc,
			})
		}
	}

	var cleanupResult *CleanupInfo
	if len(cleanupInfo.DeletedPackages) > 0 || len(cleanupInfo.DeletedModules) > 0 {
		var message strings.Builder
		message.WriteString("**SystemCleanComplete**\n\n")
		if len(cleanupInfo.DeletedPackages) > 0 {
message.WriteString(fmt.Sprintf("- deleteDone %d PieceEmptyPackage: %s\n", len(cleanupInfo.DeletedPackages), strings.Join(cleanupInfo.DeletedPackages, ", ")))
		}
		if len(cleanupInfo.DeletedModules) > 0 {
message.WriteString(fmt.Sprintf("- deleteDone %d PieceRelatedModule: %s\n", len(cleanupInfo.DeletedModules), strings.Join(cleanupInfo.DeletedModules, ", ")))
		}
		cleanupInfo.CleanupMessage = message.String()
		cleanupResult = cleanupInfo
	}

	response := &AnalyzeResponse{
		ExistingPackages:   validPackages,
		PredesignedModules: filteredModules,
		Dictionaries:       dictionaries,
		CleanupInfo:        cleanupResult,
	}

	return response, nil
}

// isPackageFolderEmpty checks whether the package folder is empty
func (g *GVAAnalyzer) isPackageFolderEmpty(packageName, template string) (bool, error) {
	// resolve the base path from the template type
	var basePath string
	if template == "plugin" {
		basePath = filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server, "plugin", packageName)
	} else {
		basePath = filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server, "api", "v1", packageName)
	}

	// check whether the folder exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return true, nil // folder does not exist; treat as empty
	} else if err != nil {
		return false, err // other errors
	}
	// recursively check for .go files
	return g.hasGoFilesRecursive(basePath)
}

// hasGoFilesRecursive recursively checks for .go files under the directory and subdirectories
func (g *GVAAnalyzer) hasGoFilesRecursive(dirPath string) (bool, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return true, err // read failed; return empty
	}

	// check .go files in the current directory
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			return false, nil // found .go files; not empty
		}
	}

	// recursively check subdirectories
	for _, entry := range entries {
		if entry.IsDir() {
			subDirPath := filepath.Join(dirPath, entry.Name())
			isEmpty, err := g.hasGoFilesRecursive(subDirPath)
			if err != nil {
				continue // ignore subdirectory errors and continue checking other directories
			}
			if !isEmpty {
				return false, nil // found .go files in a subdirectory, not empty
			}
		}
	}

	return true, nil // no .go files found; folder is empty
}

// removeEmptyPackageFolder deleteempty package folder
func (g *GVAAnalyzer) removeEmptyPackageFolder(packageName, template string) error {
	var basePath string
	if template == "plugin" {
		basePath = filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server, "plugin", packageName)
	} else {
		// for package type, multiple directories must be removed
		paths := []string{
			filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server, "api", "v1", packageName),
			filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server, "model", packageName),
			filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server, "router", packageName),
			filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server, "service", packageName),
		}
		for _, path := range paths {
			if err := g.removeDirectoryIfExists(path); err != nil {
				return err
			}
		}
		return nil
	}

	return g.removeDirectoryIfExists(basePath)
}

// removeDirectoryIfExists remove directory if it exists
func (g *GVAAnalyzer) removeDirectoryIfExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil // directory does not exist; nothing to remove
	} else if err != nil {
		return err // other errors
	}

	// check whether the directory contains Go files
	noGoFiles, err := g.hasGoFilesRecursive(dirPath)
	if err != nil {
		return err
	}
	// hasGoFilesRecursive returns false when Go files are found
	if noGoFiles {
		return os.RemoveAll(dirPath)
	}
	return nil
}

// cleanupRelatedApiAndMenus clean up related API and menu records
func (g *GVAAnalyzer) cleanupRelatedApiAndMenus(historyIDs []uint) error {
	if len(historyIDs) == 0 {
		return nil
	}

	// HereCanByAccording toNeedImplementToolBodyofAPIAndMenuCleanLogic
	// Due toInvolveToToolBodyofBusinessLogic, HereOnlyDoLogRecord
	global.GVA_LOG.Info(fmt.Sprintf("CleanHistoryRecordID %v RelatedAPIAndMenuRecord", historyIDs))

	// CanByInvokeserviceLayerofRelatedmethodPerformClean
	// e.g. service.ServiceGroupApp.SystemApiService.DeleteApisByIds(historyIDs)
	// e.g. service.ServiceGroupApp.MenuService.DeleteMenusByIds(historyIDs)

	return nil
}

// scanPredesignedModules scan predesigned modules
func (g *GVAAnalyzer) scanPredesignedModules() ([]PredesignedModuleInfo, error) {
	// getautocodeconfigurationpath
	autocodeRoot := global.GVA_CONFIG.AutoCode.Root
	if autocodeRoot == "" {
		return nil, errors.New("autocodeRootpathNotconfiguration")
	}

	var modules []PredesignedModuleInfo

	// ScanpluginDirectory
	pluginModules, err := g.scanPluginModules(filepath.Join(autocodeRoot, global.GVA_CONFIG.AutoCode.Server, "plugin"))
	if err != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("ScanpluginModulefailed: %v", err))
	} else {
		modules = append(modules, pluginModules...)
	}

	// ScanmodelDirectory
	modelModules, err := g.scanModelModules(filepath.Join(autocodeRoot, global.GVA_CONFIG.AutoCode.Server, "model"))
	if err != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("ScanmodelModulefailed: %v", err))
	} else {
		modules = append(modules, modelModules...)
	}

	return modules, nil
}

// scanPluginModules ScanPluginModule
func (g *GVAAnalyzer) scanPluginModules(pluginDir string) ([]PredesignedModuleInfo, error) {
	var modules []PredesignedModuleInfo

	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		return modules, nil // directory missing; return empty list
	}

	entries, err := os.ReadDir(pluginDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			pluginName := entry.Name()
			pluginPath := filepath.Join(pluginDir, pluginName)

			// LookupmodelDirectory
			modelDir := filepath.Join(pluginPath, "model")
			if _, err := os.Stat(modelDir); err == nil {
				// ScanmodelDirectoryDown module
				pluginModules, err := g.scanModulesInDirectory(modelDir, pluginName, "plugin")
				if err != nil {
					global.GVA_LOG.Warn(fmt.Sprintf("ScanPlugin %s  modulefailed: %v", pluginName, err))
					continue
				}
				modules = append(modules, pluginModules...)
			}
		}
	}

	return modules, nil
}

// scanModelModules ScanModelModule
func (g *GVAAnalyzer) scanModelModules(modelDir string) ([]PredesignedModuleInfo, error) {
	var modules []PredesignedModuleInfo

	if _, err := os.Stat(modelDir); os.IsNotExist(err) {
		return modules, nil // directory missing; return empty list
	}

	entries, err := os.ReadDir(modelDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			packageName := entry.Name()
			packagePath := filepath.Join(modelDir, packageName)

			// ScanPackageDirectoryDown module
			packageModules, err := g.scanModulesInDirectory(packagePath, packageName, "package")
			if err != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("ScanPackage %s  modulefailed: %v", packageName, err))
				continue
			}
			modules = append(modules, packageModules...)
		}
	}

	return modules, nil
}

// scanModulesInDirectory Scanin directory module
func (g *GVAAnalyzer) scanModulesInDirectory(dir, packageName, template string) ([]PredesignedModuleInfo, error) {
	var modules []PredesignedModuleInfo

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			moduleName := strings.TrimSuffix(entry.Name(), ".go")
			filePath := filepath.Join(dir, entry.Name())

			module := PredesignedModuleInfo{
				ModuleName:  moduleName,
				PackageName: packageName,
				Template:    template,
				FilePaths:   []string{filePath},
				Description: fmt.Sprintf("%sModuleIn%s", packageName, moduleName),
			}
			modules = append(modules, module)
		}
	}

	return modules, nil
}
