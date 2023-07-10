integration-tests: stop-tests-dependencies start-tests-dependencies tests 
	$(MAKE) stop-tests-dependencies

tests:
	- go test -p 1 -v ./...

stop-tests-dependencies:
	docker-compose -f docker-compose-test.yml down

start-tests-dependencies:
	docker-compose -f docker-compose-test.yml up -d
