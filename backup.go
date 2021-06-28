package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"google.golang.org/api/drive/v3"
)

type Backup struct {
	FileName                   string
	ContainerName              string
	UserName, DBName           string
	CredentialsFile, TokenFile string
	MaxBackups                 int
}

func (b *Backup) getContainerID() (string, error) {

	if out, err := exec.Command("docker", "ps").Output(); err != nil {
		return "", err
	} else {
		for _, row := range strings.Split(string(out), "\n") {
			items := strings.Split(row, " ")
			if strings.Contains(strings.ToLower(items[len(items)-1]), strings.ToLower(b.ContainerName)) {
				return items[0], nil
			}
		}
	}
	return "", fmt.Errorf("contaier with name %s is not found", b.ContainerName)
}

func (b *Backup) dumpDB() error {
	if id, err := b.getContainerID(); err != nil {
		return err
	} else {
		return exec.Command("sh", "-c", "docker exec -i "+id+" pg_dump -U "+b.UserName+" -d "+b.DBName+" -Fc > files/"+b.FileName).Run()
	}
}

func (b *Backup) upload() error {

	g := &GoogleDrive{CredentialPath: b.CredentialsFile, TokenPath: b.TokenFile}
	if err := g.Init(); err != nil {
		return err
	}

	f, err := os.Open("files/" + b.FileName)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = g.Service.Files.Create(&drive.File{MimeType: "*/*", Name: b.FileName}).Media(f).Do(); err != nil {
		return err
	}

	r, err := g.Service.Files.List().PageSize(100).Fields("nextPageToken, files(id, name)").Q("name contains '" + b.FileName + "'").OrderBy("createdTime desc").Do()
	if err != nil {
		return err
	}

	for i, file := range r.Files {
		if i > b.MaxBackups {
			if err := g.Service.Files.Delete(file.Id).Do(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Backup) Create() {
	if err := b.dumpDB(); err != nil {
		AdminBot.SendError(nil, fmt.Sprintf("dumpDB %s", err.Error()))
	} else if err := b.upload(); err != nil {
		AdminBot.SendError(nil, fmt.Sprintf("upload backup %s", err.Error()))
	}
}
