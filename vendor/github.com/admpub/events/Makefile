VERSION=v$(shell cat VERSION)

test:
	ginkgo -r -cover -race -progress -keepGoing -randomizeAllSpecs -slowSpecThreshold 5 -trace

release:
	git tag $(VERSION) --message "release $(VERSION) ($(shell date '+%Y-%m-%d'))"
	git push $(VERSION)