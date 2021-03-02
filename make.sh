#bin/sh
go build
rm ~/Work/Go/project/bin/terraform-provider-ksyun
cp ~/Work/Go/project/src/github.com/kingsoftcloud/terraform-provider-ksyun/terraform-provider-ksyun ~/Work/Go/project/bin/
