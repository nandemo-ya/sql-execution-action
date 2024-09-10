.PHONY: update-credits
update-credits:
	@go install github.com/Songmu/gocredits/cmd/gocredits@latest
	@gocredits -w . > CREDITS