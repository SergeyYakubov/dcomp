package utils

import (
	"fmt"
	"log"
	"os"

	"encoding/base64"

	"strings"

	"github.com/apcera/gssapi"
	"github.com/pkg/errors"
)

type Context struct {
	DebugLog       bool
	ServiceName    string
	ServiceAddress string

	gssapi.Options

	Lib       *gssapi.Lib
	libLoaded bool

	// Service credentials loaded from keytab
	credential *gssapi.CredId
}

func (c *Context) loadLib(debug bool, prefix string) (err error) {
	if c.libLoaded {
		return nil
	}
	max := gssapi.Err + 1
	if debug {
		max = gssapi.MaxSeverity
	}
	pp := make([]gssapi.Printer, 0, max)
	for i := gssapi.Severity(0); i < max; i++ {
		p := log.New(os.Stderr,
			fmt.Sprintf("%s: %s\t", prefix, i),
			log.LstdFlags)
		pp = append(pp, p)
	}
	c.Options.Printers = pp
	c.Lib, err = gssapi.Load(&c.Options)

	if err == nil {
		c.libLoaded = true
	}
	return
}

func (c *Context) prepareServiceName(name string) (*gssapi.Name, error) {

	if !c.libLoaded {
		return nil, errors.New("GSSAPI library not loaded")
	}

	if name == "" {
		return nil, errors.New("Need a --service-name")
	}

	nameBuf, err := c.Lib.MakeBufferString(name)
	if err != nil {
		return nil, err
	}
	defer nameBuf.Release()

	gssapiName, err := nameBuf.Name(c.Lib.GSS_KRB5_NT_PRINCIPAL_NAME)
	if err != nil {
		return nil, err
	}
	if gssapiName.String() != name {
		return nil, errors.Errorf("name: got %q, expected %q", gssapiName.String(), name)
	}

	return gssapiName, nil
}

func (c *Context) initSecContext(targetName *gssapi.Name) (
	*gssapi.Buffer, error) {

	_, _, token, _, _, err := c.Lib.InitSecContext(
		c.Lib.GSS_C_NO_CREDENTIAL,
		nil,
		targetName,
		c.Lib.GSS_C_NO_OID,
		0,
		0,
		c.Lib.GSS_C_NO_CHANNEL_BINDINGS,
		c.Lib.GSS_C_NO_BUFFER)
	return token, err
}

func PrepareSeverGSSAPIContext(serviceName string) (*Context, error) {

	var c = &Context{}

	err := c.loadLib(c.DebugLog, "")
	if err != nil {
		return nil, err
	}

	gssapiName, err := c.prepareServiceName(serviceName)
	if err != nil {
		return nil, err
	}

	defer gssapiName.Release()

	keytab := "/etc/krb5/krb5.keytab." + serviceName

	os.Setenv("KRB5_KTNAME", keytab)

	cred, actualMechs, _, err := c.Lib.AcquireCred(gssapiName,
		gssapi.GSS_C_INDEFINITE, c.Lib.GSS_C_NO_OID_SET, gssapi.GSS_C_ACCEPT)
	actualMechs.Release()
	if err != nil {
		return nil, err
	}

	c.credential = cred

	return c, err

}

func (c *Context) ParseToken(token string) (name string, err error) {

	tbytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}

	var inputToken *gssapi.Buffer
	if len(tbytes) != 0 {
		inputToken, err = c.Lib.MakeBufferBytes(tbytes)
		if err != nil {
			return "", err
		}
	}

	_, srcName, _, _, _, _, delegatedCredHandle, err :=
		c.Lib.AcceptSecContext(c.Lib.GSS_C_NO_CONTEXT,
			c.credential, inputToken, c.Lib.GSS_C_NO_CHANNEL_BINDINGS)

	if err != nil {
		return "", err
	}
	names := strings.Split(srcName.String(), "@")
	if len(names) > 0 {
		name = names[0]
	} else {
		return "", errors.New("Cannot extract username from GSSAPI token")
	}

	srcName.Release()
	delegatedCredHandle.Release()

	return name, nil
}

func GetGSSAPIToken(serviceName string) (data []byte, err error) {

	var c = &Context{}

	err = c.loadLib(c.DebugLog, "")
	if err != nil {
		return
	}

	gssapiName, err := c.prepareServiceName(serviceName)
	if err != nil {
		return
	}
	defer gssapiName.Release()
	token, err := c.initSecContext(gssapiName)
	if err != nil {
		e, ok := err.(*gssapi.Error)
		if ok && e.Major.ContinueNeeded() {
			err = errors.New("Unexpected GSS_S_CONTINUE_NEEDED")
			return
		}
		return
	}

	defer token.Release()

	if token.Length() == 0 {
		err = errors.New("Empty GSSAPI token")
		return
	}
	data = token.Bytes()
	return
}
