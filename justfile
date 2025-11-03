default:
  @just --list

up: start-kind-cluster \
build-deps \
install-boxer \
wait-for-apps-to-be-ready \
install-ingress-controller \
create-ingress

fresh: stop up

stop:
    kind delete cluster --name kind

start-kind-cluster:
    kind create cluster --name kind --config=integration_tests/kind.yaml

check-kind-cluster:
    kubectl cluster-info
    kubectl get nodes
    kind get kubeconfig --name kind

build-deps:
    helm dependency build ./integration_tests/helm/setup

key := `openssl rand -base64 16 | tr -dc 'a-zA-Z0-9' | fold -w 16 | head -n 1`
install-boxer:
    helm upgrade --install --namespace default integration-tests integration_tests/helm/setup \
      --set boxer-issuer.issuer.replicas=1 \
      --set-literal 'boxer-validator-nginx.validator.config.tokenSettings.keys={"default": "{{key}}"}' \
      --set 'boxer-issuer.issuer.config.listenIp=0.0.0.0' \
      --set 'boxer-validator-nginx.validator.config.listenIp=0.0.0.0' \
      --set boxer-validator-nginx.validator.replicas=1

wait-for-apps-to-be-ready:
    for i in {1..30}; do \
      kubectl get pods -l app.kubernetes.io/name=boxer-issuer && \
      kubectl get pods -l app.kubernetes.io/name=boxer-validator-nginx && \
      break; \
      echo "Waiting for apps to be launched... ($i/30)"; \
      sleep 2; \
    done
    echo "Waiting for apps to be ready:"
    kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=boxer-issuer --timeout=120s
    kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=boxer-validator-nginx --timeout=120s

install-ingress-controller:
    helm upgrade --install ingress-nginx ingress-nginx \
      --repo https://kubernetes.github.io/ingress-nginx \
      --namespace ingress-nginx --create-namespace
    for i in {1..30}; do \
      kubectl get pods -l app.kubernetes.io/name=ingress-nginx && \
      break; \
      echo "Waiting for apps to be launched... ($i/30)"; \
      sleep 2; \
    done
    echo "Waiting for apps to be ready:"
    kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=ingress-nginx --timeout=120s

create-ingress:
    kubectl apply -f ./integration_tests/ingress.yaml
