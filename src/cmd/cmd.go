/*
 * @Time : 2019-06-25 17:52
 * @Author : solacowa@gmail.com
 * @File : cmd
 * @Software: GoLand
 */

package cmd

import (
	"flag"
	"github.com/spf13/cobra"
)

func AddFlags(rootCmd *cobra.Command) {
	flag.CommandLine.VisitAll(func(gf *flag.Flag) {
		rootCmd.PersistentFlags().AddGoFlag(gf)
	})
}
