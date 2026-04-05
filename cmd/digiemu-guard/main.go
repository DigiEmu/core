package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Severity string
type Zone string

const (
	SeverityInfo     Severity = "INFO"
	SeverityWarning  Severity = "WARNING"
	SeverityCritical Severity = "CRITICAL"
)

const (
	ZoneCore     Zone = "core"
	ZoneBoundary Zone = "boundary"
	ZoneSupport  Zone = "support"
)

type Finding struct {
	Severity Severity `json:"severity"`
	File     string   `json:"file"`
	Line     int      `json:"line"`
	Column   int      `json:"column"`
	Rule     string   `json:"rule"`
	Category string   `json:"category"`
	Message  string   `json:"message"`
	Zone     Zone     `json:"zone"`
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s:%d:%d [%s/%s/%s] %s",
		f.Severity, f.File, f.Line, f.Column, f.Zone, f.Category, f.Rule, f.Message)
}

type IgnoreRule struct {
	FileContains string `json:"file_contains"`
	Rule         string `json:"rule"`
	Category     string `json:"category"`
}

type IgnoreConfig struct {
	Ignores []IgnoreRule `json:"ignores"`
}

type Config struct {
	Root          string
	FailOnWarning bool
	Verbose       bool
	JSON          bool
	CoreOnly      bool
	IgnoreFile    string
	Ignores       IgnoreConfig
}

type FileContext struct {
	File          *ast.File
	FilePath      string
	FileSet       *token.FileSet
	Info          *types.Info
	ImportAliases map[string]string
	Zone          Zone
}

func main() {
	cfg, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "guard config error: %v\n", err)
		os.Exit(2)
	}

	findings, err := scanRepo(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "guard failed: %v\n", err)
		os.Exit(2)
	}

	findings = filterIgnored(findings, cfg.Ignores)

	if findings == nil {
		findings = []Finding{}
	}

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].File != findings[j].File {
			return findings[i].File < findings[j].File
		}
		if findings[i].Line != findings[j].Line {
			return findings[i].Line < findings[j].Line
		}
		if findings[i].Column != findings[j].Column {
			return findings[i].Column < findings[j].Column
		}
		if findings[i].Severity != findings[j].Severity {
			return findings[i].Severity < findings[j].Severity
		}
		if findings[i].Category != findings[j].Category {
			return findings[i].Category < findings[j].Category
		}
		return findings[i].Rule < findings[j].Rule
	})

	if cfg.JSON {
		out, _ := json.MarshalIndent(findings, "", "  ")
		fmt.Println(string(out))
	} else {
		if len(findings) == 0 {
			fmt.Println("DigiEmu Guard v3: no findings.")
		} else {
			for _, f := range findings {
				fmt.Println(f.String())
			}
			fmt.Println()
			printSummary(findings)
		}
	}

	criticals, warnings := countSeverities(findings)
	if criticals > 0 {
		os.Exit(1)
	}
	if cfg.FailOnWarning && warnings > 0 {
		os.Exit(1)
	}
}

func parseFlags() (Config, error) {
	root := flag.String("root", ".", "repository root to scan")
	failOnWarning := flag.Bool("fail-on-warning", false, "exit non-zero on warnings too")
	verbose := flag.Bool("v", false, "verbose output")
	jsonOut := flag.Bool("json", false, "emit findings as JSON")
	coreOnly := flag.Bool("core-only", false, "scan only core-sensitive files")
	ignoreFile := flag.String("ignore-file", ".digiemu-guard-ignore.json", "path to ignore config JSON")
	flag.Parse()

	cfg := Config{
		Root:          *root,
		FailOnWarning: *failOnWarning,
		Verbose:       *verbose,
		JSON:          *jsonOut,
		CoreOnly:      *coreOnly,
		IgnoreFile:    *ignoreFile,
	}

	if _, err := os.Stat(cfg.IgnoreFile); err == nil {
		data, err := os.ReadFile(cfg.IgnoreFile)
		if err != nil {
			return cfg, fmt.Errorf("read ignore file: %w", err)
		}
		if err := json.Unmarshal(data, &cfg.Ignores); err != nil {
			return cfg, fmt.Errorf("parse ignore file: %w", err)
		}
	}

	return cfg, nil
}

