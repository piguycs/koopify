# Koopify

A sassy little web app to transfer funds from your wallet to mine and transfer snake oil from my non-existant inventory to yours!

This is a university project, I will try to keep things simple while also over-engineering where possible. The backend is in Go, and the frontend a SPA website written in Vue and Typescript which is transpiled to Javascript which generates HTML after reading JSON which is fetched after making a request to the backend. Could I have just sent the final JSON to the client along with some alpine or htmx code? No, because Vue is required for a passing grade.

## Deployments

CI/CD is setup to automatically build container images and publish them to the GitHub Container Registry. Server-side, `podman auto-update` is used for continuous delivery. Zero-downtime is not guranteed in the current setup, as only one node is used with no orchestrator.

## Local development

TODO. Probably going to be docker-compose, so it works with most popular OCI runtimes like podman (sorry lxc).
