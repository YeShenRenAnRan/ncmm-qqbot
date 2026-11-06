package handler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testLog(t *testing.T, home string, lines []string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(home, "log"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, "log", "ncm.log"), []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestReadNcmLogSummary(t *testing.T) {
	home := t.TempDir()
	testLog(t, home, []string{
		`time=2026-11-06T00:00:00+08:00 level=INFO msg="[playids] 当前账号：uid=1915065397 昵称=\"nnylo_\"" app=ncm`,
		`time=2026-11-06T00:00:01+08:00 level=INFO msg="[playids] >>>>>> 开始为账号 (/data/cookie.json) 执行模拟播放 <<<<<<" app=ncm`,
		`time=2026-11-06T00:05:00+08:00 level=INFO msg="[playids] 本首结果：第7/319首，成功" app=ncm`,
	})
	s := readNcmLogSummary(home)
	if s.UID != "1915065397" || s.Nickname != "nnylo_" {
		t.Fatalf("account parse: %+v", s)
	}
	if s.StartAt != "2026-11-06T00:00:01+08:00" || s.Progress != "7/319" {
		t.Fatalf("start/progress: %+v", s)
	}
	if s.ErrSample != "" {
		t.Fatalf("unexpected error: %q", s.ErrSample)
	}
}

func TestNcmLogErrorResetsEachRun(t *testing.T) {
	home := t.TempDir()
	testLog(t, home, []string{
		`time=2026-11-05T00:00:00+08:00 level=ERROR msg="未登录" app=ncm`,
		`time=2026-11-06T00:00:00+08:00 level=INFO msg="[playids] >>>>>> 开始为账号 (x) 执行模拟播放 <<<<<<" app=ncm`,
		`time=2026-11-06T00:05:00+08:00 level=INFO msg="[playids] 本首结果：第1/319首，成功" app=ncm`,
	})
	if s := readNcmLogSummary(home); s.ErrSample != "" {
		t.Fatalf("old error leaked into new run: %q", s.ErrSample)
	}
}
