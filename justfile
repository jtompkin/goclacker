version := "v1.4.3"
progname := "goclacker"
tmpdir := `mktemp -d`
tardir := tmpdir / progname + "-" + version
tarball := tardir + ".tar.gz"
releasedir := "release" / version

release:
    rm -f {{ tarball }}
    mkdir {{ tardir }}
    cp -R README.md CHANGELOG.md LICENSE go.mod go.sum *.go internal {{ tardir }}
    tar czvf {{ tarball }} -C {{ tmpdir }} {{ progname + "-" + version }}
    mkdir -p {{ releasedir }}
    cp {{ tarball }} {{ releasedir }}
    rm -rf {{ tarball }} {{ tardir }}
    rmdir {{ tmpdir }}
