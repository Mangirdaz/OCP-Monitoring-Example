#/bin/bash
oc expose service fe
oc expose service api-svc
oc expose service img-svc