package pages

import (
	"go-templ-mock-1/components/badge"
	"go-templ-mock-1/components/chart"
)

// Dummy platform engineering data. Everything here is static mock data.

var Environments = []string{"dev", "test", "prod"}

type Deployment struct {
	Version    string
	Status     string // healthy, deploying, degraded, failed
	Replicas   string
	DeployedAt string
	By         string
	Commit     string
}

type Vulns struct {
	Critical, High, Medium, Low int
}

func (v Vulns) Total() int { return v.Critical + v.High + v.Medium + v.Low }

type Project struct {
	Slug, Name, Team, Language, Image string
	Deployments                       map[string]Deployment
	Vulns                             Vulns
	Coverage                          int
}

type DeployEvent struct {
	Project, Env, Version, Status, By, When, Duration string
}

type ErrorEvent struct {
	Project, Env, Message string
	Count                 int
	LastSeen              string
}

type CVE struct {
	ID, Project, Severity, Package, Installed, Fixed string
}

var AllProjects = []Project{
	{
		Slug: "payments-api", Name: "payments-api", Team: "Payments", Language: "Go", Image: "registry.internal/payments-api",
		Vulns: Vulns{Critical: 1, High: 3, Medium: 7, Low: 12}, Coverage: 82,
		Deployments: map[string]Deployment{
			"dev":  {"v2.14.0-rc.2", "healthy", "2/2", "12m ago", "ada", "a1b2c3d"},
			"test": {"v2.13.1", "healthy", "3/3", "3h ago", "ada", "9f8e7d6"},
			"prod": {"v2.13.0", "healthy", "6/6", "2d ago", "release-bot", "4c5d6e7"},
		},
	},
	{
		Slug: "checkout-web", Name: "checkout-web", Team: "Storefront", Language: "TypeScript", Image: "registry.internal/checkout-web",
		Vulns: Vulns{Critical: 0, High: 5, Medium: 14, Low: 9}, Coverage: 74,
		Deployments: map[string]Deployment{
			"dev":  {"v5.2.0", "deploying", "1/2", "1m ago", "grace", "e4f5a6b"},
			"test": {"v5.1.4", "healthy", "2/2", "1d ago", "grace", "c7d8e9f"},
			"prod": {"v5.1.3", "degraded", "4/5", "4d ago", "release-bot", "0a1b2c3"},
		},
	},
	{
		Slug: "inventory-service", Name: "inventory-service", Team: "Fulfilment", Language: "Go", Image: "registry.internal/inventory-service",
		Vulns: Vulns{Critical: 0, High: 1, Medium: 4, Low: 6}, Coverage: 88,
		Deployments: map[string]Deployment{
			"dev":  {"v1.9.3", "healthy", "1/1", "5h ago", "linus", "d3e4f5a"},
			"test": {"v1.9.3", "healthy", "2/2", "5h ago", "linus", "d3e4f5a"},
			"prod": {"v1.9.2", "healthy", "4/4", "6d ago", "release-bot", "b6c7d8e"},
		},
	},
	{
		Slug: "notifications-worker", Name: "notifications-worker", Team: "Platform", Language: "Python", Image: "registry.internal/notifications-worker",
		Vulns: Vulns{Critical: 2, High: 6, Medium: 11, Low: 20}, Coverage: 61,
		Deployments: map[string]Deployment{
			"dev":  {"v0.8.1", "failed", "0/1", "40m ago", "margaret", "f7a8b9c"},
			"test": {"v0.8.0", "healthy", "1/1", "2d ago", "margaret", "1d2e3f4"},
			"prod": {"v0.7.9", "healthy", "2/2", "9d ago", "release-bot", "5a6b7c8"},
		},
	},
	{
		Slug: "auth-gateway", Name: "auth-gateway", Team: "Identity", Language: "Rust", Image: "registry.internal/auth-gateway",
		Vulns: Vulns{Critical: 0, High: 0, Medium: 2, Low: 3}, Coverage: 93,
		Deployments: map[string]Deployment{
			"dev":  {"v3.0.0-beta.4", "healthy", "2/2", "1d ago", "dennis", "8b9c0d1"},
			"test": {"v2.8.2", "healthy", "3/3", "2d ago", "dennis", "2e3f4a5"},
			"prod": {"v2.8.2", "healthy", "8/8", "2d ago", "release-bot", "2e3f4a5"},
		},
	},
}

var DeployHistory = []DeployEvent{
	{"checkout-web", "dev", "v5.2.0", "deploying", "grace", "1m ago", "-"},
	{"payments-api", "dev", "v2.14.0-rc.2", "healthy", "ada", "12m ago", "1m 48s"},
	{"notifications-worker", "dev", "v0.8.1", "failed", "margaret", "40m ago", "2m 05s"},
	{"payments-api", "test", "v2.13.1", "healthy", "ada", "3h ago", "2m 12s"},
	{"inventory-service", "test", "v1.9.3", "healthy", "linus", "5h ago", "1m 31s"},
	{"inventory-service", "dev", "v1.9.3", "healthy", "linus", "5h ago", "1m 20s"},
	{"checkout-web", "test", "v5.1.4", "healthy", "grace", "1d ago", "3m 02s"},
	{"auth-gateway", "prod", "v2.8.2", "healthy", "release-bot", "2d ago", "6m 40s"},
	{"payments-api", "prod", "v2.13.0", "healthy", "release-bot", "2d ago", "5m 55s"},
	{"checkout-web", "prod", "v5.1.3", "degraded", "release-bot", "4d ago", "4m 18s"},
	{"inventory-service", "prod", "v1.9.2", "healthy", "release-bot", "6d ago", "3m 47s"},
}

