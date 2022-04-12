COVERAGE_FILENAME=coverage.txt
BUILD_DIR=build
EXECUTABLE_NAME=healthcheck
GIT_HASH=$(shell git rev-parse HEAD)

rebuild: clean
	@echo "Making ${BUILD_DIR} directory"
	@mkdir ${BUILD_DIR}
	@echo "Building"
	@go build -o ./${BUILD_DIR}/${EXECUTABLE_NAME} ./cmd/healthcheck
	@echo "\033[32;1m>>> Built\033[0m"

cover: test
	@echo "Displaying results"
	@go tool cover -html ${COVERAGE_FILENAME}

run: rebuild
	@echo "Starting the application"
	@./${BUILD_DIR}/${EXECUTABLE_NAME}

test:
	@echo "Running tests with coverage"
	@go test -cover ./cmd/healthcheck -coverprofile ${COVERAGE_FILENAME}
	@echo "\033[32;1m>>> Tested\033[0m"

container:
	@echo "Building the healthcheck with tag '${GIT_HASH}'"
	@docker build -t healthcheck:"${GIT_HASH}" -f Dockerfile .

clean:
	@echo "Removing ${BUILD_DIR} directory"
	@rm -rf ./${BUILD_DIR}/
	@echo "Removing ${COVERAGE_FILENAME}"
	@rm -f ${COVERAGE_FILENAME}
	@echo "\033[32;1m>>> Cleared\033[0m"
