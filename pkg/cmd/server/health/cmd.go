package health

import (
	"github.com/cmgsj/blob/pkg/blob"
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
		Request: &healthv1.HealthCheckRequest{
			Service: blob.ServiceName,
		},
	}
}

func NewCmdHealth(f cmdutil.Factory, streams cmdutil.IOStreams) *cobra.Command {
	o := NewWriteOptions(streams)
	cmd := &cobra.Command{
		Use:   "health",
		Short: "health-check blob server",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			stderr := o.IOStreams.Err
			cmdutil.CheckErr(o.Complete(f, cmd, args), stderr)
			cmdutil.CheckErr(o.Validate(), stderr)
			cmdutil.CheckErr(o.Run(f, cmd), stderr)
		},
	}
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
		"service": o.Request.GetService(),
		"status":  healthv1.HealthCheckResponse_ServingStatus_name[int32(resp.GetStatus())],
	}
	err = cmdutil.PrintJSON(o.IOStreams.Out, v)
	return err
}
