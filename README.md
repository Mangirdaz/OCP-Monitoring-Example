## OCP Monitoring Demo

This is test application to show how we can use Kubernetes Liveness and Readiness probes to monitor our application and application relationships.
Application has HARD dependencies on KV Storage and API to be available. But external content is not one, which need to be available all the time.
Application can work and serve customers without it. 

### Application architecture/desing
    FrontEnd(FE) <-api-> Backend <-kv interface -> Storage(BoltDB, Consul, ETCD)
                                 <-    api      -> External static content (Web server)


#### Api Backend Endpoints
    (POST, GET) /api/v1/notes               - To manage notes on our storage via api
    (GET)               /healthz            - To check if container is healthy (Liveness probes) 
    (GET)               /readiness          - Check if Container is ready to serve trafic (Itself + downstream dependencies)

### External Service 
    (GET) /api/v1/img        - Serve static content for app 
    (GET) /healthz              - To check if container is healthy (Liveness probes) 
    
## Available Commands
    make build                - build all things
    make build-api            - build api 
    make build-api-ext        - build api external
    make build-ui             - build ui
    make run-api              - run api
    make run-api-ext          - run api external
    make run-consul           - run consul
    make run-fe               - run FrontEnd
    make run-ui               - run ui
    make docker-build-fe      - build docker image for frontend 
    make docker-build-api     - build docker image for api 
    make docker-build-api-ext     - build docker image for external api


## Images
    mangirdas/ocp-example-fe:v0.5
    mangirdas/ocp-example-api:v0.5
    mangirdas/ocp-example-api-ext:v0.5

## How To Run:

    oc cluster up
    oc create -f deployments/ocp/
    ./deployments/ocp/post.sh 

    Because we use statik for running static content via GO App YOU will need to update /frontend/ui/config/prod.env.js file with API and IMG api endpoints. 

## Rebuild UI:
   make build-ui
   make docker-build-fe