func scanRepo(cfg Config) ([]Finding, error) {
	fset := token.NewFileSet()
	findings := make([]Finding, 0)

	err := filepath.WalkDir(cfg.Root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		name := d.Name()
		if d.IsDir() {
			if shouldSkipDir(name) {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			return nil
		}

		normalized := filepath.ToSlash(path)
		if strings.Contains(normalized, "/cmd/digiemu-guard/") || strings.HasPrefix(normalized, "cmd/digiemu-guard/") {
			return nil
		}

		zone := classifyZone(path)
		if cfg.CoreOnly && zone != ZoneCore {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			findings = append(findings, makeFinding(
				token.Position{Filename: path, Line: 1, Column: 1},
				SeverityCritical,
				"parse-error",
				"parser",
				ZoneSupport,
				fmt.Sprintf("could not parse file: %v", err),
			))
			return nil
		}

		ctx := buildFileContext(fset, file, path, zone)
		fileFindings := scanFile(ctx, cfg)
		findings = append(findings, fileFindings...)
		return nil
	})

	return findings, err
}

func buildFileContext(fset *token.FileSet, file *ast.File, path string, zone Zone) FileContext {
	importAliases := map[string]string{}
	for _, imp := range file.Imports {
		importPath, _ := strconv.Unquote(imp.Path.Value)
		alias := filepath.Base(importPath)
		if imp.Name != nil && imp.Name.Name != "_" && imp.Name.Name != "." {
			alias = imp.Name.Name
		}
		importAliases[alias] = importPath
	}

	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	conf := types.Config{
		Importer: importer.Default(),
		Error:    func(err error) {},
	}
	_, _ = conf.Check(file.Name.Name, fset, []*ast.File{file}, info)

	return FileContext{
		File:          file,
		FilePath:      filepath.ToSlash(path),
		FileSet:       fset,
		Info:          info,
		ImportAliases: importAliases,
		Zone:          zone,
	}
}

func shouldSkipDir(name string) bool {
	switch name {
	case ".git", ".next", "node_modules", "vendor", "dist", "build", "tmp", "coverage":
		return true
	default:
		return false
	}
}

func classifyZone(path string) Zone {
	p := filepath.ToSlash(path)
	p = strings.TrimPrefix(p, "./")
	p = "/" + strings.Trim(p, "/") + "/"

	coreParts := []string{
		"/pkg/snapshot/",
		"/pkg/replay/",
		"/pkg/verify/",
		"/pkg/meaning/",
		"/pkg/uncertainty/",
		"/internal/verify/",
		"/internal/kernel/domain/",
		"/internal/kernel/usecases/",
	}

	for _, part := range coreParts {
		if strings.Contains(p, part) {
			return ZoneCore
		}
	}

	boundaryParts := []string{
		"/internal/httpapi/",
		"/internal/kernel/adapters/",
	}

	for _, part := range boundaryParts {
		if strings.Contains(p, part) {
			return ZoneBoundary
		}
	}

	return ZoneSupport
}

func scanFile(ctx FileContext, cfg Config) []Finding {
	findings := make([]Finding, 0)

	findings = append(findings, scanImports(ctx)...)
	findings = append(findings, scanDecls(ctx)...)

	ast.Inspect(ctx.File, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GoStmt:
			findings = append(findings, zoneFinding(ctx, x.Pos(),
				severityByZone(ctx.Zone, SeverityCritical, SeverityWarning, SeverityWarning),
				"concurrency",
				"goroutine",
				"goroutine usage found; concurrency in deterministic paths must be explicitly isolated and ordered",
			))

		case *ast.SelectStmt:
			findings = append(findings, zoneFinding(ctx, x.Pos(),
				severityByZone(ctx.Zone, SeverityCritical, SeverityWarning, SeverityWarning),
				"concurrency",
				"select-review",
				"select statement found; verify no scheduling-dependent behavior exists",
			))

		case *ast.RangeStmt:
			findings = append(findings, scanRangeStmt(ctx, x)...)

		case *ast.CallExpr:
			findings = append(findings, scanCallExpr(ctx, x)...)

		case *ast.Ident:
			findings = append(findings, scanFloatIdent(ctx, x)...)
		}
		return true
	})

	if cfg.Verbose && len(findings) == 0 {
		findings = append(findings, Finding{
			Severity: SeverityInfo,
			File:     ctx.FilePath,
			Line:     1,
			Column:   1,
			Rule:     "scanned",
			Category: "meta",
			Message:  "file scanned with no findings",
			Zone:     ctx.Zone,
		})
	}

	return findings
}

