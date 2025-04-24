{{- define "slvlib.clusterrole" -}}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: slv-clusterrole
rules:
  - apiGroups: ["slv.sh"]
    resources: ["slvs"]
    verbs:
      - "get"
      - "list"
      - "watch"
      - "update"
  - apiGroups: [""]
    resources: ["secrets"]
    verbs:
      - "create"
      - "get"
      - "list"
      - "update"
      - "delete"
      - "watch"
  - apiGroups: [""]
    resources: ["configmaps"]
    verbs:
      - "get"
      - "create"
      - "update"
{{- end -}}
