#!/usr/bin/bash
builddir="bin"
windowsarches=(amd64)
macarches=(amd64)
linuxarches=(amd64)



for arch in ${windowsarches[@]}
do
    echo "Building for Windows - ${arch}" 
    target=${builddir}/drivvnapi_windows_${arch}.exe
	env GOOS=windows GOARCH=${arch} go build -o bin/drivvnapi_windows_${arch}.exe
    echo "Built for Windows to ${target}"
done

for arch in ${linuxarches[@]}
do
    echo "Building for Linux - ${arch}" 
    target=${builddir}/drivvnapi_linux_${arch}
	env GOOS=linux GOARCH=${arch} go build -o bin/drivvnapi_linux_${arch}
    echo "Built for Linux to ${target}"
done

for arch in ${macarches[@]}
do
    echo "Building for Mac - ${arch}" 
    target=${builddir}/drivvnapi_mac_${arch}
	env GOOS=darwin GOARCH=${arch} go build -o bin/drivvnapi_darwin_${arch}
    echo "Built for Mac to ${target}"
done

echo "Done!"