func scanImports(ctx FileContext) []Finding {
	findings := make([]Finding, 0)

	for _, imp := range ctx.File.Imports {
		importPath, _ := strconv.Unquote(imp.Path.Value)

		switch importPath {
		case "unsafe":
			findings = append(findings, zoneFinding(ctx, imp.Pos(),
				SeverityCritical,
				"hidden-state",
				"unsafe-import",
				`import of "unsafe" is forbidden`,
			))

		case "sync", "sync/atomic":
			findings = append(findings, zoneFinding(ctx, imp.Pos(),
				severityByZone(ctx.Zone, SeverityCritical, SeverityWarning, SeverityInfo),
				"concurrency",
				"sync-review",
				fmt.Sprintf("import of %q found; verify no nondeterministic concurrency affects deterministic behavior", importPath),
			))
		}
	}

	return findings
}

func scanDecls(ctx FileContext) []Finding {
	findings := make([]Finding, 0)

	for _, decl := range ctx.File.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					findings = append(findings, scanTypeSpec(ctx, s)...)
				case *ast.ValueSpec:
					findings = append(findings, scanValueSpec(ctx, d, s)...)
				}
			}

		case *ast.FuncDecl:
			findings = append(findings, scanFuncDecl(ctx, d)...)
		}
	}

	return findings
}

func scanTypeSpec(ctx FileContext, spec *ast.TypeSpec) []Finding {
	st, ok := spec.Type.(*ast.StructType)
	if !ok {
		return nil
	}

	findings := make([]Finding, 0)

	for _, field := range st.Fields.List {
		if field.Tag == nil {
			continue
		}

		rawTag, err := strconv.Unquote(field.Tag.Value)
		if err != nil {
			continue
		}

		if strings.Contains(rawTag, `json:"-"`) {
			findings = append(findings, zoneFinding(ctx, field.Pos(),
				severityByZone(ctx.Zone, SeverityCritical, SeverityWarning, SeverityWarning),
				"serialization",
				"ghost-field",
				fmt.Sprintf("struct %q contains json:\"-\"; hidden fields must not affect deterministic behavior", spec.Name.Name),
			))
		}

		if strings.Contains(rawTag, "omitempty") {
			if shouldFlagOmitEmpty(ctx, spec.Name.Name, field) {
				findings = append(findings, zoneFinding(ctx, field.Pos(),
					severityByZone(ctx.Zone, SeverityWarning, SeverityWarning, SeverityInfo),
					"serialization",
					"omitempty-review",
					fmt.Sprintf("struct %q contains omitempty; verify no semantic state can disappear from snapshots or verification state", spec.Name.Name),
				))
			}
		}
	}

	return findings
}

func scanValueSpec(ctx FileContext, decl *ast.GenDecl, spec *ast.ValueSpec) []Finding {
	findings := make([]Finding, 0)

	if decl.Tok == token.VAR && !isLikelySentinelError(ctx, spec) && !isAllowedBuildMetadata(spec) {
		for _, name := range spec.Names {
			if name.Name == "_" {
				continue
			}

			findings = append(findings, zoneFinding(ctx, name.Pos(),
				severityByZone(ctx.Zone, SeverityCritical, SeverityWarning, SeverityInfo),
				"hidden-state",
				"global-var",
				fmt.Sprintf("package-level var %q found; verify it does not create hidden mutable state", name.Name),
			))
		}
	}

	for _, name := range spec.Names {
		obj := ctx.Info.Defs[name]
		if obj == nil {
			continue
		}
		if isFloatType(obj.Type()) {
			findings = append(findings, zoneFinding(ctx, name.Pos(),
				severityByZone(ctx.Zone, SeverityCritical, SeverityWarning, SeverityInfo),
				"numeric",
				"float-declaration",
				fmt.Sprintf("float declaration %q found; fixed-point/integer types are preferred in deterministic logic", name.Name),
			))
		}
	}

	return findings
}

func scanFuncDecl(ctx FileContext, fn *ast.FuncDecl) []Finding {
	findings := make([]Finding, 0)

	if fn.Recv == nil && fn.Name != nil && fn.Name.Name == "init" {
		findings = append(findings, zoneFinding(ctx, fn.Pos(),
			severityByZone(ctx.Zone, SeverityCritical, SeverityWarning, SeverityInfo),
			"hidden-state",
			"init-function",
			"init() function found; verify it does not introduce hidden state or environment-dependent behavior",
		))
	}

	if fn.Type != nil && fn.Type.Params != nil {
		for _, field := range fn.Type.Params.List {
			if exprType := ctx.Info.TypeOf(field.Type); isFloatType(exprType) {
				findings = append(findings, zoneFinding(ctx, field.Pos(),
					severityByZone(ctx.Zone, SeverityCritical, SeverityWarning, SeverityInfo),
					"numeric",
					"float-param",
					fmt.Sprintf("function %q uses float parameter type; verify deterministic thresholds and rounding", fn.Name.Name),
				))
			}
		}
	}

	if fn.Type != nil && fn.Type.Results != nil {
		for _, field := range fn.Type.Results.List {
			if exprType := ctx.Info.TypeOf(field.Type); isFloatType(exprType) {
				findings = append(findings, zoneFinding(ctx, field.Pos(),
					severityByZone(ctx.Zone, SeverityCritical, SeverityWarning, SeverityInfo),
					"numeric",
					"float-result",
					fmt.Sprintf("function %q returns float type; verify deterministic and cross-platform-safe behavior", fn.Name.Name),
				))
			}
		}
	}

	return findings
}

