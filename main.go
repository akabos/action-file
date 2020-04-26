package main

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func main() {
	var args struct {
		Content       string `envconfig:"INPUT_CONTENT"`
		Path          string `envconfig:"INPUT_PATH"`
		Encoding      string `envconfig:"INPUT_ENCODING"`
		WorkspacePath string `envconfig:"GITHUB_WORKSPACE" default:"./"`
	}
	var err error
	envconfig.MustProcess("", &args)

	debug("arguments: encoding=%q path=%q workspace-path=%q", args.Encoding, args.Path, args.WorkspacePath)

	buf := bytes.NewBuffer(nil)

	switch strings.ToLower(args.Encoding) {
	case "":
		_, err = buf.WriteString(args.Content)
	case "base64":
		_, err = buf.ReadFrom(base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(args.Content)))
	case "base32":
		_, err = buf.ReadFrom(base32.NewDecoder(base32.StdEncoding, bytes.NewBufferString(args.Content)))
	case "base58":
		_, err = buf.Write(base58.Decode(args.Content))
	default:
		err = errors.Errorf("unknown encoding: %q", args.Encoding)
	}
	if err != nil {
		die(errors.Wrap(err, "failed to decode content"))
	}

	var f *os.File

	if args.Path == "" {
		f, err = ioutil.TempFile(args.WorkspacePath, ".file-*")
	} else {
		f, err = os.Create(path.Clean(path.Join(args.WorkspacePath, args.Path)))
	}
	if err != nil {
		die(errors.Wrap(err, "failed to open output file"))
	}
	_ = f.Chown(1000, 1000)
	_ = f.Chmod(0644)
	defer f.Close()

	output("path", strings.TrimPrefix(strings.TrimPrefix(f.Name(), args.WorkspacePath), "/"))

	debug("container path: %v", f.Name())

	_, err = buf.WriteTo(f)
	if err != nil {
		die(errors.Wrap(err, "failed to write content to the file"))
	}

	os.Exit(0)
}

func die(e error) {
	_, _ = fmt.Fprintf(os.Stdout, "::error ::%s\n", e.Error())
	os.Exit(-1)
}

func output(name, value string) {
	_, _ = fmt.Fprintf(os.Stdout, "::set-output name=%s::%s\n", name, value)
}

func debug(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, "::debug ::%v\n", fmt.Sprintf(msg, args...))
}
