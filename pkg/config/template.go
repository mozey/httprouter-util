// Code generated with https://github.com/mozey/config DO NOT EDIT

package config

import (
	"bytes"
	"text/template"
)

// ExecTemplateClientDownloadUrl fills APP_TEMPLATE_CLIENT_DOWNLOAD_URL with the given params
func (c *Config) ExecTemplateClientDownloadUrl(token string) string {
	t := template.Must(template.New("templateClientDownloadUrl").Parse(c.templateClientDownloadUrl))
	b := bytes.Buffer{}
	_ = t.Execute(&b, map[string]interface{}{

		"Token": token,
	})
	return b.String()
}

// ExecTemplateClientVersionUrl fills APP_TEMPLATE_CLIENT_VERSION_URL with the given params
func (c *Config) ExecTemplateClientVersionUrl(token string) string {
	t := template.Must(template.New("templateClientVersionUrl").Parse(c.templateClientVersionUrl))
	b := bytes.Buffer{}
	_ = t.Execute(&b, map[string]interface{}{

		"Token": token,
	})
	return b.String()
}
