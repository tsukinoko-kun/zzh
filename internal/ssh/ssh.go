package ssh

import (
	"errors"
	"fmt"
	"os"

	"github.com/tsukinoko-kun/zzh/internal/config"
	"golang.org/x/crypto/ssh"
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

	if err := session.Shell(); err != nil {
		return errors.Join(errors.New("failed to start shell"), err)
	}

	return session.Wait()
}
