SHELL:=/bin/bash

proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. \
  		--go-grpc_opt=paths=source_relative article-service/pb/article/*.proto
	@ls article-service/pb/article/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'