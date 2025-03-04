package main

import (
	"github.com/lassenordahl/hermitcrab/versions"
	"github.com/spf13/cobra"
)

var (
	prevManifestPath string
	newVersionPath   string
	outputPath       string
)

var cmdValidate = &cobra.Command{
	Use:   "validate",
	Short: "Validate a version",
	Long:  `Validate a version in the system`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := versions.ValidateVersion(newVersionPath); err != nil {
			return err
		}
		return nil
	},
}

var cmdAdd = &cobra.Command{
	Use:   "add",
	Short: "Add a new version",
	Long:  `Add a new version to the system`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := versions.AddNewVersionToManifest(prevManifestPath, outputPath, newVersionPath); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	// Validate flags.
	cmdValidate.Flags().StringVarP(&newVersionPath, "version-path", "p", "", "Path to the version to validate")

	// Add flags
	cmdAdd.Flags().StringVarP(&prevManifestPath, "prev-manifest", "p", "", "Path to the previous manifest")
	cmdAdd.Flags().StringVarP(&newVersionPath, "version-path", "v", "", "Path to the new version")
	cmdAdd.Flags().StringVarP(&outputPath, "output-path", "o", "", "Path to the output manifest")
}
