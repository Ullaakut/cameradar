package curl

import (
	"testing"
)

func TestVersionInfo(t *testing.T) {
	info := VersionInfo(VERSION_FIRST)
	expectedProtocols := []string{"dict", "file", "ftp", "ftps", "gopher", "http", "https", "imap", "imaps", "ldap", "ldaps", "pop3", "pop3s", "rtmp", "rtsp", "smtp", "smtps", "telnet", "tftp", "scp", "sftp", "smb", "smbs"}
	protocols := info.Protocols
	for _, protocol := range protocols {
		found := false
		for _, expectedProtocol := range expectedProtocols {
			if expectedProtocol == protocol {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("protocol should be in %v and is %v.", expectedProtocols, protocol)
		}
	}
}
