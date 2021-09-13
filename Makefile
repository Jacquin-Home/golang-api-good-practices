all-migrations-up:
	#docker exec -it golang-api-good-practices-backend sh -c 'echo $$POSTGRES_PORT'
	docker exec -it golang-api-good-practices-backend sh -c 'migrate -database postgres://$$POSTGRES_USER:$$POSTGRES_PASSWORD@$$POSTGRES_HOST:$$POSTGRES_PORT/$$POSTGRES_DB?sslmode=disable -path store/migrations -verbose up'


hi:
	echo "hi"
	ls -l