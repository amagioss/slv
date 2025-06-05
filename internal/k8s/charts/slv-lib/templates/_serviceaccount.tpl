{{- define "slvlib.serviceaccount" -}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: slv-serviceaccount
  namespace: {{ .Release.Namespace }}
  {{- with .Values.serviceAccount.labels }}
  labels:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  {{- with .Values.serviceAccount.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
automountServiceAccountToken: true
{{- end -}}
