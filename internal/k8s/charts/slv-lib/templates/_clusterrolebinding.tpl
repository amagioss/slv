{{- define "slvlib.clusterrolebinding" -}}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: slv-rolebinding
subjects:
- kind: ServiceAccount
  name: slv-serviceaccount
  namespace: {{ .Release.Namespace }}
roleRef:
  kind: ClusterRole
  name: slv-clusterrole
  apiGroup: rbac.authorization.k8s.io
{{- end -}}
