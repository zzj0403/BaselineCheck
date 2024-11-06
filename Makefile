GO=go
APP=BaselineCheck

.PHONY: build run clean playbook check
	
build:
	$(GO) build -ldflags="-s -w"  -o ansible/roles/os-check/files/$(APP) main.go  && upx ansible/roles/os-check/files/$(APP)

run:
	$(GO) run main.go start 

start: build
	ansible/roles/os-check/files/$(APP) start 



clean:
	rm -f ansible/os-check/files/$(APP)

playbook: build
	ansible-playbook -i ansible/hosts  ansible/playbook/check-playbook.yml

check:
	ansible/roles/os-check/files/$(APP) check