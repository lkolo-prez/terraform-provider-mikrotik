package client

import (
	"fmt"
	"strings"
)

// SystemBackup represents a RouterOS system backup
type SystemBackup struct {
	Name     string
	Size     string
	Creation string
}

// SystemBackupSave represents parameters for backup save operation
type SystemBackupSave struct {
	Name        string
	Password    string
	DontEncrypt bool
}

// SaveSystemBackup creates a new system backup
func (c *Mikrotik) SaveSystemBackup(params *SystemBackupSave) error {
	mikrotikClient, err := c.getMikrotikClient()
	if err != nil {
		return err
	}

	cmd := []string{"/system/backup/save"}

	if params.Name != "" {
		cmd = append(cmd, "=name="+params.Name)
	}

	if params.Password != "" {
		cmd = append(cmd, "=password="+params.Password)
	}

	if params.DontEncrypt {
		cmd = append(cmd, "=dont-encrypt=yes")
	}

	_, err = mikrotikClient.RunArgs(cmd)
	if err != nil {
		return fmt.Errorf("failed to save backup: %w", err)
	}

	return nil
}

// DeleteSystemBackup deletes a backup file
func (c *Mikrotik) DeleteSystemBackup(name string) error {
	// Ensure .backup extension
	if !strings.HasSuffix(name, ".backup") {
		name = name + ".backup"
	}

	return c.DeleteFile(name)
}

// File operations

// File represents a RouterOS file
type File struct {
	Id           string `mikrotik:".id"`
	Name         string `mikrotik:"name"`
	Type         string `mikrotik:"type"`
	Size         string `mikrotik:"size"`
	CreationTime string `mikrotik:"creation-time"`
}

// ListFiles returns all files
func (c *Mikrotik) ListFiles() ([]*File, error) {
	mikrotikClient, err := c.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := []string{"/file/print"}
	reply, err := mikrotikClient.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	var files []*File
	for _, sentence := range reply.Re {
		file := &File{
			Id:           sentence.Map[".id"],
			Name:         sentence.Map["name"],
			Type:         sentence.Map["type"],
			Size:         sentence.Map["size"],
			CreationTime: sentence.Map["creation-time"],
		}
		files = append(files, file)
	}

	return files, nil
}

// FindFile finds a specific file by exact name
func (c *Mikrotik) FindFile(name string) (*File, error) {
	files, err := c.ListFiles()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.Name == name {
			return file, nil
		}
	}

	return nil, NewNotFound("file not found")
}

// FindSystemBackup finds a specific backup by name
func (c *Mikrotik) FindSystemBackup(name string) (*SystemBackup, error) {
	// Ensure .backup extension
	if !strings.HasSuffix(name, ".backup") {
		name = name + ".backup"
	}

	file, err := c.FindFile(name)
	if err != nil {
		return nil, err
	}

	backup := &SystemBackup{
		Name:     file.Name,
		Size:     file.Size,
		Creation: file.CreationTime,
	}

	return backup, nil
}

// DeleteFile deletes a file
func (c *Mikrotik) DeleteFile(name string) error {
	mikrotikClient, err := c.getMikrotikClient()
	if err != nil {
		return err
	}

	file, err := c.FindFile(name)
	if err != nil {
		return err
	}

	cmd := []string{"/file/remove", "=.id=" + file.Id}
	_, err = mikrotikClient.RunArgs(cmd)
	return err
}
