/**
 * @Time : 7/21/21 11:32 AM
 * @Author : solacowa@gmail.com
 * @File : server_reset
 * @Software: GoLand
 */

package server

import (
	"github.com/spf13/cobra"
)

var (
	resetCmd = &cobra.Command{
		Use:   "reset",
		Short: "重置kplcloud",
		Example: `## 重置kplcloud
kplcloud reset
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)
