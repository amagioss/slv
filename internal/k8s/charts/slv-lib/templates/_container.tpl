{{- define "slvlib.container" -}}
{{- if and (empty .Values.slvEnvironment.k8sSecret) (empty .Values.slvEnvironment.secretBinding) -}}
{{- fail "You must set at least one of slvEnvironment.k8sSecret or slvEnvironment.secretBinding" -}}
{{- end -}}
{{- $image := replace (trim .Values.config.image) "\n" "" }}
{{- $expectedTag := replace (trim (printf ":%s" .Chart.Version)) "\n" "" }}
{{- if not (hasSuffix .Chart.Version .Values.config.image ) -}}
{{- fail (printf "The image tag must be set to the Chart.Version '%s'" .Chart.Version) -}}
{{- end -}}
- name: slv
  image: {{ .Values.config.image | default (printf "ghcr.io/amagioss/slv:%s" .Chart.AppVersion) }}
  imagePullPolicy: {{ .Values.config.imagePullPolicy | default "IfNotPresent" }}
  resources:
    {{- toYaml .Values.config.resource | nindent 4 }}
  env:
    {{- if ne .Values.slvEnvironment.k8sSecret "" | default "slv" }}
    - name: SLV_K8S_ENV_SECRET
      value: {{ .Values.slvEnvironment.k8sSecret }}
    {{- end }}
    {{- if ne .Values.slvEnvironment.secretBinding "" }}
    - name: SLV_ENV_SECRET_BINDING
      value: {{ .Values.slvEnvironment.secretBinding }}
    {{- end }}
{{- end }}
