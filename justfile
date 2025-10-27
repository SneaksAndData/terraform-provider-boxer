default:
  @just --list

up: start-kind-cluster \
deploy-ingress-nginx-controller \
build-deps \
install-boxer \
wait-for-apps-to-be-ready \
create-ingress

fresh: stop-kind-cluster up

stop: stop-kind-cluster

stop-kind-cluster:
    kind delete cluster --name kind

start-kind-cluster:
    kind create cluster --name kind --config=integration_tests/kind.yaml

check-kind-cluster:
    kubectl cluster-info
    kubectl get nodes
    kind get kubeconfig --name kind

deploy-ingress-nginx-controller:
    kubectl apply -f https://kind.sigs.k8s.io/examples/ingress/deploy-ingress-nginx.yaml

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
    kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=boxer-issuer --timeout=120s
    kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=boxer-validator-nginx --timeout=120s

create-ingress:
    kubectl apply -f ./integration_tests/ingress.yaml
