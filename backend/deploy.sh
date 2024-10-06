env GOOS=linux GOARCH=amd64 go build main.go && gcloud compute scp main instance-20241005-084327:~/main --zone "europe-west3-a" --project "data-integration-development"
