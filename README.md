Simple IoT Backend using Go and hosted on Kubernetes

Folder Structure:

1. backend - The backend service. Implements a simple in-memory database using a map with RPC methods for receiving and sending data
2. frontend - The frontend service. Implements a HTTP server interacting with a JSON API. Interacts with the backend via RPC
3. yaml files- These files define the Kubernetes resources created
  3.1 backend.yaml and frontend.yaml - Define deployments of the app
  3.2 backend-service.yaml and frontend-service.yaml - Define services to expose the frontend service to the public and backend service to       the frontend service
