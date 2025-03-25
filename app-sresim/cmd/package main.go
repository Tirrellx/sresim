package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"

	"github.com/tmc/langchaingo/llms/openai"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

func main() {
	ctx := context.Background()

	// GitHub Token from environment variable
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("GITHUB_TOKEN environment variable not set")
	}

	// OpenAI API Key from environment variable
	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable not set")
	}

	// Prometheus Server URL from environment variable
	prometheusURL := os.Getenv("PROMETHEUS_URL")
	if prometheusURL == "" {
		log.Fatal("PROMETHEUS_URL environment variable not set")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	client := github.NewClient(oauth2.NewClient(ctx, ts))

	owner := "your-org" // Replace with your GitHub org/user
	repo := "your-repo" // Replace with your GitHub repository

	listOpts := &github.ListWorkflowRunsOptions{Status: "completed", ListOptions: github.ListOptions{PerPage: 5}}
	runs, _, err := client.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, listOpts)
	if err != nil {
		log.Fatalf("Error fetching workflow runs: %v", err)
	}

	fmt.Println("Analyzing recent workflow runs with AI suggestions and Prometheus metrics...")
	for _, run := range runs.WorkflowRuns {
		metrics := fetchAdvancedPrometheusMetrics(prometheusURL, ctx)
		suggestion := analyzeWorkflowRunWithAI(ctx, run, openAIKey, metrics)
		fmt.Printf("\nWorkflow: %s | Run ID: %d\n", run.GetName(), run.GetID())
		fmt.Printf("Status: %s | Conclusion: %s\n", run.GetStatus(), run.GetConclusion())
		fmt.Printf("Total Duration: %s\n", run.GetUpdatedAt().Sub(run.GetRunStartedAt()).Round(time.Second))
		fmt.Printf("AI Suggestion: %s\n", suggestion)
	}
}

func fetchAdvancedPrometheusMetrics(prometheusURL string, ctx context.Context) string {
	client, err := api.NewClient(api.Config{Address: prometheusURL})
	if err != nil {
		log.Fatalf("Error creating Prometheus client: %v", err)
	}

	v1api := v1.NewAPI(client)
	metrics := ""

	queries := map[string]string{
		"Build Duration Rate":         "rate(ci_cd_build_duration_seconds_sum[5m])",
		"CPU Usage":                   "avg(rate(container_cpu_usage_seconds_total{image!=""}[5m])) by (pod)",
		"Memory Usage":                "avg(container_memory_usage_bytes{image!=""}) by (pod)",
		"Error Rate":                  "rate(ci_cd_pipeline_errors_total[5m])",
		"Failed Jobs":                 "sum(increase(ci_cd_job_failures_total[5m]))",
		"Success Rate":                "rate(ci_cd_pipeline_success_total[5m])",
	}

	for description, query := range queries {
		result, warnings, err := v1api.Query(ctx, query, time.Now())
		if err != nil {
			log.Printf("Error querying Prometheus for %s: %v", description, err)
			continue
		}
		if len(warnings) > 0 {
			log.Printf("Prometheus warnings for %s: %v", description, warnings)
		}

		if vector, ok := result.(model.Vector); ok {
			for _, sample := range vector {
				metrics += fmt.Sprintf("%s - Metric: %s, Value: %f\n", description, sample.Metric, sample.Value)
			}
		}
	}

	return metrics
}

func analyzeWorkflowRunWithAI(ctx context.Context, run *github.WorkflowRun, apiKey string, metrics string) string {
	llm, err := openai.New(openai.WithModel("gpt-4"), openai.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to initialize OpenAI client: %v", err)
	}

	duration := run.GetUpdatedAt().Sub(run.GetRunStartedAt()).Round(time.Second)
	prompt := fmt.Sprintf(`Analyze the following CI/CD workflow run and provide suggestions for optimization:
	Workflow Name: %s
	Status: %s
	Conclusion: %s
	Total Duration: %s
	Prometheus Metrics: %s
	Suggest improvements to speed up the build process or optimize steps if needed.`, run.GetName(), run.GetStatus(), run.GetConclusion(), duration, metrics)

	response, err := llm.Call(ctx, prompt)
	if err != nil {
		log.Fatalf("Error from LLM: %v", err)
	}

	return response
}
