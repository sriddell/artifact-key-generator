/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/sriddell/artifact-tracker/internal/support"
)

// getKeyCmd represents the getKey command
var getKeyCmd = &cobra.Command{
	Use:   "getKey",
	Short: "Return the computed key for the specified artifact",
	Long:  `Generates a key consisting of {sha512}-{filesize} for the specified artifact.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		key, err := support.GenerateKey(filename)
		if err != nil {
			log.Fatal("Error generating key:", err)
			return
		}
		fmt.Println(key)
	},
}

func init() {
	rootCmd.AddCommand(getKeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getKeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getKeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
