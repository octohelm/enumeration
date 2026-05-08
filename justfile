# Go 工具链入口
[group: 'toolchain']
mod go 'tool/go/justfile'

# 列出 root 与子模块的稳定执行入口
[group('meta')]
default:
    @just --list --list-submodules
