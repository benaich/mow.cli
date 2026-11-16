package cli

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"testing"
)

func TestHelpPrinterInheritedWithoutRunningActions(t *testing.T) {
	for _, args := range [][]string{{"example", "--help"}, {"example", "group", "--help"}, {"example", "group", "write", "--help"}, {"example", "group", "write", "--invalid"}} {
		app := App("example", "example app")
		app.ErrorHandling = flag.ContinueOnError
		var out bytes.Buffer
		app.Output = &out
		path := ""
		app.HelpPrinter = func(p string, _ bool, w io.Writer) { path = p; fmt.Fprintln(w, "custom help", p) }
		app.Command("group", "group", func(cmd *Cmd) {
			cmd.Command("write", "write", func(cmd *Cmd) {
				cmd.StringArg("ID", "", "required ID")
				cmd.Action = func() { t.Fatal("help/invalid arguments ran action") }
			})
		})
		err := app.Run(args)
		if args[len(args)-1] == "--invalid" && err == nil {
			t.Fatal("parser no longer rejects invalid arguments")
		}
		if path == "" || !bytes.Contains(out.Bytes(), []byte("custom help "+path)) {
			t.Fatal(path, out.String(), err)
		}
	}
}
