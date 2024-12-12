/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/sriddell/artifact-tracking/internal/support"
)

// recallCmd represents the recall command
var recallCmd = &cobra.Command{
	Use:   "recall",
	Short: "Downloads SBOMs associated with an artifact",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKeyName, err := cmd.Flags().GetString("api-key-name")
		apiKey, err := support.GetSSMParameter(apiKeyName)
		if err != nil {
			log.Fatal("Error getting ssm key:", err)
		}
		key, err := support.GenerateKey(args[0])
		if err != nil {
			log.Fatal("Error generating key:", err)
		}

		client := &http.Client{}
		req, err := http.NewRequest("POST", "https://artifact-metadata-service.devsecops.devops.ellucian.com/api/v1/artifacts/"+key+"/associated-sboms-queries", nil)
		if err != nil {
			log.Fatal("Error creating HTTP request:", err)
		}

		req.Header.Set("X-API-Key", apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Error making HTTP request:", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 201 {
			log.Fatalf("Error: received non-201 response code: %d", resp.StatusCode)
		}
		type sbomKey struct {
			ArtifactKey string `json:"artifact_key"`
			SbomKey     string `json:"sbom_key"`
			Link        string `json:"link"`
		}
		type queryResponse struct {
			Message string    `json:"message"`
			Keys    []sbomKey `json:"keys"`
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error reading response body: ", err)
		}
		if err != nil {
			log.Fatal("Error reading response body: ", err)
		}
		reader := io.NopCloser(bytes.NewReader(body))
		fmt.Println(string(body))
		var qs queryResponse
		if err := json.NewDecoder(io.Reader(reader)).Decode(&qs); err != nil {
			log.Fatalf("Error decoding JSON response for body: %s, %v", string(body), err)
		}

		for _, entry := range qs.Keys {
			req, err = http.NewRequest("GET", entry.Link, nil)
			if err != nil {
				log.Fatalf("Error creating HTTP GET request to %s: %v", entry.Link, err)
			}
			req.Header.Set("X-API-Key", apiKey)
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalf("Error making HTTP GET request to %s: %v", entry.Link, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				log.Fatalf("Error: received non-200 response code getting url: %v %d", entry.Link, resp.StatusCode)
			}

			var resultData map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&resultData); err != nil {
				log.Fatalf("Error decoding JSON response from %s: %v", entry.Link, err)
			}

			sbom, ok := resultData["sbom"]
			if !ok {
				log.Fatalf("Error: sbom key %s not found in response from %s", entry.SbomKey, entry.Link)
			}
			file, err := os.Create(entry.SbomKey + ".json")
			if err != nil {
				log.Fatalf("Error creating file %s: %v", entry.SbomKey, err)
			}
			defer file.Close()

			sbomData, err := json.MarshalIndent(sbom, "", "  ")
			if err != nil {
				log.Fatalf("Error marshaling sbom data: %v", err)
			}

			if _, err := file.Write(sbomData); err != nil {
				log.Fatalf("Error writing to file %s: %v", entry.SbomKey, err)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(recallCmd)
	var apiKeyName string
	recallCmd.Flags().StringVarP(&apiKeyName, "api-key-name", "a", "", "SSM parameter store name for artifact service API key")
	recallCmd.MarkFlagRequired("api-key-name")
	// recallCmd.MarkFlagRequired("file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// recallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// recallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
