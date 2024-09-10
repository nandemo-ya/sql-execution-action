.PHONY: update-credits
update-credits:
	@go install github.com/Songmu/gocredits/cmd/gocredits@latest
	@gocredits -w . > CREDITS

.PHONY: update-docs
update-docs:
	@docker image build -t sql-execution-action-docs -f tools/action-docs/Dockerfile .
	@docker container run -v .:/app -it --rm sql-execution-action-docs --update-readme