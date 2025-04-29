{{- define "slvlib.clusterrolebinding" -}}
{{- if eq .Values.serviceAccountName "" -}}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: slv-rolebinding
subjects:
- kind: ServiceAccount
  name: {{ .Values.serviceAccountName | default "slv-serviceaccount" }}
  namespace: {{ .Release.Namespace }}
roleRef:
  kind: ClusterRole
  name: slv-clusterrole
  apiGroup: rbac.authorization.k8s.io
{{- end -}}
{{- end -}}

