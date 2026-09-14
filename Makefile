.PHONY: api
api:
	go build -o /home/ankit/Studies/coding/projects/distributed-media-processing-platform/services/api/api /home/ankit/Studies/coding/projects/distributed-media-processing-platform/services/api

.PHONY: auth
auth:
	go build -o /home/ankit/Studies/coding/projects/distributed-media-processing-platform/services/auth/auth /home/ankit/Studies/coding/projects/distributed-media-processing-platform/services/auth

auth-run: auth
	cd /home/ankit/Studies/coding/projects/distributed-media-processing-platform/services/auth && \
	/home/ankit/Studies/coding/projects/distributed-media-processing-platform/services/auth/auth

api-run: api
	cd /home/ankit/Studies/coding/projects/distributed-media-processing-platform/services/api && \
	/home/ankit/Studies/coding/projects/distributed-media-processing-platform/services/api/api
