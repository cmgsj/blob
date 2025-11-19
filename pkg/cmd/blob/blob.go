package blob

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/riza-io/grpc-go/credentials/basic"
	"github.com/riza-io/grpc-go/credentials/bearer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var version = "dev"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blob",
		Short: "Blob CLI",
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			viper.AutomaticEnv()
			viper.AllowEmptyEnv(true)
			viper.SetEnvPrefix("blob")
			viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

			return viper.BindPFlags(cmd.PersistentFlags())
		},
	}

	cmd.PersistentFlags().String("address", "localhost:2562", "server address")
	cmd.PersistentFlags().String("log-level", "info", "log level")
	cmd.PersistentFlags().String("log-handler", "text", "log handler")
	cmd.PersistentFlags().String("auth-token", "", "auth token")
	cmd.PersistentFlags().String("auth-username", "", "auth username")
	cmd.PersistentFlags().String("auth-password", "", "auth password")
	cmd.PersistentFlags().String("tls-cert", "", "tls cert file")
	cmd.PersistentFlags().String("tls-key", "", "tls key file")
	cmd.PersistentFlags().String("tls-ca", "", "tls ca file")
	cmd.PersistentFlags().String("tls-server-name", "", "tls server name")

	cmd.AddCommand(
		NewCommandList(),
		NewCommandGet(),
		NewCommandPut(),
		NewCommandDelete(),
		NewCommandHealth(),
		NewCommandServer(),
	)

	return cmd
}

func newLogger() (logging.Logger, error) {
	logLevel := viper.GetString("log-level")
	logHandler := viper.GetString("log-handler")

	var level slog.Level

	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		return nil, err
	}

	var handler slog.Handler

	switch logHandler {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})

	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})

	case "discard":
		handler = slog.DiscardHandler

	default:
		return nil, fmt.Errorf("unknown log handler %q", logHandler)
	}

	logger := slog.New(handler)

	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		logger.Log(ctx, slog.Level(level), msg, fields...)
	}), nil
}

func newServerAuthFunc() auth.AuthFunc {
	authToken := viper.GetString("auth-token")
	authUsername := viper.GetString("auth-username")
	authPassword := viper.GetString("auth-password")

	if authToken != "" {
		return newBearerAuthFunc(authToken)
	}

	if authUsername != "" && authPassword != "" {
		return newBasicAuthFunc(authUsername, authPassword)
	}

	return newInsecureAuthFunc()
}

func newClientPerRPCCredentials() credentials.PerRPCCredentials {
	authToken := viper.GetString("auth-token")
	authUsername := viper.GetString("auth-username")
	authPassword := viper.GetString("auth-password")

	if authToken != "" {
		return bearer.NewPerRPCCredentials(authToken)
	}

	if authUsername != "" && authPassword != "" {
		return basic.NewPerRPCCredentials(authUsername, authPassword)
	}

	return nil
}

func newClientTransportCredentials() (credentials.TransportCredentials, error) {
	return newTransportCredentials(false)
}

func newServerTransportCredentials() (credentials.TransportCredentials, error) {
	return newTransportCredentials(true)
}

func newTransportCredentials(isServer bool) (credentials.TransportCredentials, error) {
	tlsCert := viper.GetString("tls-cert")
	tlsKey := viper.GetString("tls-key")
	tlsCA := viper.GetString("tls-ca")
	tlsServerName := viper.GetString("tls-server-name")

	tlsConfig := &tls.Config{}

	hasCertificates := false
	hasCertPool := false

	if tlsCert != "" && tlsKey != "" {
		cert, err := os.ReadFile(tlsCert)
		if err != nil {
			return nil, fmt.Errorf("failed to read tls cert file: %w", err)
		}

		key, err := os.ReadFile(tlsKey)
		if err != nil {
			return nil, fmt.Errorf("failed to read tls key file: %w", err)
		}

		certificate, err := tls.X509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}

		tlsConfig.Certificates = []tls.Certificate{certificate}

		hasCertificates = true
	}

	if tlsCA != "" {
		ca, err := os.ReadFile(tlsCA)
		if err != nil {
			return nil, fmt.Errorf("failed to read tls ca file: %w", err)
		}

		certPool := x509.NewCertPool()

		if !certPool.AppendCertsFromPEM(ca) {
			return nil, errors.New("failed to append tls ca")
		}

		if isServer {
			tlsConfig.ClientCAs = certPool
		} else {
			tlsConfig.RootCAs = certPool
		}

		hasCertPool = true
	}

	if isServer {
		tlsConfig.ServerName = tlsServerName

		if hasCertificates && hasCertPool {
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
	}

	if hasCertificates || hasCertPool {
		return credentials.NewTLS(tlsConfig), nil
	}

	return insecure.NewCredentials(), nil
}
