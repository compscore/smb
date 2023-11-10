package smb

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"

	"github.com/hirochachacha/go-smb2"
)

type optionsStruct struct {
	// Domain is the domain to authenticate against
	Domain string `compscore:"domain"`

	// Share is the share to connect to
	Share string `compscore:"share"`

	// Check if the targeted file exists
	Exists bool `compscore:"exists"`

	// Check if contents of file matches a regex
	RegexMatch bool `compscore:"regex_match"`

	// Check if contents of file matches a string
	SubstringMatch bool `compscore:"match"`

	// Check if contents of file matches a string exactly
	Match bool `compscore:"match"`

	// Sha256 hash of the expected output
	Sha256 bool `compscore:"sha256"`

	// Md5 hash of the expected output
	Md5 bool `compscore:"md5"`

	// Sha1 hash of the expected output
	Sha1 bool `compscore:"sha1"`
}

func (o *optionsStruct) Unmarshal(options map[string]interface{}) {
	domainInterface, ok := options["domain"]
	if ok {
		domain, ok := domainInterface.(string)
		if ok {
			o.Domain = domain
		}
	}

	shareInterface, ok := options["share"]
	if ok {
		share, ok := shareInterface.(string)
		if ok {
			o.Share = share
		}
	}

	_, ok = options["exists"]
	if ok {
		o.Exists = true
	}

	_, ok = options["regex_match"]
	if ok {
		o.RegexMatch = true
	}

	_, ok = options["match"]
	if ok {
		o.Match = true
	}

	_, ok = options["sha256"]
	if ok {
		o.Sha256 = true
	}

	_, ok = options["md5"]
	if ok {
		o.Md5 = true
	}

	_, ok = options["sha1"]
	if ok {
		o.Sha1 = true
	}
}

func (o *optionsStruct) Check(expectedOutput string, file *smb2.File) error {
	if file == nil {
		return fmt.Errorf("no content found")
	}

	bodyBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	body := string(bodyBytes)

	if o.Exists {
		if body == "" {
			return fmt.Errorf("file is empty or does not exist")
		}
	}

	if o.RegexMatch {
		regexp, err := regexp.Compile(expectedOutput)
		if err != nil {
			return err
		}

		if !regexp.MatchString(body) {
			return fmt.Errorf("regex mismatch: expected \"%s\", got \"%s\"", expectedOutput, body)
		}
	}

	if o.SubstringMatch {
		if !strings.Contains(body, expectedOutput) {
			return fmt.Errorf("substring mismatch: expected \"%s\", got \"%s\"", expectedOutput, body)
		}
	}

	if o.Match {
		if body != expectedOutput {
			return fmt.Errorf("mismatch: expected \"%s\", got \"%s\"", expectedOutput, body)
		}
	}

	if o.Sha256 {
		hash := fmt.Sprintf("%x", sha256.Sum256(bodyBytes))
		if hash != expectedOutput {
			return fmt.Errorf("sha256 mismatch: expected \"%s\", got \"%s\"", expectedOutput, hash)
		}
	}

	if o.Md5 {
		hash := fmt.Sprintf("%x", sha256.Sum256(bodyBytes))
		if hash != expectedOutput {
			return fmt.Errorf("md5 mismatch: expected \"%s\", got \"%s\"", expectedOutput, hash)
		}
	}

	if o.Sha1 {
		hash := fmt.Sprintf("%x", sha256.Sum256(bodyBytes))
		if hash != expectedOutput {
			return fmt.Errorf("sha1 mismatch: expected \"%s\", got \"%s\"", expectedOutput, hash)
		}
	}

	return nil
}

func clean(input string) string {
	return strings.ReplaceAll(input, "\x00", "")
}

func Run(ctx context.Context, target string, command string, expectedOutput string, username string, password string, options map[string]interface{}) (bool, string) {
	if !strings.Contains(target, ":") {
		target = fmt.Sprintf("%s:445", target)
	}

	var o optionsStruct
	o.Unmarshal(options)

	conn, err := net.Dial("tcp", target)
	if err != nil {
		return false, clean(err.Error())
	}
	defer conn.Close()

	smbConn := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     username,
			Password: password,
			Domain:   o.Domain,
		},
	}

	s, err := smbConn.DialContext(ctx, conn)
	if err != nil {
		return false, clean(err.Error())
	}
	defer s.Logoff()

	fs, err := s.Mount(
		fmt.Sprintf(
			`\\%s\%s`,
			strings.Split(target, ":")[0],
			o.Share,
		),
	)
	if err != nil {
		return false, clean(err.Error())
	}
	defer fs.Umount()

	f, err := fs.Open(command)
	if err != nil {
		return false, clean(err.Error())
	}
	defer f.Close()

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return false, clean(err.Error())
	}

	err = o.Check(expectedOutput, f)
	if err != nil {
		return false, clean(err.Error())
	}

	return true, ""
}
