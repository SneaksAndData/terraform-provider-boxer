default:
  @just --list

up: start-kind-cluster \
configure-rbac \
build-deps \
install-boxer \
install-ingress-controller \
wait-for-services \
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

wait-for-services:
    kubectl rollout status deployment/boxer-issuer --timeout=180s
    kubectl rollout status deployment/boxer-validator-nginx --timeout=180s
    kubectl rollout status deployment/ingress-nginx-controller --namespace ingress-nginx --timeout=180s

install-ingress-controller:
    kubectl apply -f https://kind.sigs.k8s.io/examples/ingress/deploy-ingress-nginx.yaml

create-ingress:
    kubectl apply -f ./integration_tests/ingress.yaml

configure-rbac: # see:  https://github.com/kubernetes/kubernetes/issues/130781 for details
    kubectl apply -f ./integration_tests/rbac.yaml