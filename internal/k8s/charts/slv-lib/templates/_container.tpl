{{- define "slvlib.container" -}}
{{- if and (.Values.image) (not (hasSuffix .Chart.Version .Values.image)) -}}
{{- fail (printf "The image tag must be set to the Chart.Version '%s'" .Chart.Version) -}}
{{- end -}}
- name: slv
  image: {{ .Values.image | default (printf "ghcr.io/amagioss/slv:%s" .Chart.AppVersion) }}
  imagePullPolicy: {{ .Values.imagePullPolicy | default "IfNotPresent" }}
  resources:
    {{- toYaml .Values.resource | nindent 4 }}
  env:
    {{- if ne .Values.k8sSecret ""}}
    - name: SLV_K8S_ENV_SECRET
      value: {{ .Values.k8sSecret }}
    {{- end }}
    {{- if ne .Values.secretBinding "" }}
    - name: SLV_ENV_SECRET_BINDING
      value: {{ .Values.secretBinding }}
    {{- end }}
{{- end }}
