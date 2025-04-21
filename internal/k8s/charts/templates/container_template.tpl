{{- define "slv.container" -}}
{{- if and (empty .Values.slvEnvironment.k8sSecret) (empty .Values.slvEnvironment.secretBinding) -}}
{{- fail "You must set at least one of slvEnvironment.k8sSecret or slvEnvironment.secretBinding" -}}
{{- end -}}
{{- $image := replace (trim .Values.runnerConfig.image) "\n" "" }}
{{- $expectedTag := replace (trim (printf ":%s" .Chart.Version)) "\n" "" }}
{{- if not (hasSuffix .Chart.Version .Values.runnerConfig.image ) -}}
{{- fail (printf "The image tag must be set to the Chart.Version '%s'" .Chart.Version) -}}
{{- end -}}
- name: slv
  image: {{ .Values.runnerConfig.image | default (printf "ghcr.io/amagioss/slv:%s" .Chart.Version) }}
  imagePullPolicy: {{ .Values.runnerConfig.imagePullPolicy | default "IfNotPresent" }}
  resources:
    {{- toYaml .Values.runnerConfig.resource | nindent 4 }}
  env:
    {{- if eq .Values.config.mode "operator" }}
    - name: SLV_MODE
      value: "k8s_operator"
    {{- end }}
    {{- if eq .Values.config.mode "job" }}
    - name: SLV_MODE
      value: "k8s_job"
    {{- end }}
    {{- if eq .Values.config.mode "cronjob" }}
    - name: SLV_MODE
      value: "k8s_job"
    {{- end }}
    {{- if and (false) (eq .Values.config.enableWebhook true) (eq .Values.config.mode "operator") }}
    - name: SLV_K8S_ENABLE_WEBHOOKS
      value: "true"
    {{- end }}
    {{- if ne .Values.slvEnvironment.k8sSecret "" | default "slv" }}
    - name: SLV_K8S_ENV_SECRET
      value: {{ .Values.slvEnvironment.k8sSecret }}
    {{- end }}
    {{- if ne .Values.slvEnvironment.secretBinding "" }}
    - name: SLV_ENV_SECRET_BINDING
      value: {{ .Values.slvEnvironment.secretBinding }}
    {{- end }}
{{- end }}
