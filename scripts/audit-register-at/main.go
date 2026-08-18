package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	tp "github.com/tychy/toukibo-parser"
	"github.com/tychy/toukibo-parser/internal/toukibo"
)

var (
	warekiDateRe = regexp.MustCompile(`(?:明治|大正|昭和|平成|令和)[　 ]*(?:元|[０-９0-9]+)年[　 ]*[０-９0-9]+月[　 ]*[０-９0-9]+日`)
	registerRe   = regexp.MustCompile(`((?:明治|大正|昭和|平成|令和)[　 ]*(?:元|[０-９0-9]+)年[　 ]*[０-９0-9]+月[　 ]*[０-９0-9]+日)[　 ]*登記`)
	appointRe    = regexp.MustCompile(`((?:明治|大正|昭和|平成|令和)[　 ]*(?:元|[０-９0-9]+)年[　 ]*[０-９0-9]+月[　 ]*[０-９0-9]+日)[　 ]*(就任|重任)`)
	otherDateRe  = regexp.MustCompile(`((?:明治|大正|昭和|平成|令和)[　 ]*(?:元|[０-９0-9]+)年[　 ]*[０-９0-9]+月[　 ]*[０-９0-9]+日)[　 ]*(辞任|退任|死亡|抹消|廃止|解任|退社|移記|更正|資格変更|責任変更)`)
	splitMarkRe  = regexp.MustCompile(`├[─－]+`)
	positionRe   = regexp.MustCompile(`(?:代表取締役|取締役・監査等|取締役|監査役|会計監査人|代表理事|理事長|理事|監事|代表社員|業務執行社員|会長|代表清算人|清算人|代表役員|会計参与|無限責任社員|有限責任社員|破産管財人|評議員|代表者|会頭|学長|代表執行役|執行役|報酬委員|監査委員|指名委員|職務執行者|保全管財人|社員)`)
	yamlNameRe   = regexp.MustCompile(`^\s+- Name:\s*(.*)$`)
	yamlPosRe    = regexp.MustCompile(`^\s+Position:\s*(.*)$`)
	yamlRegRe    = regexp.MustCompile(`^\s+RegisterAt:`)
)

type finding struct {
	Sample   int
	Name     string
	Position string
	Got      string
	Kind     string
	Detail   string
}

func normalizeDate(s string) string {
	s = strings.ReplaceAll(s, "　", "")
	s = strings.ReplaceAll(s, " ", "")
	return toukibo.ZenkakuToHankaku(s)
}

func dateSet(re *regexp.Regexp, text string) map[string]string {
	out := map[string]string{}
	for _, m := range re.FindAllStringSubmatch(text, -1) {
		raw := strings.ReplaceAll(m[1], "　", "")
		raw = strings.ReplaceAll(raw, " ", "")
		out[normalizeDate(m[1])] = raw
	}
	return out
}

func executiveSection(content string) string {
	start := strings.Index(content, "役員に関する事項")
	if start == -1 {
		start = strings.Index(content, "社員に関する事項")
	}
	if start == -1 {
		return ""
	}
	rest := content[start:]
	for _, mark := range []string{"登記記録に関する", "┗━━━━━━━━"} {
		if i := strings.Index(rest[20:], mark); i >= 0 {
			return rest[:20+i]
		}
	}
	if len(rest) > 20000 {
		return rest[:20000]
	}
	return rest
}

func nameRe(name string) *regexp.Regexp {
	var b strings.Builder
	for _, r := range name {
		b.WriteString(regexp.QuoteMeta(string(r)))
		b.WriteString(`[　 ]*`)
	}
	return regexp.MustCompile(b.String())
}

func windowsForName(section, name string) []string {
	if name == "" {
		return nil
	}
	re := nameRe(name)
	idxs := re.FindAllStringIndex(section, -1)
	if len(idxs) == 0 {
		return nil
	}
	var windows []string
	for _, idx := range idxs {
		start := idx[0]
		end := len(section)
		if loc := splitMarkRe.FindStringIndex(section[idx[1]:]); loc != nil {
			if idx[1]+loc[1] < end {
				end = idx[1] + loc[1]
			}
		}
		if loc := positionRe.FindStringIndex(section[idx[1]:]); loc != nil {
			// keep a bit of the next row's right column; cut before the next role label
			cut := idx[1] + loc[0]
			if cut > start && cut < end {
				end = cut
			}
		}
		if end-start > 800 {
			end = start + 800
		}
		windows = append(windows, section[start:end])
	}
	return windows
}

func nearbyRegister(section, name string) (string, bool) {
	for _, w := range windowsForName(section, name) {
		ms := registerRe.FindAllStringSubmatch(w, -1)
		if len(ms) == 0 {
			continue
		}
		raw := strings.ReplaceAll(ms[len(ms)-1][1], "　", "")
		raw = strings.ReplaceAll(raw, " ", "")
		return raw, true
	}
	return "", false
}

