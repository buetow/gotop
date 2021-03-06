install:
	go install github.com/snonux/gotop/gotop
run:
	go run gotop/main.go -m=2
docu: install
	sh -c '($(GOPATH)/bin/gotop -h 2>&1)|sed 1d > help.txt;exit 0'
cross:
	goxc -bc="linux"
pub:
	pub $(GOPATH)/bin/gotop
