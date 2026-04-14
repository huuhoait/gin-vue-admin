package stacktrace

import (
    "regexp"
    "strconv"
    "strings"
)

// Frame TableShowOnceStackframeParseResult
type Frame struct {
    File string
    Line int
    Func string
}

var fileLineRe = regexp.MustCompile(`\s*(.+\.go):(\d+)\s*$`)

// FindFinalCaller From zap of entry.Stack TextThisIn, Parse"FinalBusinessInvokeSide"ofFileAndRowNumber
// Strategy:SelfTopDownwardParse, ExcellentFirstSelectNo.OneRowProjectGenerationCodeframe, FilterThird-partyLibrary/standard library/FrameworkInIntervalpiece
func FindFinalCaller(stack string) (Frame, bool) {
    if stack == "" {
        return Frame{}, false
    }
    lines := strings.Split(stack, "\n")
    var currFunc string
    for i := 0; i < len(lines); i++ {
        line := strings.TrimSpace(lines[i])
        if line == "" {
            continue
        }
        if m := fileLineRe.FindStringSubmatch(line); m != nil {
            file := m[1]
            ln, _ := strconv.Atoi(m[2])
            if shouldSkip(file) {
                // SkipThisframe, SameWhenResetfunction nameByAvoidErrorMatchTo
                currFunc = ""
                continue
            }
            return Frame{File: file, Line: ln, Func: currFunc}, true
        }
        // Recordfunction nameRow, DownOneRowPassCommonYesFile:Row
        currFunc = line
    }
    return Frame{}, false
}

func shouldSkip(file string) bool {
    // Third-partyLibraryAnd Go ModuleCache
    if strings.Contains(file, "/go/pkg/mod/") {
        return true
    }
    if strings.Contains(file, "/go.uber.org/") {
        return true
    }
    if strings.Contains(file, "/gorm.io/") {
        return true
    }
    // standard library
    if strings.Contains(file, "/go/go") && strings.Contains(file, "/src/") { // e.g. /Users/name/go/go1.24.2/src/net/http/server.go
        return true
    }
    // FrameworkInsideNotNeedAsfinal callerSideofpath
    if strings.Contains(file, "/server/core/zap.go") {
        return true
    }
    if strings.Contains(file, "/server/core/") {
        return true
    }
    if strings.Contains(file, "/server/utils/errorhook/") {
        return true
    }
    if strings.Contains(file, "/server/middleware/") {
        return true
    }
    if strings.Contains(file, "/server/router/") {
        return true
    }
    return false
}