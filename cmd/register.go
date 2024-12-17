/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"compress/gzip"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/sriddell/artifact-tracker/internal/support"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Registers the specified artifact's metadata",
	Long: `Computes a unique hash key for the specified artifact and 
	       registers the metadata with the Artifact Metadata Service.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKeyName, err := cmd.Flags().GetString("api-key-name")
		if err != nil {
			log.Fatal("Error getting api-key-name flag:", err)
		}
		sbomFile, err := cmd.Flags().GetString("sbom")
		if err != nil {
			log.Fatal("Error getting sbom flag:", err)
		}
		apiKey, err := support.GetSSMParameter(apiKeyName)
		if err != nil {
			log.Fatal("Error getting ssm key:", err)
		}
		key, err := support.GenerateKey(args[0])
		if err != nil {
			log.Fatal("Error generating key:", err)
		}
		content, err := os.ReadFile(sbomFile)
		if err != nil {
			log.Fatal("Error reading sbom file:", err)
		}
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		if _, err := gz.Write(content); err != nil {
			log.Fatal("Error compressing sbom content:", err)
		}
		if err := gz.Close(); err != nil {
			log.Fatal("Error closing gzip writer:", err)
		}
		req, err := http.NewRequest("POST", "https://artifact-metadata-service.devsecops.devops.ellucian.com/api/v1/artifacts/"+key+"/associated-sboms", &buf)
		if err != nil {
			log.Fatal("Error creating request for sbom content:", err)
		}
		req.Header.Set("X-API-Key", apiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")
		if err != nil {
			log.Fatal("Error posting sbom content:", err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Error making HTTP request:", err)
		}
		if resp.StatusCode != 201 {
			log.Fatalf("Failed to post sbom content, status code: %d", resp.StatusCode)
		}

	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	var apiKeyName string
	var sbomFile string
	registerCmd.Flags().StringVarP(&apiKeyName, "api-key-name", "a", "", "SSM parameter store name for artifact service API key")
	registerCmd.MarkFlagRequired("api-key-name")

	registerCmd.Flags().StringVarP(&sbomFile, "sbom", "s", "", "An sbom json file to register against the artifact")
	registerCmd.MarkFlagRequired("sbom")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
