# pwstore
a bad project for storing pws (please use with caution!)

url:
- https://pwstore-ydrkzrffeq-lz.a.run.app/

## Tech stack

Framework:
- [Gofiber](https://gofiber.io/)
- [Goth](https://github.com/markbates/goth)
- [Tailwindcss](https://tailwindcss.com/)
- [HTMX](https://htmx.org/)

Services:
- Cockroachdb: PSQL
  - https://cockroachlabs.cloud/cluster/500fd670-179e-4cd5-a528-5d4da6ed3b5c/overview?cluster-type=serverless
- GCP: cloud run
- GCP: Artifact Registry
  - europe-north1-docker.pkg.dev/pwstore/pwstore/pwstore:latest
- AWS: SES

