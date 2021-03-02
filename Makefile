build: clean
	go build -o _out/kubectl-schema cmd/kubectl-schema.go

clean:
	rm -rf _out
