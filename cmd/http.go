package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var HttpCmd = &cobra.Command{
	Use:   "http serve",
	Short: "Run Http API",
	Long:  "Run Http API",
	RunE: func(cmd *cobra.Command, args []string) error {
		echan := make(chan error)
		go func() {
			echan <- initHTTP()
		}()

		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)

		select {
		case <-term:
			logrus.Infoln("signal terminated detected")
			return nil
		case err := <-echan:
			return errors.Wrap(err, "service runtime error")
		}
	},
}