func patchYAML(path string, execs []toukibo.HoujinExecutiveValue) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	out := make([]string, 0, len(lines)+len(execs))
	inExec := false
	execIdx := -1
	var names []string
	var positions []string
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "HoujinExecutiveValues:") {
			inExec = true
			execIdx = -1
			out = append(out, line)
			continue
		}
		if inExec && line != "" && line[0] != ' ' && line[0] != '\t' {
			inExec = false
		}
		if inExec {
			if m := yamlNameRe.FindStringSubmatch(line); m != nil {
				execIdx++
				names = append(names, strings.TrimSpace(m[1]))
			}
			if m := yamlPosRe.FindStringSubmatch(line); m != nil {
				positions = append(positions, strings.TrimSpace(m[1]))
				out = append(out, line)
				register := ""
				if execIdx >= 0 && execIdx < len(execs) {
					register = execs[execIdx].RegisterAt
				}
				if i+1 < len(lines) && yamlRegRe.MatchString(lines[i+1]) {
					i++
				}
				out = append(out, "    RegisterAt: "+register)
				continue
			}
			if yamlRegRe.MatchString(line) {
				continue
			}
		}
		out = append(out, line)
	}

	lookup := map[string]string{}
	for _, ev := range execs {
		key := ev.Name + "\t" + ev.Position
		lookup[key] = ev.RegisterAt
	}
	used := map[string]int{}
	rewritten := make([]string, 0, len(out)+len(names))
	execIdx = -1
	inExec = false
	for i := 0; i < len(out); i++ {
		line := out[i]
		if strings.HasPrefix(line, "HoujinExecutiveValues:") {
			inExec = true
			rewritten = append(rewritten, line)
			continue
		}
		if inExec && line != "" && line[0] != ' ' && line[0] != '\t' {
			inExec = false
		}
		if inExec && yamlRegRe.MatchString(line) {
			execIdx++
			register := ""
			if execIdx < len(names) && execIdx < len(positions) {
				key := names[execIdx] + "\t" + positions[execIdx]
				if execIdx < len(execs) && execs[execIdx].Name == names[execIdx] && execs[execIdx].Position == positions[execIdx] {
					register = execs[execIdx].RegisterAt
				} else if v, ok := lookup[key]; ok {
					register = v
				}
				used[key]++
			}
			rewritten = append(rewritten, "    RegisterAt: "+register)
			continue
		}
		rewritten = append(rewritten, line)
	}

	text := strings.Join(rewritten, "\n")
	return os.WriteFile(path, []byte(text), 0644)
}

func auditOne(i int) (stats map[string]int, findings []finding, err error) {
	stats = map[string]int{}
	pdfPath := fmt.Sprintf("testdata/pdf/sample%d.pdf", i)
	yamlPath := fmt.Sprintf("testdata/yaml/sample%d.yaml", i)
	h, err := tp.ParseByPDFPath(pdfPath)
	if err != nil {
		return stats, nil, fmt.Errorf("parse: %w", err)
	}
	execs, err := h.GetHoujinExecutives()
	if err != nil {
		if !strings.Contains(err.Error(), "not found executives") {
			return stats, nil, fmt.Errorf("execs: %w", err)
		}
		stats["no_executives"] = 1
		execs = nil
	}
	if err := patchYAML(yamlPath, execs); err != nil {
		return stats, nil, fmt.Errorf("yaml: %w", err)
	}
	if len(execs) == 0 {
		return stats, nil, nil
	}

	content, err := tp.GetContentByPDFPath(pdfPath)
	if err != nil {
		return stats, nil, fmt.Errorf("content: %w", err)
	}
	section := executiveSection(content)
	sectionRegisters := dateSet(registerRe, section)
	sectionAppoints := dateSet(appointRe, section)
	sectionOthers := dateSet(otherDateRe, section)

	stats["executives"] = len(execs)
	for _, ev := range execs {
		if ev.RegisterAt == "" {
			stats["empty"]++
			if nearby, ok := nearbyRegister(section, ev.Name); ok {
				stats["empty_but_nearby_register"]++
				findings = append(findings, finding{
					Sample: i, Name: ev.Name, Position: ev.Position,
					Kind: "empty_but_nearby_register", Got: nearby,
					Detail: "parser empty, name window has " + nearby + "登記",
				})
			}
			continue
		}
		stats["filled"]++
		norm := normalizeDate(ev.RegisterAt)
		if !warekiDateRe.MatchString(ev.RegisterAt) && !warekiDateRe.MatchString(toukibo.ZenkakuToHankaku(ev.RegisterAt)) {
			stats["invalid_format"]++
			findings = append(findings, finding{
				Sample: i, Name: ev.Name, Position: ev.Position,
				Kind: "invalid_format", Got: ev.RegisterAt,
			})
		}
		if _, ok := sectionRegisters[norm]; !ok {
			stats["not_in_exec_section"]++
			findings = append(findings, finding{
				Sample: i, Name: ev.Name, Position: ev.Position,
				Kind: "not_in_exec_section", Got: ev.RegisterAt,
				Detail: "date not found as 登記 in executive section",
			})
		}
		if _, isAppoint := sectionAppoints[norm]; isAppoint {
			if _, isReg := sectionRegisters[norm]; !isReg {
				stats["appointment_as_register"]++
				findings = append(findings, finding{
					Sample: i, Name: ev.Name, Position: ev.Position,
					Kind: "appointment_as_register", Got: ev.RegisterAt,
					Detail: "date exists as 就任/重任 but not as 登記",
				})
			}
		}
		if _, isOther := sectionOthers[norm]; isOther {
			if _, isReg := sectionRegisters[norm]; !isReg {
				stats["other_event_as_register"]++
				findings = append(findings, finding{
					Sample: i, Name: ev.Name, Position: ev.Position,
					Kind: "other_event_as_register", Got: ev.RegisterAt,
					Detail: "date exists as 辞任/退任/死亡/移記/更正 etc but not as 登記",
				})
			}
		}
		if nearby, ok := nearbyRegister(section, ev.Name); ok && normalizeDate(nearby) != norm {
			stats["nearby_differs"]++
			findings = append(findings, finding{
				Sample: i, Name: ev.Name, Position: ev.Position,
				Kind: "nearby_differs", Got: ev.RegisterAt,
				Detail: "name window last 登記 is " + nearby,
			})
		}
	}
	return stats, findings, nil
}

