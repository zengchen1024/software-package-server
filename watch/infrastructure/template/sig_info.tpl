
- repo:
  - src-openeuler/{{.PkgName}}
  committers:
  {{- range .Committers }}
  - gitee_id: {{.GiteeId}}
    name: {{.Name}}
    email: {{.Email}}
  {{- end}}
