build:
	@go mod tidy
	@go build -o target/pomodoro

run: build
	@./target/pomodoro

clean:
	rm -rf target/