func sampleCount() int {
	matches, err := filepath.Glob("testdata/pdf/sample*.pdf")
	if err != nil {
		return 0
	}
	max := 0
	for _, m := range matches {
		base := strings.TrimSuffix(filepath.Base(m), ".pdf")
		n, err := strconv.Atoi(strings.TrimPrefix(base, "sample"))
		if err == nil && n > max {
			max = n
		}
	}
	return max
}

func main() {
	n := sampleCount()
	if n == 0 {
		fmt.Fprintln(os.Stderr, "no sample PDFs; run make get/sample")
		os.Exit(1)
	}

	workers := runtime.GOMAXPROCS(0)
	if workers < 2 {
		workers = 2
	}
	type result struct {
		i        int
		stats    map[string]int
		findings []finding
		err      error
	}
	jobs := make(chan int)
	out := make(chan result)
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				st, fs, err := auditOne(i)
				out <- result{i: i, stats: st, findings: fs, err: err}
			}
		}()
	}
	go func() {
		for i := 1; i <= n; i++ {
			if _, err := os.Stat(fmt.Sprintf("testdata/pdf/sample%d.pdf", i)); err != nil {
				continue
			}
			jobs <- i
		}
		close(jobs)
		wg.Wait()
		close(out)
	}()

	totals := map[string]int{}
	var findings []finding
	var errors []string
	var done atomic.Int64
	for r := range out {
		done.Add(1)
		if r.err != nil {
			errors = append(errors, fmt.Sprintf("sample%d: %v", r.i, r.err))
			continue
		}
		for k, v := range r.stats {
			totals[k] += v
		}
		findings = append(findings, r.findings...)
		if d := done.Load(); d%200 == 0 {
			fmt.Fprintf(os.Stderr, "processed %d\n", d)
		}
	}

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Kind != findings[j].Kind {
			return findings[i].Kind < findings[j].Kind
		}
		if findings[i].Sample != findings[j].Sample {
			return findings[i].Sample < findings[j].Sample
		}
		return findings[i].Name < findings[j].Name
	})
	sort.Strings(errors)

	fmt.Println("=== totals ===")
	keys := make([]string, 0, len(totals))
	for k := range totals {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s: %d\n", k, totals[k])
	}
	fmt.Printf("yaml_errors: %d\n", len(errors))
	fmt.Printf("findings: %d\n", len(findings))

	if len(errors) > 0 {
		fmt.Println("\n=== errors ===")
		for _, e := range errors {
			fmt.Println(e)
		}
	}

	fmt.Println("\n=== findings ===")
	byKind := map[string][]finding{}
	for _, f := range findings {
		byKind[f.Kind] = append(byKind[f.Kind], f)
	}
	kinds := make([]string, 0, len(byKind))
	for k := range byKind {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)
	for _, k := range kinds {
		fs := byKind[k]
		fmt.Printf("\n-- %s (%d) --\n", k, len(fs))
		limit := len(fs)
		if limit > 80 {
			limit = 80
		}
		for _, f := range fs[:limit] {
			fmt.Printf("sample%d\t%s\t%s\tgot=%q\t%s\n", f.Sample, f.Name, f.Position, f.Got, f.Detail)
		}
		if len(fs) > limit {
			fmt.Printf("... %d more\n", len(fs)-limit)
		}
	}
}
