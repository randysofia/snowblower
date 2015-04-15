package main

import "github.com/spf13/cobra"

func main() {

	var collectorCmd = &cobra.Command{
		Use:   "collect",
		Short: "Run the collector",
		Run: func(cmd *cobra.Command, args []string) {
			startCollector()
		},
	}

	var rootCmd = &cobra.Command{Use: "snowblower"}
	rootCmd.AddCommand(collectorCmd)
	rootCmd.Execute()

}
