package blob

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"slices"
	"strings"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/validate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	blobhandler "github.com/cmgsj/blob/pkg/blob/handler"
	blobstorage "github.com/cmgsj/blob/pkg/blob/storage"
	blobstorageriver "github.com/cmgsj/blob/pkg/blob/storage/driver"
	"github.com/cmgsj/blob/pkg/blob/storage/driver/azblob"
	"github.com/cmgsj/blob/pkg/blob/storage/driver/gcs"
	"github.com/cmgsj/blob/pkg/blob/storage/driver/memory"
	"github.com/cmgsj/blob/pkg/blob/storage/driver/minio"
	"github.com/cmgsj/blob/pkg/blob/storage/driver/s3"
	"github.com/cmgsj/blob/pkg/docs"
	"github.com/cmgsj/blob/pkg/proto/blob/api/v1/apiv1connect"
)

func NewCommandServer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start blob server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)

			address := viper.GetString("address")
			debugExpvar := viper.GetBool("debug-expvar")
			debugPprof := viper.GetBool("debug-pprof")

			blobStorageDriver, err := newBlobStorageDriver(ctx)
			if err != nil {
				return err
			}

			blobStorage, err := blobstorage.NewStorage(ctx, blobStorageDriver)
			if err != nil {
				return err
			}

			blobHandler := blobhandler.NewHandler(blobStorage)

			services := []string{
				grpchealth.HealthV1ServiceName,
				grpcreflect.ReflectV1ServiceName,
				grpcreflect.ReflectV1AlphaServiceName,
				apiv1connect.BlobServiceName,
			}

			grpcHealthChecker := grpchealth.NewStaticChecker(services...)
			grpcReflectReflector := grpcreflect.NewStaticReflector(services...)

			handlerOpts := []connect.HandlerOption{
				connect.WithInterceptors(
					validate.NewInterceptor(
						validate.WithValidateResponses(),
					),
				),
			}

			mux := http.NewServeMux()

			mux.Handle(grpchealth.NewHandler(grpcHealthChecker, handlerOpts...))
			mux.Handle(grpcreflect.NewHandlerV1(grpcReflectReflector, handlerOpts...))
			mux.Handle(grpcreflect.NewHandlerV1Alpha(grpcReflectReflector, handlerOpts...))
			mux.Handle(apiv1connect.NewBlobServiceHandler(blobHandler, handlerOpts...))

			mux.Handle("GET /docs/", http.StripPrefix("/docs", http.FileServerFS(docs.Assets())))

			if debugExpvar {
				mux.Handle("GET /debug/vars", expvar.Handler())
			}

			if debugPprof {
				mux.Handle("GET /debug/pprof/", http.HandlerFunc(pprof.Index))
				mux.Handle("GET /debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
				mux.Handle("GET /debug/pprof/profile", http.HandlerFunc(pprof.Profile))
				mux.Handle("GET /debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
				mux.Handle("GET /debug/pprof/trace", http.HandlerFunc(pprof.Trace))
			}

			httpServer := &http.Server{
				Addr:    address,
				Handler: h2c.NewHandler(mux, &http2.Server{}),
			}

			go func() {
				slog.Info("starting http server", "address", address)

				err := httpServer.ListenAndServe()
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					slog.Error("failed to run http server", "address", address, "error", err)
				}

				cancel()
			}()

			<-ctx.Done()

			slog.Info("shutting down http server", "address", address)

			err = httpServer.Shutdown(ctx)
			if err != nil {
				slog.Error("failed to shut down http server", "address", address, "error", err)

				return err
			}

			return nil
		},
	}

	cmd.Flags().String("address", ":8080", "blob service address")

	cmd.Flags().String("azblob-uri", "", "azblob uri")
	cmd.Flags().String("azblob-account-name", "", "azblob account name")
	cmd.Flags().String("azblob-account-key", "", "azblob account key")

	cmd.Flags().String("gcs-uri", "", "gcs uri")

	cmd.Flags().String("s3-uri", "", "s3 uri")
	cmd.Flags().String("s3-access-key", "", "s3 access key")
	cmd.Flags().String("s3-secret-key", "", "s3 secret key")

	cmd.Flags().String("minio-address", "", "minio address")
	cmd.Flags().String("minio-access-key", "", "minio access key")
	cmd.Flags().String("minio-secret-key", "", "minio secret key")
	cmd.Flags().String("minio-bucket", "", "minio bucket")
	cmd.Flags().String("minio-object-prefix", "", "minio object prefix")
	cmd.Flags().Bool("minio-secure", false, "minio secure")

	cmd.Flags().Bool("debug-expvar", false, "debug expvar")
	cmd.Flags().Bool("debug-pprof", false, "debug pprof")

	_ = viper.BindPFlags(cmd.Flags())

	return cmd
}

func newBlobStorageDriver(ctx context.Context) (blobstorageriver.Driver, error) {
	driverTypes := []string{
		"azblob",
		"gcs",
		"s3",
		"minio",
	}

	driverTypeSet := make(map[string]struct{})

	for _, key := range viper.AllKeys() {
		for _, driverType := range driverTypes {
			if strings.HasPrefix(key, driverType) && viper.IsSet(key) {
				driverTypeSet[driverType] = struct{}{}

				break
			}
		}
	}

	driverTypes = slices.Collect(maps.Keys(driverTypeSet))

	if len(driverTypes) > 1 {
		return nil, fmt.Errorf("multiple blob storage drivers set: [%s]", strings.Join(driverTypes, ", "))
	}

	if len(driverTypes) == 0 {
		driverTypes = []string{"memory"}
	}

	driverType := driverTypes[0]

	var (
		driver blobstorageriver.Driver
		err    error
	)

	switch driverType {
	case "memory":
		driver = memory.NewDriver()

	case "azblob":
		driver, err = azblob.NewDriver(ctx, azblob.DriverOptions{
			URI:         viper.GetString("azblob-uri"),
			AccountName: viper.GetString("azblob-account-name"),
			AccountKey:  viper.GetString("azblob-account-key"),
		})

	case "gcs":
		driver, err = gcs.NewDriver(ctx, gcs.DriverOptions{
			URI: viper.GetString("gcs-uri"),
		})

	case "s3":
		driver, err = s3.NewDriver(ctx, s3.DriverOptions{
			URI:       viper.GetString("s3-uri"),
			AccessKey: viper.GetString("s3-access-key"),
			SecretKey: viper.GetString("s3-secret-key"),
		})

	case "minio":
		driver, err = minio.NewDriver(ctx, minio.DriverOptions{
			Address:      viper.GetString("minio-address"),
			AccessKey:    viper.GetString("minio-access-key"),
			SecretKey:    viper.GetString("minio-secret-key"),
			Bucket:       viper.GetString("minio-bucket"),
			ObjectPrefix: viper.GetString("minio-object-prefix"),
			Secure:       viper.GetBool("minio-secure"),
		})

	default:
		return nil, fmt.Errorf("unknown blob storage driver %q", driverType)
	}

	if err != nil {
		return nil, err
	}

	slog.Info("using blob storage driver", "type", driverType)

	return driver, nil
}
