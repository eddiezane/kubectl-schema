package main

import (
	"os"

	"github.com/spf13/pflag"

	"github.com/eddiezane/kubectl-schema/pkg/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-schema", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := schema.NewCmdSchema(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
