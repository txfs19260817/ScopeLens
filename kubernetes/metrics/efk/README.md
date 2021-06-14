# EFK Stack

Install ECK CRDs first (https://www.elastic.co/guide/en/cloud-on-k8s/current/k8s-deploy-eck.html).

The applying order is `elasticsearch.yaml` -> `kibana.yaml` -> `fluentd-es-configmap.yaml` -> `fluentd-es-ds.yaml`.