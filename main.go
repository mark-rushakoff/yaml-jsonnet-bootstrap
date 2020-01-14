package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) != 1 {
		fmt.Fprintf(os.Stderr, `USAGE: %[1]s

%[1]s accepts a YAML document on stdin and emits valid Jsonnet to stdout.
This simply bootstraps the initial conversion from YAML to Jsonnet.
You will want to run jsonnetfmt against the output,
and then you will want to start renaming variables and adding parameters as appropriate.
`, filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)

	n, err := declareLocalObjects(in, out)
	if err != nil {
		out.Flush()
		fmt.Fprintln(os.Stderr, "fatal: "+err.Error())
		os.Exit(1)
	}

	if err := manifest(out, n); err != nil {
		out.Flush()
		fmt.Fprintln(os.Stderr, "fatal: "+err.Error())
		os.Exit(1)
	}

	out.Flush()
	os.Exit(0)
}

// declareLocalObjects decodes each value in the YAML stream until EOF,
// and emits valid Jsonnet to stdout in form:
//    local var0 = { /* ... */ };
//    local var1 = { /* ... */ };
//    // etc.
func declareLocalObjects(in *bufio.Reader, out *bufio.Writer) (n int, err error) {
	dec := yaml.NewDecoder(in)
	for ; ; n++ {
		var obj interface{}
		err = dec.Decode(&obj)
		if err != nil {
			if err == io.EOF {
				// All decoded successfully.
				return n, nil
			}

			return n, err
		}

		// Emit the local variable assignment.
		if _, err := fmt.Fprintf(out, "local var%d = ", n); err != nil {
			return n, err
		}

		// Use MarshalIndent instead of an Encoder, so that we don't encode a trailing newline.
		// But, this has the side effect of encoding "HTML-safe" JSON strings.
		buf, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			return n, err
		}
		if _, err := io.Copy(out, bytes.NewReader(buf)); err != nil {
			return n, err
		}

		if _, err := out.WriteString(";\n\n"); err != nil {
			return n, err
		}
	}
}

// manifest emits the Jsonnet manifest, the final emitted value from evaluating a .jsonnet value.
// (I'm not sure if manifest is the right word for this.)
// To save a little boilerplate, this is in form:
//    {
//      Objects(conf):: [
//        var0,
//        var1,
//        // ...
//      ],
//    }
func manifest(out *bufio.Writer, n int) error {
	if _, err := out.WriteString("{\nObjects(conf):: ["); err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		if _, err := fmt.Fprintf(out, "\nvar%d,", i); err != nil {
			return err
		}
	}

	if _, err := out.WriteString("\n],\n}\n"); err != nil {
		return err
	}

	return nil
}
