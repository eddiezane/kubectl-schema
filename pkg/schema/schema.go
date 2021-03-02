package schema

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// SchemaOptions
type SchemaOptions struct {
	ConfigFlags *genericclioptions.ConfigFlags

	SchemaType string

	genericclioptions.IOStreams
}

func NewSchemaOptions(streams genericclioptions.IOStreams) *SchemaOptions {
	return &SchemaOptions{
		ConfigFlags: genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag(),
		IOStreams:   streams,
	}
}

func NewCmdSchema(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewSchemaOptions(streams)

	cmd := &cobra.Command{
		Use:   "kubectl schema <type>",
		Short: "prints out the Kubernetes schema",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}

func (o *SchemaOptions) Complete(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Type of schema required")
	}

	o.SchemaType = strings.ToLower(args[0])

	return nil
}

func (o *SchemaOptions) Validate() error {
	if o.SchemaType != "json" {
		return errors.New("schema must be 'json'")
	}

	return nil
}

func (o *SchemaOptions) Run() error {
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(o.ConfigFlags)
	factory := cmdutil.NewFactory(matchVersionKubeConfigFlags)

	client, err := factory.RESTClient()
	if err != nil {
		return err
	}

	res, err := client.Get().AbsPath("/openapi/v2").DoRaw(context.TODO())
	if err != nil {
		return err
	}

	jsonSchema := &apiextensionsv1.JSONSchemaProps{}

	err = json.Unmarshal(res, jsonSchema)
	if err != nil {
		return err
	}

	oneOf := make([]apiextensionsv1.JSONSchemaProps, 0, len(jsonSchema.Definitions))
	for k := range jsonSchema.Definitions {
		oneOf = append(oneOf, apiextensionsv1.JSONSchemaProps{Ref: strptr("#/definitions/" + k)})
	}

	jsonSchema.OneOf = oneOf

	j, err := json.Marshal(jsonSchema)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(o.Out, string(j))

	return nil
}

func strptr(s string) *string {
	return &s
}
