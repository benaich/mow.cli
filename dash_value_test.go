package cli

import (
	"flag"
	"io"
	"strings"
	"testing"
)

func TestStandaloneDashValues(t *testing.T) {
	for _, args := range [][]string{
		{"--source", "-", "--verbose", "destination"},
		{"-s", "-", "-v", "destination"},
		{"-vs", "-", "destination"},
		{"--source=-", "destination"},
		{"-s=-", "destination"},
		{"-s-", "destination"},
		{"--source", "-", "-"},
		{"--source", "-", "--", "--literal"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			app := App("test", "")
			app.ErrorHandling = flag.ContinueOnError
			app.Output = io.Discard
			called := false
			app.Command("copy", "", func(cmd *Cmd) {
				source := cmd.StringOpt("s source", "", "")
				cmd.BoolOpt("v verbose", false, "")
				dest := cmd.StringArg("DEST", "", "")
				cmd.Spec = "[OPTIONS] [-- ] DEST"
				cmd.Action = func() {
					called = true
					if *source != "-" || *dest != args[len(args)-1] {
						t.Fatalf("source=%q destination=%q", *source, *dest)
					}
				}
			})
			if err := app.Run(append([]string{"test", "copy"}, args...)); err != nil || !called {
				t.Fatalf("args=%q called=%v err=%v", args, called, err)
			}
		})
	}
}

func TestStandaloneDashDoesNotConsumeOptionsAsValues(t *testing.T) {
	for _, args := range [][]string{
		{"--source", "--verbose"}, {"-s", "-v"},
		{"--source", "--"}, {"-s", "--"},
		{"--source", "-unknown"}, {"--unknown", "-"},
	} {
		app := App("test", "")
		app.ErrorHandling = flag.ContinueOnError
		app.Output = io.Discard
		app.StringOpt("s source", "", "")
		app.BoolOpt("v verbose", false, "")
		app.Spec = "[OPTIONS]"
		app.Action = func() { t.Error("invalid arguments reached action") }
		if err := app.Run(append([]string{"test"}, args...)); err == nil {
			t.Errorf("invalid arguments accepted: %q", args)
		}
	}
}
