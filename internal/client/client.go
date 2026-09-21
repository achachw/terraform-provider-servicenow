package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
}

type ServiceOffering struct {
	SysID             string `json:"sys_id"`
	Name              string `json:"name"`
	ShortDescription  string `json:"short_description"`
	OperationalStatus string `json:"operational_status"`
	OwnedBy           string `json:"owned_by"`
	ManagedBy         string `json:"managed_by"`
	Raw               map[string]any
}

type UserGroup struct {
	SysID       string `json:"sys_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Manager     string `json:"manager"`
	Parent      string `json:"parent"`
	Active      string `json:"active"`
	Raw         map[string]any
}

type tableResponse struct {
	Result []map[string]any `json:"result"`
}

func New(baseURL, username, password string) (*Client, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return nil, fmt.Errorf("base_url is required")
	}
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}
	if password == "" {
		return nil, fmt.Errorf("password is required")
	}

	return &Client{
		baseURL:  baseURL,
		username: username,
		password: password,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (c *Client) GetServiceOfferingByName(ctx context.Context, name, view string) (*ServiceOffering, error) {
	offering, lookupURL, err := c.getSingleServiceOffering(ctx, "name="+name, view)
	if err != nil {
		return nil, err
	}

	if offering == nil {
		matches, matchErr := c.findServiceOfferings(ctx, "nameLIKE"+name, view, 5)
		if matchErr != nil {
			return nil, fmt.Errorf("service offering %q not found using URL %s; also failed to search similar names: %w", name, lookupURL, matchErr)
		}
		if len(matches) == 0 {
			return nil, fmt.Errorf("service offering %q not found using URL %s; no similar service offerings found with nameLIKE", name, lookupURL)
		}

		return nil, fmt.Errorf("service offering %q not found using URL %s; similar service offerings: %s", name, lookupURL, strings.Join(matches, ", "))
	}

	return offering, nil
}

func (c *Client) GetUserGroupByName(ctx context.Context, name, view string) (*UserGroup, error) {
	group, lookupURL, err := c.getSingleUserGroup(ctx, "name="+name, view)
	if err != nil {
		return nil, err
	}

	if group == nil {
		matches, matchErr := c.findUserGroups(ctx, "nameLIKE"+name, view, 5)
		if matchErr != nil {
			return nil, fmt.Errorf("user group %q not found using URL %s; also failed to search similar names: %w", name, lookupURL, matchErr)
		}
		if len(matches) == 0 {
			return nil, fmt.Errorf("user group %q not found using URL %s; no similar user groups found with nameLIKE", name, lookupURL)
		}

		return nil, fmt.Errorf("user group %q not found using URL %s; similar user groups: %s", name, lookupURL, strings.Join(matches, ", "))
	}

	return group, nil
}

func (c *Client) getSingleServiceOffering(ctx context.Context, sysparmQuery, view string) (*ServiceOffering, string, error) {
	rows, requestURL, err := c.queryTable(ctx, "service_offering", sysparmQuery, view, 0)
	if err != nil {
		return nil, requestURL, err
	}

	if len(rows) == 0 {
		return nil, requestURL, nil
	}
	if len(rows) > 1 {
		return nil, requestURL, fmt.Errorf("service offering query %q returned %d results, expected exactly one", sysparmQuery, len(rows))
	}

	return serviceOfferingFromRow(rows[0]), requestURL, nil
}

func (c *Client) getSingleUserGroup(ctx context.Context, sysparmQuery, view string) (*UserGroup, string, error) {
	rows, requestURL, err := c.queryTable(ctx, "sys_user_group", sysparmQuery, view, 0)
	if err != nil {
		return nil, requestURL, err
	}

	if len(rows) == 0 {
		return nil, requestURL, nil
	}
	if len(rows) > 1 {
		return nil, requestURL, fmt.Errorf("user group query %q returned %d results, expected exactly one", sysparmQuery, len(rows))
	}

	return userGroupFromRow(rows[0]), requestURL, nil
}

func (c *Client) findServiceOfferings(ctx context.Context, sysparmQuery, view string, limit int) ([]string, error) {
	rows, _, err := c.queryTable(ctx, "service_offering", sysparmQuery, view, limit)
	if err != nil {
		return nil, err
	}

	matches := make([]string, 0, len(rows))
	for _, row := range rows {
		name := stringValue(row["name"])
		sysID := stringValue(row["sys_id"])
		if sysID != "" {
			matches = append(matches, fmt.Sprintf("%s (%s)", name, sysID))
		} else {
			matches = append(matches, name)
		}
	}

	return matches, nil
}

func (c *Client) findUserGroups(ctx context.Context, sysparmQuery, view string, limit int) ([]string, error) {
	rows, _, err := c.queryTable(ctx, "sys_user_group", sysparmQuery, view, limit)
	if err != nil {
		return nil, err
	}

	matches := make([]string, 0, len(rows))
	for _, row := range rows {
		name := stringValue(row["name"])
		sysID := stringValue(row["sys_id"])
		if sysID != "" {
			matches = append(matches, fmt.Sprintf("%s (%s)", name, sysID))
		} else {
			matches = append(matches, name)
		}
	}

	return matches, nil
}

func (c *Client) queryTable(ctx context.Context, table, sysparmQuery, view string, limit int) ([]map[string]any, string, error) {
	endpoint, err := url.Parse(c.baseURL + "/api/now/table/" + table)
	if err != nil {
		return nil, "", fmt.Errorf("invalid ServiceNow URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("sysparm_query", sysparmQuery)
	if view != "" {
		query.Set("sysparm_view", view)
	}
	if limit > 0 {
		query.Set("sysparm_limit", fmt.Sprint(limit))
	}
	endpoint.RawQuery = encodeQuery(query)
	requestURL := endpoint.Redacted()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, requestURL, fmt.Errorf("create ServiceNow request: %w", err)
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, requestURL, fmt.Errorf("call ServiceNow %s table: %w", table, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, requestURL, fmt.Errorf("ServiceNow returned HTTP %d for %s lookup", resp.StatusCode, table)
	}

	var payload tableResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, requestURL, fmt.Errorf("decode ServiceNow response: %w", err)
	}

	return payload.Result, requestURL, nil
}

func serviceOfferingFromRow(row map[string]any) *ServiceOffering {
	return &ServiceOffering{
		SysID:             stringValue(row["sys_id"]),
		Name:              stringValue(row["name"]),
		ShortDescription:  stringValue(row["short_description"]),
		OperationalStatus: stringValue(row["operational_status"]),
		OwnedBy:           stringValue(row["owned_by"]),
		ManagedBy:         stringValue(row["managed_by"]),
		Raw:               row,
	}
}

func userGroupFromRow(row map[string]any) *UserGroup {
	return &UserGroup{
		SysID:       stringValue(row["sys_id"]),
		Name:        stringValue(row["name"]),
		Description: stringValue(row["description"]),
		Email:       stringValue(row["email"]),
		Manager:     stringValue(row["manager"]),
		Parent:      stringValue(row["parent"]),
		Active:      stringValue(row["active"]),
		Raw:         row,
	}
}

func encodeQuery(query url.Values) string {
	// ServiceNow sysparm_query is more reliable with spaces encoded as %20.
	return strings.ReplaceAll(query.Encode(), "+", "%20")
}

func stringValue(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case map[string]any:
		if displayValue, ok := typed["display_value"].(string); ok {
			return displayValue
		}
		if rawValue, ok := typed["value"].(string); ok {
			return rawValue
		}
	}

	if value == nil {
		return ""
	}

	return fmt.Sprint(value)
}