func scanRangeStmt(ctx FileContext, stmt *ast.RangeStmt) []Finding {
	findings := make([]Finding, 0)

	severity := severityByZone(ctx.Zone, SeverityWarning, SeverityWarning, SeverityInfo)
	msg := "range loop found; review whether this affects deterministic ordering"

	if t := ctx.Info.TypeOf(stmt.X); t != nil {
		if _, ok := deref(t).Underlying().(*types.Map); ok {
			severity = severityByZone(ctx.Zone, SeverityCritical, SeverityWarning, SeverityInfo)
			msg = "range over map found; map iteration order is nondeterministic and must not affect deterministic behavior"
		}
	}

	findings = append(findings, zoneFinding(ctx, stmt.Pos(),
		severity,
		"determinism",
		"range-review",
		msg,
	))
	return findings
}

func scanCallExpr(ctx FileContext, call *ast.CallExpr) []Finding {
	findings := make([]Finding, 0)

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return findings
	}

	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return findings
	}

	importPath, ok := ctx.ImportAliases[pkgIdent.Name]
	if !ok {
		return findings
	}

	switch importPath {
	case "time":
		switch sel.Sel.Name {
		case "Now", "Sleep", "After", "AfterFunc", "NewTicker", "NewTimer", "Since":
			findings = append(findings, zoneFinding(ctx, call.Pos(),
				severityByZone(ctx.Zone, SeverityCritical, SeverityWarning, SeverityInfo),
				"time",
				"time-usage",
				fmt.Sprintf("time.%s() found; deterministic logic must not depend on live clock behavior", sel.Sel.Name),
			))
		}

	case "math/rand", "crypto/rand":
		findings = append(findings, zoneFinding(ctx, call.Pos(),
			severityByZone(ctx.Zone, SeverityCritical, SeverityWarning, SeverityInfo),
			"randomness",
			"randomness",
			fmt.Sprintf("%s.%s() introduces randomness and is forbidden in deterministic core logic", pkgIdent.Name, sel.Sel.Name),
		))

	case "os":
		switch sel.Sel.Name {
		case "Getenv", "LookupEnv", "Environ":
			findings = append(findings, zoneFinding(ctx, call.Pos(),
				severityByZone(ctx.Zone, SeverityCritical, SeverityWarning, SeverityInfo),
				"environment",
				"environment-access",
				fmt.Sprintf("os.%s() is environment-dependent and must not affect deterministic logic", sel.Sel.Name),
			))
		}
	}

	return findings
}

func scanFloatIdent(ctx FileContext, id *ast.Ident) []Finding {
	obj := ctx.Info.Uses[id]
	if obj == nil {
		obj = ctx.Info.Defs[id]
	}
	if obj == nil || obj.Type() == nil {
		return nil
	}
	if !isFloatType(obj.Type()) {
		return nil
	}
	if id.Name == "float32" || id.Name == "float64" {
		return nil
	}

	return []Finding{
		zoneFinding(ctx, id.Pos(),
			severityByZone(ctx.Zone, SeverityCritical, SeverityWarning, SeverityInfo),
			"numeric",
			"float-usage",
			fmt.Sprintf("float-typed identifier %q used; verify deterministic rounding, thresholds, and cross-platform behavior", id.Name),
		),
	}
}
func shouldFlagOmitEmpty(ctx FileContext, _ string, field *ast.Field) bool {
	switch ctx.Zone {
	case ZoneCore:
		return true
	case ZoneBoundary:
		t := ctx.Info.TypeOf(field.Type)
		return isRiskyOmitEmptyType(t)
	default:
		return false
	}
}

func isRiskyOmitEmptyType(t types.Type) bool {
	if t == nil {
		return false
	}
	t = deref(t)

	switch u := t.Underlying().(type) {
	case *types.Basic:
		info := u.Info()
		if (info&types.IsBoolean) != 0 || (info&types.IsInteger) != 0 || (info&types.IsFloat) != 0 {
			return true
		}
	case *types.Struct:
		return true
	}
	return false
}

