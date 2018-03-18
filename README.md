# golang_demo
go语言学习的demo仓库

1. 平台区分
- 文件名_平台。
- 例： file_windows.go, file_unix.go
- 可选为：windows, unix, posix, plan9, darwin, bsd, linux, freebsd, nacl, netbsd, openbsd, solaris, dragonfly, bsd, notbsd， android，stubs


2. 测试单元
- 文件名_test.go或者 文件名_平台_test.go。
- 例： path_test.go,  path_windows_test.go


3. 版本区分(猜测)
- 文件名_版本号等。
- 例：trap_windows_1.4.go


4. CPU类型区分, 汇编用的多
- 文件名_(平台:可选)_CPU类型.
- 例：vdso_linux_amd64.go
- 可选：amd64, none, 386, arm, arm64, mips64, s390,mips64x,ppc64x, nonppc64x, s390x, x86,amd64p32


以上是根据go源码中收集整理的，难免有错。有些还未得到证实。
