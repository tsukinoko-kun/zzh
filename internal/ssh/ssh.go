package ssh

import (
	"errors"
	"fmt"
	"os"

	"github.com/tsukinoko-kun/zzh/internal/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func Interactive(config *config.Config) error {
	sshConfig := &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", config.Host, config.Port), sshConfig)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to dial %s", config.Display()), err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return errors.Join(errors.New("failed to start ssh session"), err)
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	fd := int(os.Stdin.Fd())
	width, height, err := term.GetSize(fd)
	if err != nil {
		width = 80
		height = 80
	}
	if err := session.RequestPty("xterm", height, width, modes); err != nil {
		return fmt.Errorf("request for pseudo terminal failed: %w", err)
	}

	if err := session.Shell(); err != nil {
		return errors.Join(errors.New("failed to start shell"), err)
	}

	return session.Wait()
}
