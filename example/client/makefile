# 定义变量
PROTO_DIR=../proto
GEN_DIR=gen/user
PROTOC=protoc
PROTOC_GEN_GO=protoc-gen-go
PROTOC_GEN_GO_GRPC=protoc-gen-go-grpc

# 默认目标
all: generate

# 生成 Go 和 gRPC 代码
generate:
	@echo "Generating Go and gRPC code..."
	@mkdir -p $(GEN_DIR)
	$(PROTOC) --go_out=$(GEN_DIR) --go-grpc_out=$(GEN_DIR) --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative -I $(PROTO_DIR) $(PROTO_DIR)/*.proto
	@echo "Code generation complete!"

# 清理生成的文件
clean:
	@echo "Cleaning generated files..."
	@rm -rf $(GEN_DIR)
	@echo "Clean complete!"

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  all       - Alias for generate"
	@echo "  generate  - Generate Go and gRPC code"
	@echo "  clean     - Clean generated files"
	@echo "  help      - Show this help message"

.PHONY: all generate clean help