DIRS=loginserver superserver routeserver gatewayserver recordserver test 

all: .proto .debug

test: .test

clean:
	
	@for dir in $(DIRS); \
	do \
		rm -f $$dir/$$dir; \
	done

	@rm -f command/*.pb.go
	@rm -rf $(GOPATH)/pkg

.debug:

	@for dir in $(DIRS); \
	do \
		go build -v -i -o $$dir/$$dir $$dir; \
	done
	

.test:

	@for dir in $(DIRS); \
	do \
		go build -race -v -i -o $$dir/$$dir $$dir; \
	done


.proto:
	@echo 编译proto文件
	@cd $(GOPATH)/proto/msg/ && chmod 777 build.sh && ./build.sh


