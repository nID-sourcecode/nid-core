Dit hoofdstuk beschrijft alle functionaliteit met betrekking tot infrastructuur voor het TWI platform.

## Versie beheer

Er wordt ontwikkeld middels Gitlab. In [Gitlab](https://about.gitlab.com/) is een CI/CD pipeline ingericht die de volgende stappen uitvoert:

1. Validate
   - Authorization policy generation
   - Go [WASM](https://github.com/envoyproxy/envoy-wasm) filters build
   - GraphQL generation
   - [Linting](https://github.com/golangci/golangci-lint)
   - Proto generation
   - Proto lint
   - Unit tests(248 unit tests)
2. Build services(nID-core: 12 micro services)
   - Build
   - Docker build
3. Release
   - [Go doc](https://godoc.org/) as pages
   - [gRPC transcoders](https://cloud.google.com/endpoints/docs/grpc/transcoding)
   - [Helm](https://helm.sh/)
   - [NPM](https://www.npmjs.com/) Package
   - [Swagger](https://swagger.io/) docs
   - WASM filters
4. Deploy
   - [Terraform](https://www.terraform.io/) plan
   - Terraform apply
   - Configure cluster
5. Integration test
   - Run 3x 28 integration tests firing 1000+ HTTP/gRPC requests

Voor issues/user stories worden development branches aangemaakt. Wanneer een developer een feature af denkt te hebben wordt hiervoor een merge request geopend. Hiervoor wordt de validate step van onze CI/CD stap gedraaid. Vervolgens moet er tenminste 1 andere developer een code review uitvoeren alvorens de nieuwe functionaliteit richting de master branch mag.
De master branch voert onze volledig CI/CD straat uit en deployed de code naar de staging omgeving op [Azure](https://portal.azure.com/).
Wanneer een nieuwe release gedaan wordt, wordt er vanaf master naar de prod branch een mr aangemaakt. Wanneer er akkoord wordt gegeven voor een nieuwe release zal de CI/CD straat dezelfde stappen uitvoeren, maar dan op het productie cluster.

## Cluster

Door middel van onze CI/CD worden [kubernetes](https://kubernetes.io/) clusters opgespint. Hierop wordt [Istio](https://istio.io/latest/) geïnstalleerd als service mesh. [Istiod](https://istio.io/latest/blog/2020/istiod/) zorgt voor [mTLS](https://en.wikipedia.org/wiki/Mutual_authentication) communicatie tussen alle services die draaien op het cluster. Zowel front als backend services draaien als [docker](https://www.docker.com/) containers binnen kubernetes. Load balancing van inkomend verkeer wordt geregeld door de Istio [ingress gateway](https://istio.io/latest/docs/tasks/traffic-management/ingress/ingress-control/). Met [cert-manager](https://cert-manager.io/docs/) wordt TLS verkeer geregeld. Uitgaand verkeer uit het cluster gaat via de [egress gateway](https://istio.io/latest/docs/tasks/traffic-management/egress/egress-gateway/).

Backend services zijn ingericht als [gRPC](https://grpc.io/) [microservices](https://en.wikipedia.org/wiki/Microservices) geïmplementeerd in [Go](https://golang.org/) en worden middels transcoding zowel als gRPC als REST JSON exposed. Daarnaast wordt standaard CRUD functionaliteit beschikbaar gemaakt middels GraphQL. Data wordt bewaard in [PostgreSQL](https://www.postgresql.org/) databases. Frontend services worden geïmplementeerd in [Typescript](https://www.typescriptlang.org/) middels de [Open WC](https://open-wc.org/) standaard.

## Metrics

Om metrics te kunnen bekijken binnen het cluster zijn de volgende tools geïnstalleerd:

- [Jaeger](https://www.jaegertracing.io/): Network tracing
- [Prometheus](https://prometheus.io/): Metric gathering
- [Grafana](https://grafana.com/): Metric dashboards
- [Kiali](https://kiali.io/): Inzicht service mesh
- [ELK](https://www.elastic.co/what-is/elk-stack): Logging
