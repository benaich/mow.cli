package cli

import (
	"bytes"
	"flag"
	"testing"
)

func TestDescribeIsolatedDeclarationsAndOutput(t *testing.T) {
	app := App("example", "example app")
	app.ErrorHandling = flag.ContinueOnError
	var diagnostics bytes.Buffer
	app.Output = &diagnostics
	actions := 0
	app.Command("write", "write a resource", func(cmd *Cmd) {
		count := cmd.IntOpt("n count", 3, "item count")
		cmd.StringArg("ID", "", "resource ID")
		cmd.Spec = "[OPTIONS] ID"
		cmd.Action = func() {
			if *count != 7 {
				t.Error(*count)
			}
			actions++
		}
	})
	for i := 0; i < 2; i++ {
		tree, err := app.Cmd.Describe()
		if err != nil {
			t.Fatal(err)
		}
		child := tree.Commands[0]
		if len(child.Options) != 1 || child.Options[0].Type != "integer" || !child.Arguments[0].Required {
			t.Fatal(child)
		}
	}
	if actions != 0 {
		t.Fatal("snapshot ran an action")
	}
	if err := app.Run([]string{"example", "write", "--count", "7", "resource"}); err != nil {
		t.Fatal(err)
	}
	if actions != 1 {
		t.Fatal(actions)
	}
	if diagnostics.Len() != 0 {
		t.Fatal(diagnostics.String())
	}
}

func TestChildDiagnosticsUseAppOutput(t *testing.T) {
	app := App("example", "example app")
	app.ErrorHandling = flag.ContinueOnError
	var diagnostics bytes.Buffer
	app.Output = &diagnostics
	app.Command("child", "child command", func(cmd *Cmd) { cmd.Spec = ""; cmd.Action = func() { t.Fatal("invalid input ran action") } })
	if err := app.Run([]string{"example", "child", "--invalid"}); err == nil {
		t.Fatal("expected usage error")
	}
	if diagnostics.Len() == 0 {
		t.Fatal("diagnostics bypassed app output")
	}
}
