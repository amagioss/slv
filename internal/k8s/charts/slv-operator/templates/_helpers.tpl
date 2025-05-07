{{- define "installCertScript" -}}
- /bin/sh
- -c
- |
  set -e
  
  echo "Installing Packages..."
  apk add --no-cache openssl kubectl coreutils
  
  echo "Generating TLS certificates..."
  
  mkdir -p /certs && cd /certs
  
  CN="slv-operator-webhook.{{ .Release.Namespace }}.svc"
  DNS1="${CN}"
  DNS2="${CN}.cluster.local"
  DAYS={{ .Values.webhook.duration | default 10950 }} # 30 years
  VWC="slv-operator-validating-webhook"
  CERT_SECRET="slv-webhook-server-cert"
  RENEW_BEFORE={{ .Values.webhook.renewBefore | default 15 }} # days
  
  echo "🔍 Checking if certificate exists and is near expiry..."
  if kubectl get secret $CERT_SECRET -n {{ .Release.Namespace }} >/dev/null 2>&1; then
  kubectl get secret $CERT_SECRET -n {{ .Release.Namespace }} -o jsonpath="{.data.tls\.crt}" | base64 -d > current.crt
  EXPIRY=$(openssl x509 -enddate -noout -in current.crt | cut -d= -f2)
  EXPIRY_TS=$(date -d "$EXPIRY" +%s)
  NOW_TS=$(date +%s)
  DIFF_DAYS=$(( (EXPIRY_TS - NOW_TS) / 86400 ))
  echo "Cert expires in $DIFF_DAYS days"
  
  if [ "$DIFF_DAYS" -gt "$RENEW_BEFORE" ]; then
      echo "✅ Cert is still valid. Skipping regeneration."
      exit 0
  else
      echo "⚠️  Cert expiring soon. Regenerating..."
  fi
  else
  echo "⚠️  No cert found. Creating new one..."
  fi
  
  # 1. Generate CA
  openssl genrsa -out ca.key 2048
  openssl req -x509 -new -nodes -key ca.key -subj "/CN=slv-webhook-ca" -days $DAYS -out ca.crt
  
  # 2. Server key + CSR
  openssl genrsa -out tls.key 2048
  openssl req -new -key tls.key -subj "/CN=${CN}" -out server.csr
  
  # 3. CSR config
  cat > cert.conf <<EOF
  [req]
  distinguished_name = req_distinguished_name
  req_extensions = v3_req
  prompt = no
  
  [req_distinguished_name]
  CN = ${CN}
  
  [v3_req]
  keyUsage = keyEncipherment, digitalSignature
  extendedKeyUsage = serverAuth
  subjectAltName = @alt_names
  
  [alt_names]
  DNS.1 = ${DNS1}
  DNS.2 = ${DNS2}
  EOF
  
  # 4. Sign cert
  openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
  -out tls.crt -days $DAYS -extensions v3_req -extfile cert.conf
  
  # 5. Create TLS secret with ca.crt
  kubectl delete secret $CERT_SECRET -n {{ .Release.Namespace }} --ignore-not-found
  kubectl create secret generic $CERT_SECRET \
  --from-file=tls.crt --from-file=tls.key --from-file=ca.crt -n {{ .Release.Namespace }} \
  
  echo "TLS secret 'slv-webhook-server-cert' created successfully."
  
  CABUNDLE=$(base64 -w0 < ca.crt)
  
  kubectl patch validatingwebhookconfiguration $VWC \
  --type='json' \
  -p="[{\"op\": \"replace\", \"path\": \"/webhooks/0/clientConfig/caBundle\", \"value\":\"$CABUNDLE\"}]"
  echo "Webhook configuration patched successfully."
{{- end }}
