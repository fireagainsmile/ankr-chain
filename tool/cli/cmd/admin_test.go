package cmd

import (
	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	adminPrivate = "0mqsOtVueE7uq/I5J/dAhesumWXTu619xXuRgtj4l0d0ELMH6X9ZjGqT6Lnhrhp13LVeGIgrm3QgBnk4q16BZg=="
)

func TestSetValidator(t *testing.T) {
	convey.Convey("test setValidator", t, func() {
		//start test case
		args := []string{"admin", "setvalidator", "--pubkey", "FSyq/mTVPO/WdxMNCMEKiA5UVBFXVL8OAnDspO+buZY=","--power","20" ,"--nodeurl", "http://localhost:26657", "--privkey", adminPrivate}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestSetCert(t *testing.T) {
	convey.Convey("test setCert", t, func() {
		//start test case
		args := []string{"admin", "setcert", "--dcname", "dc-name","--perm","perm-string" ,"--nodeurl", "http://localhost:26657", "--privkey", adminPrivate}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestRemoveCert(t *testing.T) {
	convey.Convey("test removeCert", t, func() {

		args := []string{"admin", "removecert", "--dcname", "my-dcname","--nodeurl", "http://localhost:26657", "--privkey", adminPrivate}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}