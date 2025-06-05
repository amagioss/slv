{{- define "slvlib.notes" -}}
SLV Install Successful.

{{- if and (empty .Values.k8sSecret) (empty .Values.secretBinding) }}
WARNING: You have not set the value for ".Values.secretBinding" or "Values.slvEnvironment.k8sSecret".
SLV will now look for a secret named "slv" in the "{{ .Release.Namespace }}" namespace.
If a secret is not found, SLV will not run as expected and return an error.

Ensure that you have set atleast one of the following
- secret key for the environment (under key "SecretKey") 
- secret binding for the environment (under key "SecretBinding"),
under the secret name "slv" 
in namespace "{{ .Release.Namespace }}"
{{- end -}}

{{- if ne .Values.k8sSecret "" }}
SLV will get the environment secret key/binding from a preloaded kubernetes secret.

Ensure that you have set atleast one of the following
- secret binding for the environment (under key "SecretBinding"),
- secret key for the environment (under key "SecretKey")
under the secret name "{{ .Values.k8sSecret }}"
in namespace "{{ .Release.Namespace }}"
{{- end }}

{{- end }}
