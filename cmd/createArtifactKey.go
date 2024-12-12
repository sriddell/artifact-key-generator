/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/sriddell/artifact-tracking/internal/support"

	"github.com/spf13/cobra"
)

// createArtifactKeyCmd represents the createArtifactKey command
var createArtifactKeyCmd = &cobra.Command{
	Use:   "createArtifactKey",
	Short: "Generates a unique artifact key",
	Long:  `Generates a key consisting of {sha512}-{filesize} for the specified artifact.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		key, err := support.GenerateKey(filename)
		if err != nil {
			log.Fatal("Error generating key:", err)
			return
		}
		fmt.Println("Generated Key:", key)
	},
}

func init() {
	rootCmd.AddCommand(createArtifactKeyCmd)
}
