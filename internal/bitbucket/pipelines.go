package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type Pipeline struct {
	UUID            string `json:"uuid"`
	BuildNumber     int    `json:"build_number"`
	CreatedOn       string `json:"created_on"`
	CompletedOn     string `json:"completed_on"`
	BuildSecondsUsed int   `json:"build_seconds_used"`
	Creator         struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
	} `json:"creator"`
	Target struct {
		Type     string `json:"type"`
		RefType  string `json:"ref_type"`
		RefName  string `json:"ref_name"`
		Selector struct {
			Type    string `json:"type"`
			Pattern string `json:"pattern"`
		} `json:"selector"`
		Commit struct {
			Hash string `json:"hash"`
		} `json:"commit"`
	} `json:"target"`
	Trigger struct {
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"trigger"`
	State struct {
		Name   string `json:"name"`
		Result *struct {
			Name string `json:"name"`
		} `json:"result"`
		Stage *struct {
			Name string `json:"name"`
		} `json:"stage"`
	} `json:"state"`
}

type ListPipelinesOptions struct {
	Status string
	Sort   string
}

func (c *Client) ListPipelines(workspace, repoSlug string, opts *ListPipelinesOptions) ([]Pipeline, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Status != "" {
			params.Set("status", opts.Status)
		}
		if opts.Sort != "" {
			params.Set("sort", opts.Sort)
		}
	}
	params.Set("sort", "-created_on")
	params.Set("pagelen", "20")

	path := fmt.Sprintf("/repositories/%s/%s/pipelines",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	data, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	var pipelines []Pipeline
	if err := json.Unmarshal(page.Values, &pipelines); err != nil {
		return nil, fmt.Errorf("parsing pipelines: %w", err)
	}

	return pipelines, nil
}

func (c *Client) GetPipeline(workspace, repoSlug, pipelineUUID string) (*Pipeline, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pipelines/%s",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(pipelineUUID))

	data, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var pipeline Pipeline
	if err := json.Unmarshal(data, &pipeline); err != nil {
		return nil, fmt.Errorf("parsing pipeline: %w", err)
	}

	return &pipeline, nil
}

type RunPipelineRequest struct {
	Target struct {
		Type     string `json:"type"`
		RefType  string `json:"ref_type"`
		RefName  string `json:"ref_name"`
		Selector struct {
			Type    string `json:"type"`
			Pattern string `json:"pattern"`
		} `json:"selector,omitempty"`
	} `json:"target"`
}

func NewBranchPipelineRequest(branch string) *RunPipelineRequest {
	req := &RunPipelineRequest{}
	req.Target.Type = "pipeline_ref_target"
	req.Target.RefType = "branch"
	req.Target.RefName = branch
	return req
}

func NewCustomPipelineRequest(branch, pattern string) *RunPipelineRequest {
	req := NewBranchPipelineRequest(branch)
	req.Target.Selector.Type = "custom"
	req.Target.Selector.Pattern = pattern
	return req
}

func (c *Client) RunPipeline(workspace, repoSlug string, req *RunPipelineRequest) (*Pipeline, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/repositories/%s/%s/pipelines",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	data, err := c.post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var pipeline Pipeline
	if err := json.Unmarshal(data, &pipeline); err != nil {
		return nil, fmt.Errorf("parsing pipeline: %w", err)
	}
	return &pipeline, nil
}

func (c *Client) StopPipeline(workspace, repoSlug, pipelineUUID string) error {
	path := fmt.Sprintf("/repositories/%s/%s/pipelines/%s/stopPipeline",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(pipelineUUID))
	return c.postNoContent(path, nil)
}

type PipelineStep struct {
	UUID              string `json:"uuid"`
	Name              string `json:"name"`
	StartedOn         string `json:"started_on"`
	CompletedOn       string `json:"completed_on"`
	DurationInSeconds int    `json:"duration_in_seconds"`
	BuildSecondsUsed  int    `json:"build_seconds_used"`
	RunNumber         int    `json:"run_number"`
	MaxTime           int    `json:"max_time"`
	State             struct {
		Name   string `json:"name"`
		Result *struct {
			Name string `json:"name"`
		} `json:"result"`
	} `json:"state"`
	SetupCommands []struct {
		Name    string `json:"name"`
		Command string `json:"command"`
	} `json:"setup_commands"`
	ScriptCommands []struct {
		Name    string `json:"name"`
		Command string `json:"command"`
	} `json:"script_commands"`
	Image struct {
		Name string `json:"name"`
	} `json:"image"`
}

func (c *Client) ListPipelineSteps(workspace, repoSlug, pipelineUUID string) ([]PipelineStep, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pipelines/%s/steps",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(pipelineUUID))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var steps []PipelineStep
	if err := json.Unmarshal(page.Values, &steps); err != nil {
		return nil, fmt.Errorf("parsing steps: %w", err)
	}
	return steps, nil
}

func (c *Client) GetPipelineStep(workspace, repoSlug, pipelineUUID, stepUUID string) (*PipelineStep, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pipelines/%s/steps/%s",
		url.PathEscape(workspace), url.PathEscape(repoSlug),
		url.PathEscape(pipelineUUID), url.PathEscape(stepUUID))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var step PipelineStep
	if err := json.Unmarshal(data, &step); err != nil {
		return nil, fmt.Errorf("parsing step: %w", err)
	}
	return &step, nil
}

func (c *Client) GetStepLog(workspace, repoSlug, pipelineUUID, stepUUID string) (string, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pipelines/%s/steps/%s/log",
		url.PathEscape(workspace), url.PathEscape(repoSlug),
		url.PathEscape(pipelineUUID), url.PathEscape(stepUUID))
	data, err := c.get(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type PipelineVariable struct {
	UUID    string `json:"uuid"`
	Key     string `json:"key"`
	Value   string `json:"value"`
	Secured bool   `json:"secured"`
}

func (c *Client) ListPipelineVariables(workspace, repoSlug string) ([]PipelineVariable, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pipelines_config/variables",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var vars []PipelineVariable
	if err := json.Unmarshal(page.Values, &vars); err != nil {
		return nil, fmt.Errorf("parsing variables: %w", err)
	}
	return vars, nil
}

func (c *Client) CreatePipelineVariable(workspace, repoSlug string, key, value string, secured bool) (*PipelineVariable, error) {
	body := map[string]interface{}{
		"key":     key,
		"value":   value,
		"secured": secured,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/repositories/%s/%s/pipelines_config/variables",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	data, err := c.post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var v PipelineVariable
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("parsing variable: %w", err)
	}
	return &v, nil
}

func (c *Client) DeletePipelineVariable(workspace, repoSlug, variableUUID string) error {
	path := fmt.Sprintf("/repositories/%s/%s/pipelines_config/variables/%s",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(variableUUID))
	return c.deleteNoContent(path)
}
