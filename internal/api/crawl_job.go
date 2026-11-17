package api

import (
	"context"
	"fmt"
	"regexp"

	"google.golang.org/api/run/v2"
)

// CrawlIDEnvVar is how a crawl job execution learns which crawl to run.
const CrawlIDEnvVar = "CRAWL_ID"

// A Cloud Run Job's full resource name
var crawlJobNamePattern = regexp.MustCompile(`^projects/[^/]+/locations/[^/]+/jobs/[^/]+$`)

// cloudRunJobRunner starts crawls as executions of one Cloud Run Job.
type cloudRunJobRunner struct {
	jobs    *run.ProjectsLocationsJobsService
	jobName string // projects/{project}/locations/{region}/jobs/{job}
}

// newCloudRunJobRunner builds a runner for the given job, authenticating with
// the environment's default Google credentials (on Cloud Run, the service's
// own service account).
func newCloudRunJobRunner(ctx context.Context, jobName string) (*cloudRunJobRunner, error) {
	if !crawlJobNamePattern.MatchString(jobName) {
		return nil, fmt.Errorf("crawl job %q must look like projects/PROJECT/locations/REGION/jobs/NAME", jobName)
	}
	service, err := run.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("create Cloud Run client: %w", err)
	}
	return &cloudRunJobRunner{jobs: service.Projects.Locations.Jobs, jobName: jobName}, nil
}

// Run starts one execution of the job for the given crawl. It returns once
// Cloud Run has accepted the execution; it does not wait for the crawl.
func (r *cloudRunJobRunner) Run(ctx context.Context, crawlID string) error {
	request := &run.GoogleCloudRunV2RunJobRequest{
		Overrides: &run.GoogleCloudRunV2Overrides{
			ContainerOverrides: []*run.GoogleCloudRunV2ContainerOverride{{
				Env: []*run.GoogleCloudRunV2EnvVar{{Name: CrawlIDEnvVar, Value: crawlID}},
			}},
		},
	}
	if _, err := r.jobs.Run(r.jobName, request).Context(ctx).Do(); err != nil {
		return fmt.Errorf("run job %s: %w", r.jobName, err)
	}
	return nil
}
