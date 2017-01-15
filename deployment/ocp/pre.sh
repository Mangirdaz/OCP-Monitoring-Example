#/bin/bash
oc import-image debug --from=docker.io/mangirdaz/docker-debug-container --confirm
oc new-app debug --image=debug

