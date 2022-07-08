# faas-image-converter
Faas Image generator converter


Fedora setup
```
export VIPS_VERSION=8.12.2
dnf install -y libwebp libwebp-devel expat expat-devel glib2-devel libaom-devel libavif-tools libavif-devel libglib2.0-dev openjpeg2 openjpeg2-devel libdav1d-devel libde265-devel libheif libheif-devel pkg-config x265-devel

wget https://github.com/libvips/libvips/releases/download/v${VIPS_VERSION}/vips-${VIPS_VERSION}.tar.gz
tar -xf vips-${VIPS_VERSION}.tar.gz
cd vips-${VIPS_VERSION}/
./configure --prefix=/usr/local
make
make install
```