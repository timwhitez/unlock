# Unlock

- readme.md 由ai生成

Unlock 是一个用 Go 语言编写的工具,用于读取 Windows 系统上被锁定的文件。

## 功能特点

- 解锁被其他进程占用的文件
- 支持命令行参数,可以指定要解锁的文件路径
- 使用 Windows API 和系统调用来实现文件解锁功能

## 使用方法

1. 编译程序:
   ```
   go build -o unlock.exe main.go
   ```

2. 运行程序:
   ```
   unlock.exe [文件路径]
   ```
   如果不提供文件路径参数,程序将提示用户输入。

## 注意事项

- 此工具仅适用于 Windows 操作系统
- 解锁过程可能会影响正在使用该文件的其他程序

## 依赖项

- Go 1.15 或更高版本
- Windows 系统 API

