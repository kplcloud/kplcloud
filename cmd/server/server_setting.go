/**
 * @Time : 2020/7/10 3:13 PM
 * @Author : solacowa@gmail.com
 * @File : service_setting
 * @Software: GoLand
 */

package server

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-kit/kit/log/level"
	"github.com/spf13/cobra"
)

var (
	settingCmd = &cobra.Command{
		Use:               "setting command <args> [flags]",
		Short:             "置置设置命令",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `## 命令参考
kit-admin setting -h
`,
	}

	settingAddCmd = &cobra.Command{
		Use:               `add <args> [flags]`,
		Short:             "增加设置",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
kit-admin setting add HELLO world --desc "这是干啥的"
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 关闭资源连接
			defer func() {
				_ = level.Debug(logger).Log("db", "close", "err", db.Close())
				if rds != nil {
					_ = level.Debug(logger).Log("redis", "close", "err", rds.Close(context.Background()))
				}
			}()
			if len(args) < 2 {
				fmt.Println("至少需要两个参数")
				return errors.New("参数错误")
			}
			return addSetting(args[0], args[1], desc)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepare()
		},
	}

	settingDelCmd = &cobra.Command{
		Use:               `delete <args> [flags]`,
		Short:             "删除设置",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
kit-admin setting delete HELLO
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 关闭资源连接
			defer func() {
				_ = level.Debug(logger).Log("db", "close", "err", db.Close())
				if rds != nil {
					_ = level.Debug(logger).Log("redis", "close", "err", rds.Close(context.Background()))
				}
			}()
			if len(args) < 1 {
				fmt.Println("至少需要一个参数")
				return errors.New("参数错误")
			}
			return deleteSetting(args[0])
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepare()
		},
	}

	settingUpdateCmd = &cobra.Command{
		Use:               `update <args> [flags]`,
		Short:             "更新设置",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
kit-admin setting update HELLO 12324 --desc "这是干啥的"
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 关闭资源连接
			defer func() {
				_ = level.Debug(logger).Log("db", "close", "err", db.Close())
				if rds != nil {
					_ = level.Debug(logger).Log("redis", "close", "err", rds.Close(context.Background()))
				}
			}()
			if len(args) < 2 {
				fmt.Println("至少需要两个参数")
				return errors.New("参数错误")
			}
			return updateSetting(args[0], args[1], desc)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepare()
		},
	}

	settingGetCmd = &cobra.Command{
		Use:               `get <args> [flags]`,
		Short:             "删除设置",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
kit-admin setting get HELLO
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 关闭资源连接
			defer func() {
				_ = level.Debug(logger).Log("db", "close", "err", db.Close())
				if rds != nil {
					_ = level.Debug(logger).Log("redis", "close", "err", rds.Close(context.Background()))
				}
			}()
			if len(args) < 1 {
				fmt.Println("至少需要一个参数")
				return errors.New("参数错误")
			}
			return getSetting(args[0])
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepare()
		},
	}
)

func getSetting(key string) error {
	res, err := store.SysSetting().Find(context.Background(), strings.ToUpper(key))
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Key:\t%s", res.Key))
	fmt.Println(fmt.Sprintf("Value:\t%s", res.Value))
	fmt.Println(fmt.Sprintf("Description:\t%s", res.Description))
	fmt.Println(fmt.Sprintf("CreatedAt:\t%s", res.CreatedAt))
	fmt.Println(fmt.Sprintf("UpdatedAt:\t%s", res.UpdatedAt))

	return nil
}

func addSetting(key, value, desc string) error {
	return store.SysSetting().Add(context.Background(), strings.ToUpper(key), value, desc)
}

func deleteSetting(key string) error {
	return store.SysSetting().Delete(context.Background(), strings.ToUpper(key))
}

func updateSetting(key, value, desc string) error {
	res, err := store.SysSetting().Find(context.Background(), strings.ToUpper(key))
	if err != nil {
		return err
	}

	res.Value = value
	res.Description = desc

	return store.SysSetting().Update(context.Background(), &res)
}
