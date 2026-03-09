package acli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/spf13/cobra"
)

var bbPipelineCmd = &cobra.Command{
	Use:     "pipeline",
	Short:   "Manage pipelines",
	Aliases: []string{"pipe"},
	RunE:    helpRunE,
}

func init() {
	// pipeline list
	pipelineListCmd := &cobra.Command{
		Use:     "list <workspace> <repo-slug>",
		Short:   "List pipelines",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			status, _ := cmd.Flags().GetString("status")
			pipelines, err := client.ListPipelines(args[0], args[1], &bitbucket.ListPipelinesOptions{
				Status: status,
			})
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "BUILD#\tSTATUS\tTRIGGER\tTARGET\tCREATED")
			for _, p := range pipelines {
				status := p.State.Name
				if p.State.Result != nil {
					status = p.State.Result.Name
				}
				target := ""
				if p.Target.RefName != "" {
					target = p.Target.RefName
				}
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
					p.BuildNumber, status, p.Trigger.Name, target, p.CreatedOn)
			}
			return w.Flush()
		},
	}
	pipelineListCmd.Flags().String("status", "", "Filter by status (PENDING, BUILDING, PASSED, FAILED, etc.)")
	bbPipelineCmd.AddCommand(pipelineListCmd)

	// pipeline get
	bbPipelineCmd.AddCommand(&cobra.Command{
		Use:   "get <workspace> <repo-slug> <pipeline-uuid>",
		Short: "Get pipeline details",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			pipeline, err := client.GetPipeline(args[0], args[1], args[2])
			if err != nil {
				return err
			}

			status := pipeline.State.Name
			if pipeline.State.Result != nil {
				status = pipeline.State.Result.Name
			}

			fmt.Printf("Build #:      %d\n", pipeline.BuildNumber)
			fmt.Printf("UUID:         %s\n", pipeline.UUID)
			fmt.Printf("Status:       %s\n", status)
			fmt.Printf("Trigger:      %s\n", pipeline.Trigger.Name)
			fmt.Printf("Target:       %s\n", pipeline.Target.RefName)
			fmt.Printf("Commit:       %s\n", pipeline.Target.Commit.Hash)
			fmt.Printf("Creator:      %s\n", pipeline.Creator.DisplayName)
			fmt.Printf("Created:      %s\n", pipeline.CreatedOn)
			if pipeline.CompletedOn != "" {
				fmt.Printf("Completed:    %s\n", pipeline.CompletedOn)
			}
			if pipeline.BuildSecondsUsed > 0 {
				fmt.Printf("Build Time:   %ds\n", pipeline.BuildSecondsUsed)
			}
			return nil
		},
	})

	// pipeline run
	pipelineRunCmd := &cobra.Command{
		Use:   "run <workspace> <repo-slug>",
		Short: "Run a pipeline",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			branch, _ := cmd.Flags().GetString("branch")
			custom, _ := cmd.Flags().GetString("custom")

			if branch == "" {
				return fmt.Errorf("--branch is required")
			}

			var req *bitbucket.RunPipelineRequest
			if custom != "" {
				req = bitbucket.NewCustomPipelineRequest(branch, custom)
			} else {
				req = bitbucket.NewBranchPipelineRequest(branch)
			}

			pipeline, err := client.RunPipeline(args[0], args[1], req)
			if err != nil {
				return err
			}

			fmt.Printf("Started pipeline build #%d\n", pipeline.BuildNumber)
			fmt.Printf("UUID: %s\n", pipeline.UUID)
			return nil
		},
	}
	pipelineRunCmd.Flags().String("branch", "", "Branch to run pipeline on (required)")
	pipelineRunCmd.Flags().String("custom", "", "Custom pipeline pattern to run")
	bbPipelineCmd.AddCommand(pipelineRunCmd)

	// pipeline stop
	bbPipelineCmd.AddCommand(&cobra.Command{
		Use:   "stop <workspace> <repo-slug> <pipeline-uuid>",
		Short: "Stop a running pipeline",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			if err := client.StopPipeline(args[0], args[1], args[2]); err != nil {
				return err
			}
			fmt.Printf("Stopped pipeline %s\n", args[2])
			return nil
		},
	})

	// pipeline steps
	bbPipelineCmd.AddCommand(&cobra.Command{
		Use:   "steps <workspace> <repo-slug> <pipeline-uuid>",
		Short: "List steps for a pipeline",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			steps, err := client.ListPipelineSteps(args[0], args[1], args[2])
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "UUID\tNAME\tSTATUS\tDURATION\tIMAGE")
			for _, s := range steps {
				status := s.State.Name
				if s.State.Result != nil {
					status = s.State.Result.Name
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%ds\t%s\n",
					s.UUID, s.Name, status, s.DurationInSeconds, s.Image.Name)
			}
			return w.Flush()
		},
	})

	// pipeline log
	bbPipelineCmd.AddCommand(&cobra.Command{
		Use:   "log <workspace> <repo-slug> <pipeline-uuid> <step-uuid>",
		Short: "Get log output for a pipeline step",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			log, err := client.GetStepLog(args[0], args[1], args[2], args[3])
			if err != nil {
				return err
			}
			fmt.Print(log)
			return nil
		},
	})

	// pipeline variables
	bbPipelineCmd.AddCommand(&cobra.Command{
		Use:   "variables <workspace> <repo-slug>",
		Short: "List pipeline variables",
		Aliases: []string{"vars"},
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			vars, err := client.ListPipelineVariables(args[0], args[1])
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "UUID\tKEY\tVALUE\tSECURED")
			for _, v := range vars {
				value := v.Value
				if v.Secured {
					value = "***"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", v.UUID, v.Key, value, v.Secured)
			}
			return w.Flush()
		},
	})

	// pipeline add-variable
	addVarCmd := &cobra.Command{
		Use:   "add-variable <workspace> <repo-slug>",
		Short: "Create a pipeline variable",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}

			key, _ := cmd.Flags().GetString("key")
			value, _ := cmd.Flags().GetString("value")
			secured, _ := cmd.Flags().GetBool("secured")

			if key == "" || value == "" {
				return fmt.Errorf("--key and --value are required")
			}

			v, err := client.CreatePipelineVariable(args[0], args[1], key, value, secured)
			if err != nil {
				return err
			}
			fmt.Printf("Created variable: %s (UUID: %s)\n", v.Key, v.UUID)
			return nil
		},
	}
	addVarCmd.Flags().String("key", "", "Variable key (required)")
	addVarCmd.Flags().String("value", "", "Variable value (required)")
	addVarCmd.Flags().Bool("secured", false, "Mark variable as secured")
	bbPipelineCmd.AddCommand(addVarCmd)

	// pipeline delete-variable
	bbPipelineCmd.AddCommand(&cobra.Command{
		Use:   "delete-variable <workspace> <repo-slug> <variable-uuid>",
		Short: "Delete a pipeline variable",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := bitbucket.NewClient()
			if err != nil {
				return err
			}
			if err := client.DeletePipelineVariable(args[0], args[1], args[2]); err != nil {
				return err
			}
			fmt.Printf("Deleted variable %s\n", args[2])
			return nil
		},
	})
}