func isLikelySentinelError(ctx FileContext, spec *ast.ValueSpec) bool {
	if len(spec.Names) != 1 || len(spec.Values) != 1 {
		return false
	}

	name := spec.Names[0].Name
	if !strings.HasPrefix(name, "Err") {
		return false
	}

	if obj := ctx.Info.Defs[spec.Names[0]]; obj != nil {
		if obj.Type() != nil && obj.Type().String() == "error" {
			return true
		}
	}

	call, ok := spec.Values[0].(*ast.CallExpr)
	if !ok {
		return false
	}

	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		pkg, ok := fun.X.(*ast.Ident)
		if !ok {
			return false
		}
		return pkg.Name == "errors" && fun.Sel.Name == "New"
	case *ast.Ident:
		return fun.Name == "newError"
	default:
		return false
	}
}

func isAllowedBuildMetadata(spec *ast.ValueSpec) bool {
	allowed := map[string]bool{
		"Version": true,
		"Commit":  true,
		"Date":    true,
	}

	if len(spec.Names) == 0 {
		return false
	}

	for _, name := range spec.Names {
		if !allowed[name.Name] {
			return false
		}
	}
	return true
}

func isFloatType(t types.Type) bool {
	if t == nil {
		return false
	}
	b, ok := deref(t).Underlying().(*types.Basic)
	if !ok {
		return false
	}
	return (b.Info() & types.IsFloat) != 0
}

func deref(t types.Type) types.Type {
	if p, ok := t.(*types.Pointer); ok {
		return deref(p.Elem())
	}
	return t
}

func severityByZone(zone Zone, core Severity, boundary Severity, support Severity) Severity {
	switch zone {
	case ZoneCore:
		return core
	case ZoneBoundary:
		return boundary
	default:
		return support
	}
}

func zoneFinding(ctx FileContext, pos token.Pos, severity Severity, category, rule, message string) Finding {
	p := ctx.FileSet.Position(pos)
	return makeFinding(p, severity, rule, category, ctx.Zone, message)
}

func makeFinding(p token.Position, severity Severity, rule, category string, zone Zone, message string) Finding {
	return Finding{
		Severity: severity,
		File:     filepath.ToSlash(p.Filename),
		Line:     p.Line,
		Column:   p.Column,
		Rule:     rule,
		Category: category,
		Message:  message,
		Zone:     zone,
	}
}

func filterIgnored(findings []Finding, ignores IgnoreConfig) []Finding {
	if len(ignores.Ignores) == 0 {
		return findings
	}

	out := make([]Finding, 0, len(findings))
	for _, f := range findings {
		if isIgnored(f, ignores) {
			continue
		}
		out = append(out, f)
	}
	return out
}

func isIgnored(f Finding, ignores IgnoreConfig) bool {
	for _, rule := range ignores.Ignores {
		if rule.Rule != "" && rule.Rule != f.Rule {
			continue
		}
		if rule.Category != "" && rule.Category != f.Category {
			continue
		}
		if rule.FileContains != "" && !strings.Contains(filepath.ToSlash(f.File), filepath.ToSlash(rule.FileContains)) {
			continue
		}
		return true
	}
	return false
}

func countSeverities(findings []Finding) (criticals int, warnings int) {
	for _, f := range findings {
		switch f.Severity {
		case SeverityCritical:
			criticals++
		case SeverityWarning:
			warnings++
		}
	}
	return
}

func printSummary(findings []Finding) {
	criticals, warnings := countSeverities(findings)

	byCategory := map[string]int{}
	byZone := map[Zone]int{}
	for _, f := range findings {
		byCategory[f.Category]++
		byZone[f.Zone]++
	}

	fmt.Printf("Summary: %d critical, %d warning, %d info\n", criticals, warnings, countInfos(findings))

	if len(byCategory) > 0 {
		fmt.Println("By category:")
		keys := sortedKeys(byCategory)
		for _, k := range keys {
			fmt.Printf("  - %s: %d\n", k, byCategory[k])
		}
	}

	if len(byZone) > 0 {
		fmt.Println("By zone:")
		zones := []Zone{ZoneCore, ZoneBoundary, ZoneSupport}
		for _, z := range zones {
			if n, ok := byZone[z]; ok {
				fmt.Printf("  - %s: %d\n", z, n)
			}
		}
	}
}

func countInfos(findings []Finding) int {
	n := 0
	for _, f := range findings {
		if f.Severity == SeverityInfo {
			n++
		}
	}
	return n
}

func sortedKeys[K ~string, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

var _ = errors.New
