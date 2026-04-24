mod go 'tool/go/justfile'

# 列出 root 与子模块的稳定执行入口。
default:
    just --list --list-submodules