var RecentErrors = []ErrorEvent{
	{"checkout-web", "prod", "TypeError: Cannot read properties of undefined (reading 'sku')", 482, "2m ago"},
	{"notifications-worker", "dev", "ConnectionRefusedError: redis:6379 connection refused", 117, "9m ago"},
	{"payments-api", "prod", "context deadline exceeded calling fraud-check", 38, "21m ago"},
	{"inventory-service", "test", "pq: deadlock detected while updating stock_levels", 6, "1h ago"},
	{"checkout-web", "prod", "502 Bad Gateway from upstream payments-api", 24, "2h ago"},
	{"auth-gateway", "test", "jwt: token signature invalid (key rotation lag)", 3, "5h ago"},
}

var CVEs = []CVE{
	{"CVE-2025-48734", "payments-api", "critical", "golang.org/x/crypto", "0.31.0", "0.35.0"},
	{"CVE-2025-22871", "payments-api", "high", "stdlib net/http", "1.23.4", "1.23.8"},
	{"CVE-2025-30204", "payments-api", "high", "github.com/golang-jwt/jwt", "5.2.1", "5.2.2"},
	{"CVE-2025-27152", "checkout-web", "high", "axios", "1.7.9", "1.8.2"},
	{"CVE-2025-29927", "checkout-web", "high", "next", "14.2.24", "14.2.25"},
	{"CVE-2024-45296", "checkout-web", "medium", "path-to-regexp", "0.1.7", "0.1.10"},
	{"CVE-2025-32415", "inventory-service", "high", "libxml2", "2.12.9", "2.12.10"},
	{"CVE-2025-32962", "notifications-worker", "critical", "jinja2", "3.1.4", "3.1.6"},
	{"CVE-2025-47273", "notifications-worker", "critical", "setuptools", "69.5.1", "78.1.1"},
	{"CVE-2025-43859", "notifications-worker", "high", "h11", "0.14.0", "0.16.0"},
	{"CVE-2025-24970", "auth-gateway", "medium", "openssl", "3.0.15", "3.0.16"},
}

func FindProject(slug string) (Project, bool) {
	for _, p := range AllProjects {
		if p.Slug == slug {
			return p, true
		}
	}
	return Project{}, false
}

// FilterHistory returns deploy events for env ("" or "all" means every env) and project ("" means all).
func FilterHistory(env, project string) []DeployEvent {
	var out []DeployEvent
	for _, e := range DeployHistory {
		if env != "" && env != "all" && e.Env != env {
			continue
		}
		if project != "" && e.Project != project {
			continue
		}
		out = append(out, e)
	}
	return out
}

func projectErrors(slug string) []ErrorEvent {
	var out []ErrorEvent
	for _, e := range RecentErrors {
		if e.Project == slug {
			out = append(out, e)
		}
	}
	return out
}

func projectCVEs(slug string) []CVE {
	var out []CVE
	for _, c := range CVEs {
		if c.Project == slug {
			out = append(out, c)
		}
	}
	return out
}

func totalVulns() Vulns {
	var t Vulns
	for _, p := range AllProjects {
		t.Critical += p.Vulns.Critical
		t.High += p.Vulns.High
		t.Medium += p.Vulns.Medium
		t.Low += p.Vulns.Low
	}
	return t
}

func countStatus(status string) int {
	n := 0
	for _, p := range AllProjects {
		for _, d := range p.Deployments {
			if d.Status == status {
				n++
			}
		}
	}
	return n
}

func totalErrors() int {
	n := 0
	for _, e := range RecentErrors {
		n += e.Count
	}
	return n
}

func statusVariantFor(status string) badge.Variant {
	switch status {
	case "healthy":
		return badge.VariantDefault
	case "deploying":
		return badge.VariantSecondary
	case "degraded":
		return badge.VariantOutline
	}
	return badge.VariantDestructive
}

// statusClass overrides the default (black) badge look so healthy reads green.
func statusClass(status string) string {
	if status == "healthy" {
		return "bg-green-600! text-white! dark:bg-green-500!"
	}
	return ""
}

func severityVariant(sev string) badge.Variant {
	switch sev {
	case "critical", "high":
		return badge.VariantDestructive
	case "medium":
		return badge.VariantSecondary
	}
	return badge.VariantOutline
}

// Chart data

var vulnConfig = chart.Config{
	{Key: "critical", Label: "Critical", Color: "#dc2626"},
	{Key: "high", Label: "High", Color: "#ea580c"},
	{Key: "medium", Label: "Medium", Color: "#eab308"},
	{Key: "low", Label: "Low", Color: "#94a3b8"},
}

func vulnChartData() []chart.Datum {
	data := make([]chart.Datum, len(AllProjects))
	for i, p := range AllProjects {
		data[i] = chart.Datum{
			"project": p.Name, "critical": p.Vulns.Critical, "high": p.Vulns.High,
			"medium": p.Vulns.Medium, "low": p.Vulns.Low,
		}
	}
	return data
}

var deployFreqConfig = chart.Config{
	{Key: "dev", Label: "dev", Color: "#60a5fa"},
	{Key: "test", Label: "test", Color: "#a78bfa"},
	{Key: "prod", Label: "prod", Color: "#16a34a"},
}

var deployFreqData = []chart.Datum{
	{"week": "W1", "dev": 41, "test": 22, "prod": 6},
	{"week": "W2", "dev": 38, "test": 25, "prod": 7},
	{"week": "W3", "dev": 52, "test": 28, "prod": 9},
	{"week": "W4", "dev": 47, "test": 24, "prod": 8},
	{"week": "W5", "dev": 58, "test": 31, "prod": 11},
	{"week": "W6", "dev": 63, "test": 33, "prod": 10},
}
