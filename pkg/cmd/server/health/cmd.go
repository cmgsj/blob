package health

import (
	cmdutil "github.com/cmgsj/blob/pkg/cmd/util"
	"github.com/spf13/cobra"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

type HealthOptions struct {
	IOStreams cmdutil.IOStreams
	Request   *healthv1.HealthCheckRequest
}

func NewWriteOptions(streams cmdutil.IOStreams) *HealthOptions {
	return &HealthOptions{
		IOStreams: streams,
		Request:   &healthv1.HealthCheckRequest{},
	}
}

func NewCmdHealth(f cmdutil.Factory, streams cmdutil.IOStreams) *cobra.Command {
	o := NewWriteOptions(streams)
	cmd := &cobra.Command{
		Use:  "health",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			stderr := o.IOStreams.Err
			cmdutil.CheckErr(o.Complete(f, cmd, args), stderr)
			cmdutil.CheckErr(o.Validate(), stderr)
			cmdutil.CheckErr(o.Run(f, cmd), stderr)
		},
	}
	cmd.Flags().StringVar(&o.Request.Service, "service", o.Request.Service, "service to health check")
	return cmd
}

func (o *HealthOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	return nil
}

func (o *HealthOptions) Validate() error {
	return nil
}

func (o *HealthOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	resp, err := f.HealthClient().Check(cmd.Context(), o.Request)
	if err != nil {
		return err
	}
	v := map[string]string{
		"status": healthv1.HealthCheckResponse_ServingStatus_name[int32(resp.GetStatus())],
	}
	err = cmdutil.PrintJSON(o.IOStreams.Out, v)
	return err
}
