bootstrap: teardown
	docker-compose up -d grafana
	@echo "Go to http://localhost:3000/"

teardown:
	docker-compose down --remove-orphans --volumes --timeout=2
