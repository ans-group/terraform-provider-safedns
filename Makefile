testacc:
	TF_ACC=1 \
		go test -v \
		./safedns \
		-timeout 120m \
		-run=TestAcc${TEST}

testacc-all:
	TF_ACC=1 \
		go test -v \
		./safedns \
		-timeout 120m