# Ongi Backend Helm Chart

Kubernetes용 Ongi Backend API 서버 Helm Chart입니다.

## 사전 요구사항

- Kubernetes 1.19+
- Helm 3.0+
- PostgreSQL (내장 또는 외부)

## 설치

### 1. Docker 이미지 빌드

```bash
# 프로젝트 루트에서
docker build -t ongi-back:latest .

# 또는 태그 지정
docker build -t your-registry/ongi-back:v1.0.0 .
docker push your-registry/ongi-back:v1.0.0
```

### 2. Helm Chart 설치

```bash
# 기본 설치
helm install ongi-back ./helm/ongi-back

# 네임스페이스 지정
helm install ongi-back ./helm/ongi-back -n production --create-namespace

# values.yaml 오버라이드
helm install ongi-back ./helm/ongi-back -f custom-values.yaml
```

### 3. 특정 값만 오버라이드

```bash
helm install ongi-back ./helm/ongi-back \
  --set image.repository=your-registry/ongi-back \
  --set image.tag=v1.0.0 \
  --set secrets.DB_PASSWORD=123 \
  --set env.DB_HOST=external-postgres.example.com
```

## 설정

주요 설정값들은 `values.yaml`에서 확인할 수 있습니다:

### 이미지 설정

```yaml
image:
  repository: ongi-back
  tag: "latest"
  pullPolicy: IfNotPresent
```

### 데이터베이스 설정

```yaml
env:
  DB_HOST: postgres-service
  DB_PORT: "5432"
  DB_USER: ongi
  DB_NAME: ongi_db

secrets:
  DB_PASSWORD: "your-password"
```

### 오토스케일링

```yaml
autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80
```

### Ingress 설정

```yaml
ingress:
  enabled: true
  className: "nginx"
  hosts:
    - host: api.ongi.example.com
      paths:
        - path: /
          pathType: Prefix
```

## 업그레이드

```bash
# Chart 업그레이드
helm upgrade ongi-back ./helm/ongi-back

# 값 변경과 함께 업그레이드
helm upgrade ongi-back ./helm/ongi-back --set image.tag=v1.1.0
```

## 삭제

```bash
# Release 삭제
helm uninstall ongi-back

# 네임스페이스 지정
helm uninstall ongi-back -n production
```

## 상태 확인

```bash
# Release 상태
helm status ongi-back

# 배포된 리소스 확인
kubectl get all -l app.kubernetes.io/name=ongi-back

# Pod 로그 확인
kubectl logs -l app.kubernetes.io/name=ongi-back -f
```

## 프로덕션 배포 예시

### custom-prod-values.yaml

```yaml
replicaCount: 3

image:
  repository: your-registry/ongi-back
  tag: "v1.0.0"
  pullPolicy: Always

secrets:
  DB_PASSWORD: "super-secure-password"

env:
  DB_HOST: "postgres-prod.example.com"
  ENVIRONMENT: "production"

ingress:
  enabled: true
  hosts:
    - host: api.ongi.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: ongi-tls
      hosts:
        - api.ongi.com

resources:
  limits:
    cpu: 1000m
    memory: 1Gi
  requests:
    cpu: 500m
    memory: 512Mi

autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 20
```

배포:

```bash
helm install ongi-back ./helm/ongi-back \
  -f custom-prod-values.yaml \
  -n production \
  --create-namespace
```

## 문제 해결

### Pod가 시작되지 않는 경우

```bash
# Pod 상태 확인
kubectl get pods -l app.kubernetes.io/name=ongi-back

# 이벤트 확인
kubectl describe pod <pod-name>

# 로그 확인
kubectl logs <pod-name>
```

### 데이터베이스 연결 문제

```bash
# Secret 확인
kubectl get secret ongi-back-secret -o yaml

# 환경변수 확인
kubectl exec <pod-name> -- env | grep DB
```

## 내장 PostgreSQL 사용

```yaml
postgresql:
  enabled: true
  auth:
    username: ongi
    password: ongi123
    database: ongi_db
  primary:
    persistence:
      enabled: true
      size: 20Gi
```

## 외부 PostgreSQL 사용

```yaml
postgresql:
  enabled: false

env:
  DB_HOST: "external-postgres.example.com"
  DB_PORT: "5432"
  DB_USER: "ongi"
  DB_NAME: "ongi_db"

secrets:
  DB_PASSWORD: "external-db-password"
```
