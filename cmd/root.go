package cmd

import (
	"alexejk.io/mod-up/internal/version"
	"alexejk.io/mod-up/upgrade"
	"github.com/spf13/cobra"
)

type opts struct {}

func RootCmd(v *version.Info) *cobra.Command {

	o := &opts{}

	cmd := &cobra.Command{
		Use: "mod-up [dir]",
		Short: "Upgrade all modules in directory",
		Args: cobra.ExactArgs(1),
		RunE: o.RunE,
	}

	cmd.Version = v.String()
	cmd.SilenceErrors = true

	return cmd
}

func (o *opts) RunE(cmd *cobra.Command, args []string) error {

	workDir := args[0]

	u, err := upgrade.New(workDir)
	if err != nil {
		return err
	}

	if err := u.Upgrade(false); err != nil {
		return err
	}

	if err := u.Tidy(); err != nil {
		return err
	}

	return nil
}
