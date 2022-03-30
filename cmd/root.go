/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"image/png"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "appicon",
	Short: "iOS AppIcon 生成器",
	Long:  `通过一张尺寸为 1024*1024 的图片制作iOS应用的 AppIcon`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
			return
		}
		fmt.Println(path)
		file, err := os.Open(path)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		img, err := png.Decode(file)
		file.Close()
		fmt.Println(img.Bounds())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.appicon.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("version", "v", false, "Version information")
	rootCmd.Flags().StringP("path", "p", "", "path of image")
	rootCmd.Flags().StringP("output", "o", "", "output path of AppIcon.appiconset  (default current directory)")
}
