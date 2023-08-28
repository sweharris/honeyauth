SRC:=$(shell echo *.go)

honeyauth: $(SRC)
	go build

keys: ssh_host_rsa_key ssh_host_ecdsa_key ssh_host_ed25519_key

ssh_host_rsa_key:
	ssh-keygen -N '' -f $@ -t rsa

ssh_host_ecdsa_key:
	ssh-keygen -N '' -f $@ -t ecdsa

ssh_host_ed25519_key:
	ssh-keygen -N '' -f $@ -t ed25519

