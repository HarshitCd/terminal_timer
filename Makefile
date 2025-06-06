AUDIO_PATH?="$(HOME)/.config/terminal_timer/audio/time_up.mp3"

build:
	@go mod tidy
	@go build -ldflags="-X 'main.audioPath=$(AUDIO_PATH)'" -o ./target/terminal_timer

install:
	@mkdir -p ~/.config/terminal_timer
	@cp -R audio ~/.config/terminal_timer/

run: build
	@./target/pomodoro

clean:
	rm -rf target/
