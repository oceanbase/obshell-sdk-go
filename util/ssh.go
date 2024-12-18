/*
 * Copyright (c) 2024 OceanBase.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"math"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/cavaliergopher/cpio"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"github.com/oceanbase/obshell-sdk-go/internal/util"
	"github.com/oceanbase/obshell-sdk-go/log"
	"github.com/oceanbase/obshell-sdk-go/model"
	"github.com/oceanbase/obshell-sdk-go/services"
)

const (
	DEFALUT_SSH_PORT     = 22
	DEFALUT_OBSHELL_PORT = 2886
	DEFALUT_SSH_PATH     = ".ssh"
)

var (
	obshellPidFiles = []string{"daemon.pid", "obshell.pid"}
	pidFiles        = append(obshellPidFiles, "observer.pid")

	// excludeFile is a list of files that not related to private keys
	excludeFile = []string{"authorized_keys", "config", "id_rsa.pub", "known_hosts"}

	// localAddresses is a list of local IP addresses
	localAddresses []*net.IPNet

	// Unlike the Python SDK, the Go SDK does not use rsync by default, but instead uses sftp for chunked transfer.
	// This is because the rsync implementation in the GO SDK is not fully developed and can only be used with default key configuration for passwordless login.
	// At the same time, the performance of the sftp chunked transfer in the GO SDK is comparable to rsync.
	UseRsync          = false
	remoteRsyncStatus = make(map[string]bool, 0)

	// The size of each sftp chunked transfer, default is 64M.
	// Since the maximum size of a single file in the current OB, observer, is around 450M, with 64M, it can be divided into 7-8 chunks.
	// This would require 7-8 concurrent connections, and since the default MaxSessions in the sshd configuration is 10, this value is appropriate.
	// If you want to improve the performance of sftp chunked transfer, you can reduce this value to increase the number of concurrent connections.
	// However, you will need to correspondingly increase the MaxSessions configuration on the target machine.
	CHUNK_SIZE = 1024 * 1024 * 64
	// The maximum number of parallel SFTP transfers to avoid exceeding the MaxSessions limit
	PARALLEL_SFTP_MAX = 8

	// Since using in-memory backup for small files results in higher batch transfer efficiency, only files larger than SCP_THRESHOLD will be transferred using SCP.
	// You can disable SCP by setting SCP_THRESHOLD=0.
	SCP_THRESHOLD int64 = 1024 * 1024 * 1
)

func initLocalAddresses() {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Warnf("Failed to get network interfaces: %v", err)
		return
	}

	for _, iface := range interfaces {
		// Skip inactive interfaces
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		// Get the list of addresses of the interface
		addrs, err := iface.Addrs()
		if err != nil {
			log.Warnf("Failed to get addresses of interface %s: %v", iface.Name, err)
			return
		}

		// Traverse each address of the interface
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			localAddresses = append(localAddresses, ipNet)
		}
	}
}

func loadDefaultPrivateKeys() ([]ssh.Signer, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	defaultDir := filepath.Join(home, DEFALUT_SSH_PATH)

	var signers []ssh.Signer
	files, err := os.ReadDir(defaultDir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if util.ContainsString(excludeFile, file.Name()) {
			continue
		}
		keyPath := filepath.Join(defaultDir, file.Name())
		// If an error occurs while loading the private key, ignore it, because it may have no permission
		signer, _ := loadPrivateKey(keyPath)
		if signer != nil {
			signers = append(signers, signer)
		}
	}
	return signers, nil
}

func loadPrivateKey(keyPath string) (ssh.Signer, error) {
	file, err := os.Open(keyPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	keyData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	privateKey, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func isLocalAddress(address string) (bool, error) {
	parsedIP := net.ParseIP(address)
	if parsedIP == nil {
		return false, fmt.Errorf("invalid IP address: %s", address)
	}

	if parsedIP.IsLoopback() {
		return true, nil
	}

	for _, ipNet := range localAddresses {
		if ipNet.IP.Equal(parsedIP) {
			return true, nil
		}
	}

	return false, nil
}

type NodeConfig struct {
	SSHConfig
	obshellPort int
	workDir     string
	ip          string
}

type SSHConfig struct {
	*ssh.ClientConfig
	sshPort int
}

type NodeClient struct {
	*ssh.Client
	NodeConfig
	isLocal bool
	address string
}

type SshRetun struct {
	Stdout string
	Stderr string
	Code   int
}

func newSshReturn(stdout, stderr string, code int) SshRetun {
	return SshRetun{
		Stdout: stdout,
		Stderr: stderr,
		Code:   code,
	}
}

func (client *NodeClient) ExecuteCommand(cmd string) SshRetun {
	if client.isLocal {
		return executeLocal("bash", "-c", cmd)
	}
	return executeRemote(client.Client, cmd)
}

func (client *NodeClient) wirteFileByScp(f io.Reader, filePath string, mode fs.FileMode) error {
	log.Infof("Write file: %s, use scp", filePath)
	cmd := fmt.Sprintf("mkdir -p %s", filepath.Dir(filePath))
	if ret := client.ExecuteCommand(cmd); ret.Code != 0 {
		return fmt.Errorf("failed to create directory: %s", ret.Stderr)
	}

	if client.isLocal {
		context, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		return writeFileLocal(context, filePath, mode)
	}

	sshConfig := client.NodeConfig.SSHConfig.ClientConfig
	scpClient := scp.NewClient(client.address, sshConfig)
	if err := scpClient.Connect(); err != nil {
		return err
	}
	defer scpClient.Close()

	return scpClient.CopyFile(context.Background(), f, filePath, fmt.Sprintf("0%o", mode.Perm()))
}

func (client *NodeClient) WriteFile(context []byte, filePath string, mode fs.FileMode) (err error) {
	log.Infof("Write file: %s, use rsync: %v", filePath, client.useRsync())
	if client.isLocal {
		return writeFileLocal(context, filePath, mode)
	}

	cmd := fmt.Sprintf("mkdir -p %s", filepath.Dir(filePath))
	if ret := client.ExecuteCommand(cmd); ret.Code != 0 {
		return fmt.Errorf("failed to create directory: %s", ret.Stderr)
	}

	if client.useRsync() {
		err = client.writeFileByRsync(context, filePath)
	} else {
		err = client.writeFileBySFTP(context, filePath)
	}
	if err != nil {
		return err
	}

	cmd = fmt.Sprintf("chmod %o %s", mode.Perm(), filePath)
	if ret := client.ExecuteCommand(cmd); ret.Code != 0 {
		return fmt.Errorf("failed to update file mode: %s", ret.Stderr)
	}
	return nil
}

func (client *NodeClient) useRsync() bool {
	if !UseRsync {
		return false
	}

	if use, ok := remoteRsyncStatus[client.ip]; ok {
		return use
	} else {
		ret := client.ExecuteCommand("rsync -h")
		use = ret.Code == 0
		if !use {
			log.Info("rsync is not installed on %s: %s\n", client.ip, ret.Stderr)
		}
		remoteRsyncStatus[client.ip] = use
		return use
	}
}

func (client *NodeClient) writeFileBySFTP(context []byte, filePath string) error {
	size := len(context)
	if size < CHUNK_SIZE {
		return client.writeChunk(context, filePath)
	}

	defer func() {
		client.ExecuteCommand(fmt.Sprintf("rm -f %s.*", filePath))
	}()

	var errs []error
	wg := sync.WaitGroup{}
	for i := 0; i < size; i += CHUNK_SIZE {
		if i/CHUNK_SIZE/PARALLEL_SFTP_MAX > 0 && i/CHUNK_SIZE%PARALLEL_SFTP_MAX == 0 {
			// Wait for the completion of the previous batch of sftp operations
			wg.Wait()
		}

		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			end := start + CHUNK_SIZE
			if end > size {
				end = size
			}
			chunkFile := fmt.Sprintf("%s.%d", filePath, start)
			if newClient, err := NewNodeClient(client.NodeConfig); err != nil {
				errs = append(errs, err)
				log.Warn("Failed to create new client: %s", err)
			} else if err := newClient.writeChunk(context[start:end], chunkFile); err != nil {
				errs = append(errs, err)
				log.Warn("Failed to write chunk: %s", err)
			}
		}(i)
	}
	wg.Wait()

	if len(errs) > 0 {
		return errors.Errorf("failed to write file by sftp: %v", errs)
	}

	for i := 0; i < size; i += CHUNK_SIZE {
		chunkFile := fmt.Sprintf("%s.%d", filePath, i)
		cmd := fmt.Sprintf("cat %s >> %s; rm -f %s", chunkFile, filePath, chunkFile)
		if ret := client.ExecuteCommand(cmd); ret.Code != 0 {
			return fmt.Errorf("failed to merge chunk: %s", ret.Stderr)
		}
	}
	return nil
}

func (client *NodeClient) writeChunk(context []byte, filePath string) error {
	log.Debug("write chunk: ", filePath)
	sftpClient, err := sftp.NewClient(client.Client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	dstFile, err := sftpClient.Create(filePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	srcFile := io.NopCloser(bytes.NewReader(context))
	_, err = io.Copy(dstFile, srcFile)
	return err
}

func (client *NodeClient) writeFileByRsync(context []byte, filePath string) error {
	tempFile, err := os.CreateTemp("", "obshell=go-sdk-temp-")
	if err != nil {
		return err
	}
	defer tempFile.Close()

	if _, err = tempFile.Write(context); err != nil {
		return err
	}

	identityOption := "-o StrictHostKeyChecking=no "
	if client.sshPort != DEFALUT_SSH_PORT {
		identityOption += fmt.Sprintf("-p %d ", client.sshPort)
	}

	rsyncTarget := fmt.Sprintf("%s@%s:%s", client.SSHConfig.User, client.ip, filePath)
	cmd := fmt.Sprintf("yes | rsync -a -W -L -e 'ssh %s' %s %s", identityOption, tempFile.Name(), rsyncTarget)
	ret := executeLocal("bash", "-c", cmd)

	if ret.Code != 0 {
		return fmt.Errorf("failed to write file by rsync: %s", ret.Stderr)
	}
	return nil
}

func writeFileLocal(context []byte, filePath string, mode fs.FileMode) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.Write(context); err != nil {
		return err
	}

	if err = file.Chmod(mode); err != nil {
		return err
	}
	return nil
}

func executeLocal(name string, arg ...string) SshRetun {
	command := exec.Command(name, arg...)

	// Create buffers to capture output
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			stderr.WriteString(err.Error())
			exitCode = 127
		}
	}

	return newSshReturn(stdout.String(), stderr.String(), exitCode)
}

func executeRemote(client *ssh.Client, cmd string) SshRetun {
	session, err := client.NewSession()
	if err != nil {
		return newSshReturn("", err.Error(), 127)
	}
	defer session.Close()
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(cmd)
	exitCode := 0
	if err != nil {
		// Get exit code
		if exitErr, ok := err.(*ssh.ExitError); ok {
			exitCode = exitErr.ExitStatus()
		} else {
			exitCode = 127
		}
	}

	return newSshReturn(stdout.String(), stderr.String(), exitCode)
}

var defaultSSHConfig SSHConfig

func init() {
	defaultSSHConfig = SSHConfig{
		ClientConfig: &ssh.ClientConfig{
			User:            os.Getenv("USER"),
			Auth:            []ssh.AuthMethod{ssh.PublicKeys()},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
		sshPort: DEFALUT_SSH_PORT,
	}

	if signers, err := loadDefaultPrivateKeys(); err == nil {
		defaultSSHConfig.ClientConfig.Auth = []ssh.AuthMethod{ssh.PublicKeys(signers...)}
	} else {
		defaultSSHConfig.ClientConfig.Auth = []ssh.AuthMethod{ssh.Password("")}
	}

	initLocalAddresses()

	if ret := executeLocal("rsync", "-h"); ret.Code != 0 {
		log.Debug("rsync may not be installed: %s", ret.Stderr)
		UseRsync = false
	}
}

func NewSSHConfig(config ssh.ClientConfig, sshPort ...int) SSHConfig {
	sshConfig := SSHConfig{
		ClientConfig: &config,
	}
	if len(sshPort) > 0 {
		sshConfig.sshPort = sshPort[0]
	} else {
		sshConfig.sshPort = DEFALUT_SSH_PORT
	}
	return sshConfig
}

func NewNodeConfig(ip string, workDir string, obshellPort ...int) NodeConfig {
	return NewNodeConfigWithSShConfig(ip, workDir, defaultSSHConfig, obshellPort...)
}

func NewNodeConfigWithSShConfig(ip string, workDir string, config SSHConfig, obshellPort ...int) NodeConfig {
	obshellPort = append(obshellPort, DEFALUT_OBSHELL_PORT)
	return NodeConfig{
		ip:          ip,
		SSHConfig:   config,
		obshellPort: obshellPort[0],
		workDir:     workDir,
	}
}

func NewNodeClient(config NodeConfig) (*NodeClient, error) {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.ip, config.sshPort), config.ClientConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "dial to %s:%d failed", config.ip, config.sshPort)
	}

	isLocal := false
	if config.User == defaultSSHConfig.User {
		if isLocal, err = isLocalAddress(config.ip); err != nil {
			return nil, err
		}
	}

	var address string
	if strings.Contains(config.ip, ":") {
		address = fmt.Sprintf("[%s]:%d", config.ip, config.sshPort)
	} else {
		address = fmt.Sprintf("%s:%d", config.ip, config.sshPort)
	}

	return &NodeClient{
		Client:     client,
		NodeConfig: config,
		isLocal:    isLocal,
		address:    address,
	}, nil
}

func cleanNode(client *NodeClient) error {
	for _, pidFile := range pidFiles {
		if err := stopProcess(client, pidFile); err != nil {
			log.Warn(err)
		}
	}

	ret := client.ExecuteCommand(fmt.Sprintf("rm -fr %s", client.workDir))
	if ret.Code != 0 {
		return fmt.Errorf("failed to clean %s work dir %s: %s", client.ip, client.workDir, ret.Stderr)
	}
	return nil
}

func stopObshell(client *NodeClient) {
	for _, pidFile := range obshellPidFiles {
		if err := stopProcess(client, pidFile); err != nil {
			log.Warn(err)
		}
	}
}

func stopProcess(client *NodeClient, pidFile string) error {
	path := filepath.Join(client.workDir, "run", pidFile)
	ret := client.ExecuteCommand(fmt.Sprintf("[ -f %s ]", path))
	if ret.Code != 0 {
		return nil
	}

	ret = client.ExecuteCommand(fmt.Sprintf("kill -9 `cat %s`", path))
	if ret.Code != 0 {
		return errors.Errorf("Failed to kill %s(%s): %s", client.ip, pidFile, ret.Stderr)
	}
	return nil
}

func createClientMap(configs ...NodeConfig) (map[NodeConfig]*NodeClient, error) {
	clientMap := make(map[NodeConfig]*NodeClient)
	for _, config := range configs {
		log.Infof("Connecting to %s\n", config.ip)
		client, err := NewNodeClient(config)
		if err != nil {
			return nil, err
		}
		clientMap[config] = client
	}
	return clientMap, nil
}

// 将Obshell安装到指定的服务器上, 自会将obshell RPM包
func InstallObshell(rpmPackagePath string, configs ...NodeConfig) error {
	clientMap := make(map[NodeConfig]*NodeClient)
	defer func() {
		for _, client := range clientMap {
			client.Close()
		}
	}()
	// Get the connection to the servers.
	for _, config := range configs {
		log.Info("Connecting to %s\n", config.ip)
		client, err := NewNodeClient(config)
		if err != nil {
			log.Error("Failed to connect to ", config.ip)
			return err
		}
		clientMap[config] = client
	}

	return installRpmPackages([]string{rpmPackagePath}, clientMap)
}

func InitNodes(rpmPackagePaths []string, forceClean bool, configs ...NodeConfig) error {
	var clientMap map[NodeConfig]*NodeClient
	defer func() {
		for _, client := range clientMap {
			client.Close()
		}
	}()
	clientMap, err := createClientMap(configs...)
	if err != nil {
		return err
	}

	// Check the work directory. If the work directory is exist and not empty, clean it.
	for config, client := range clientMap {
		if forceClean {
			log.Debugf("Force clean the work directory: %s\n", config.workDir)
			if err := cleanNode(client); err != nil {
				return err
			}
		} else {
			isEmpty, err := checkRemoteDirEmpty(client.Client, config.workDir)
			if err != nil {
				return err
			}
			if !isEmpty {
				return fmt.Errorf("%s:%s is not empty, please clean it first", config.ip, config.workDir)
			}
		}
		client.ExecuteCommand(fmt.Sprintf("mkdir -p %s", config.workDir))
	}

	return installRpmPackages(rpmPackagePaths, clientMap)
}

func installRpmPackages(rpmPackagePaths []string, clientMap map[NodeConfig]*NodeClient) error {
	addressMap := make(map[string][]NodeConfig)
	for config, client := range clientMap {
		configs, ok := addressMap[client.address]
		if !ok {
			configs = make([]NodeConfig, 0)
		}
		addressMap[client.address] = append(configs, config)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(clientMap))
	pacakgeInstalledFiles := make(map[string][]string, 0)
	pacakgeLinkedFiles := make(map[string]map[string]string, 0)
	for _, configs := range addressMap {
		config := configs[0]
		client := clientMap[config]

		wg.Add(1)
		go func(config NodeConfig, client *NodeClient) {
			defer wg.Done()
			for _, rmpPkg := range rpmPackagePaths {
				if installedFiles, linkMap, err := installRpmPackage(rmpPkg, *client); err != nil {
					errChan <- err
				} else {
					pacakgeInstalledFiles[rmpPkg] = installedFiles
					pacakgeLinkedFiles[rmpPkg] = linkMap
				}
			}
		}(config, client)
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		return err
	}

	// Copy installed files from one node to other nodes on the same machine to improve installation efficiency.
	for _, configs := range addressMap {
		for _, config := range configs[1:] {
			client := clientMap[config]

			for _, installedFiles := range pacakgeInstalledFiles {
				for _, file := range installedFiles {
					srcPath := getDestPath(configs[0], file)
					dstPath := getDestPath(config, file)
					cmd := fmt.Sprintf("mkdir -p %s; cp %s %s", filepath.Dir(dstPath), srcPath, dstPath)
					if ret := client.ExecuteCommand(cmd); ret.Code != 0 {
						return fmt.Errorf("failed to copy files to %s:%s: %s", config.ip, config.workDir, ret.Stderr)
					}
				}
			}

			for _, linkMap := range pacakgeLinkedFiles {
				for target, source := range linkMap {
					srcPath := getDestPath(configs[0], source)
					dstPath := getDestPath(config, target)
					cmd := fmt.Sprintf("mkdir -p %s; ln -sf %s %s", filepath.Dir(dstPath), srcPath, dstPath)
					if ret := client.ExecuteCommand(cmd); ret.Code != 0 {
						return fmt.Errorf("failed to create link to %s:%s: %s", config.ip, target, ret.Stderr)
					}
				}
			}
		}
	}
	return nil
}

func installRpmPackage(rmpPkg string, client NodeClient) ([]string, map[string]string, error) {
	log.Info("Load rpm package:", rmpPkg)
	pkg := newRpmPackage(rmpPkg)
	if err := pkg.open(); err != nil {
		return nil, nil, errors.Wrap(err, "open rpm package failed")
	}
	defer pkg.close()

	installedFiles := make([]string, 0)
	linkMap := make(map[string]string)
	fileContents := make([]*fileContent, 0)
	for {
		header, err := pkg.next()
		if err != nil {
			return nil, nil, errors.Wrap(err, "cpio read failed")
		} else if header == nil {
			break
		}

		if header.Linkname != "" {
			linkMap[header.Name] = header.Linkname
			continue
		} else if header.Mode.IsDir() {
			continue
		} else if header.Size < SCP_THRESHOLD {
			content := make([]byte, header.Size)
			if _, err := io.ReadFull(pkg.cpioReader, content); err != nil {
				return nil, nil, errors.Wrap(err, "read full failed")
			}
			fileContents = append(fileContents, newFileContent(header, content))
			installedFiles = append(installedFiles, header.Name)
			continue
		}

		installedFiles = append(installedFiles, header.Name)
		destPath := getDestPath(client.NodeConfig, header.Name)
		if err = client.wirteFileByScp(pkg.cpioReader, destPath, fs.FileMode(header.Mode)); err != nil {
			return nil, nil, err
		}
	}

	paralleWriteFiles(&client, fileContents)

	// Create link files
	for target, source := range linkMap {
		log.Infof("link to %s", target)
		config := client.NodeConfig
		srcPath := getDestPath(config, source)
		dstPath := getDestPath(config, target)
		cmd := fmt.Sprintf("mkdir -p %s; ln -sf %s %s", filepath.Dir(dstPath), srcPath, dstPath)

		if ret := client.ExecuteCommand(cmd); ret.Code != 0 {
			return nil, nil, fmt.Errorf("failed to create link: %s", ret.Stderr)
		}
	}
	return installedFiles, linkMap, nil
}

type fileContent struct {
	header  *cpio.Header
	content []byte
}

func newFileContent(header *cpio.Header, content []byte) *fileContent {
	return &fileContent{
		header:  header,
		content: content,
	}
}

func paralleWriteFiles(client *NodeClient, fileContents []*fileContent) error {
	count := len(fileContents)
	cpuCount := runtime.NumCPU()
	size := count/cpuCount + 1

	errChan := make(chan error, cpuCount)
	var wg sync.WaitGroup
	for i := 0; i < count; i += size {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			end := start + size
			if end > count {
				end = count
			}
			for _, f := range fileContents[start:end] {
				destPath := getDestPath(client.NodeConfig, f.header.Name)
				if err := client.WriteFile(f.content, destPath, fs.FileMode(f.header.Mode)); err != nil {
					errChan <- err
				}
			}
		}(i)
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		return err
	}
	return nil
}

func getDestPath(config NodeConfig, name string) string {
	// Remove the prefix "./home/admin/oceanbase" if it exists
	fileName, found := strings.CutPrefix(name, "./home/admin/oceanbase")
	if !found {
		// Remove the prefix "./usr" if it exists
		fileName, _ = strings.CutPrefix(name, "./usr")
	}
	return filepath.Join(config.workDir, fileName)
}

func StartObshell(configs ...NodeConfig) error {
	clientMap := make(map[*ssh.Client]NodeConfig)
	defer func() {
		for client := range clientMap {
			client.Close()
		}
	}()

	for _, config := range configs {
		log.Infof("Connecting to %s\n", config.ip)
		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.ip, config.sshPort), config.ClientConfig)
		if err != nil {
			return errors.Wrapf(err, "dial to %s:%d failed", config.ip, config.sshPort)
		}
		clientMap[client] = config
	}

	for client, config := range clientMap {
		if err := startRemoteObshell(client, config.workDir, config.ip, config.obshellPort, nil); err != nil {
			return err
		}
	}

	return nil
}

type checkItem struct {
	message string
	suggest string
}

func CheckNodes(configs ...NodeConfig) (errors, warns []checkItem){
	var clientMap map[NodeConfig]*NodeClient
	defer func() {
		for _, client := range clientMap {
			client.Close()
		}
	}()
	clientMap, err := createClientMap(configs...)
	if err != nil {
		errors = append(errors, checkItem{err.Error(), ""})
		return errors, warns
	}
	serverNum := len(clientMap)
	for _, client := range clientMap {
		if checkNodesFirewalld(client) != nil {
			errors = append(errors, checkItem{"the firewalld service is up.", "please stop the firewalld service."})
			return errors, warns
		}
		if checkNodesSelinux(client) != nil {
			errors = append(errors, checkItem{"the selinux is not Disabled.", "please disabled the selinux."})
			return errors, warns
		}
		if  checkNodesClock(client) != nil {
			errors = append(errors, checkItem{"clock not sync.", "please sync clock."})
			return errors, warns
		}
		if message, err := checkNodesKernelParams(client); err != nil {
			for _, str := range message {
				errors = append(errors, checkItem{str, ""})
			}
			return errors, warns
		}
		if message, err := checkNodesUlimitParams(client,serverNum); err != nil {
			for _, str := range message {
				errors = append(errors, checkItem{str, ""})
			}
			return errors, warns
		}
	}

	return nil, nil
}

func checkNodesFirewalld(client *NodeClient) error{
	osType, err := getRemoteLinuxType(client)
	if err != nil {
		return err
	} else if strings.Contains(osType, "ubuntu") || strings.Contains(osType, "debian"){
		// Assume Ubuntu or Debian uses UFW
		ret := client.ExecuteCommand("ufw status")
		if strings.Contains(ret.Stdout, "Status: active"){
			return fmt.Errorf("the firewalld service is up")
		}
	} else if strings.Contains(osType, "fedora") || strings.Contains(osType, "centos") || strings.Contains(osType, "redhat"){
		// Assume Fedora, CentOS, or RedHat uses firewalld
		ret := client.ExecuteCommand("systemctl status firewalld")
		if strings.Contains(ret.Stdout, "Active: active"){
			return fmt.Errorf("the firewalld service is up")
		}
	} else {
		// For other Linux distributions, default to checking iptables
		ret := client.ExecuteCommand("iptables -L -n")
		if strings.Contains(ret.Stdout, "Chain INPUT") && strings.Contains(ret.Stdout, "Chain FORWARD") && strings.Contains(ret.Stdout, "Chain OUTPUT"){
			return fmt.Errorf("the firewalld service is up")
		}
	}
	return nil
}

func getRemoteLinuxType(client *NodeClient) (string, error){
	var ret SshRetun
	var systemType string
	if ret = client.ExecuteCommand("cat /etc/os-release"); ret.Code != 0 {
		err := fmt.Errorf("failed get os type: %s", ret.Stderr)
		return "", err
	}
	for _, line := range strings.Split(ret.Stdout, "\n") {
		if strings.HasPrefix(line, "ID=") {
			systemType = strings.ReplaceAll(strings.Split(line, "=")[1], "\"", "")
			return systemType, nil
		}
	}
	return "", nil
}

func checkNodesSelinux(client *NodeClient) error{
	if ret := client.ExecuteCommand("/usr/sbin/getenforce"); ret.Code != 0 {
		return fmt.Errorf("failed get selinux status: %s", ret.Stderr)
	} else if strings.Contains(ret.Stdout, "Enforcing"){
		return fmt.Errorf("the selinux is not Disabled")
	}
	return nil
}

func checkNodesClock(client *NodeClient) error{
	addrs, err := net.InterfaceAddrs()
    if err != nil {
		log.Infof("Get host addr error: %s\n", err.Error())
        os.Exit(1)
    }
	
    for _, address := range addrs {
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
				ret := client.ExecuteCommand(fmt.Sprintf("sudo /usr/sbin/clockdiff -o %s", ipnet.IP.String()))
				clockDiffStr := strings.Fields(ret.Stdout)[1]
				clockDiffNum, err := strconv.Atoi(clockDiffStr)
				if err != nil {
					return err
				}
				clockDiffNumAbs := math.Abs(float64(clockDiffNum))
				if clockDiffNumAbs/1000 > 2 {
					return fmt.Errorf("the time difference between two servers exceeds 2")
				}
			}
		}
	}
	return nil
}

func checkNodesKernelParams(client *NodeClient) ([]string, error){
	INF := int(math.Inf(1))
	var ret SshRetun

	kernelCheckItems := []map[string]any{
        {"check_item": "vm.max_map_count", "need": []int{327600, 1310720}, "recommend": 655360},
        {"check_item": "vm.min_free_kbytes", "need": []int{32768, 2097152}, "recommend": 2097152},
        {"check_item": "vm.overcommit_memory", "need": 0, "recommend": 0},
        {"check_item": "fs.file-max", "need": []int{6573688, INF}, "recommend": 6573688},
	}
	if ret = client.ExecuteCommand("sudo /usr/sbin/sysctl -a"); ret.Code != 0 {
		err := fmt.Errorf("sysctl command error: %s", ret.Stderr)
		return nil, err
	}
	kernelParams := make(map[string][]string)
	
	kernelParamSrc := strings.Split(ret.Stdout, "\n")
	for _, kernel := range kernelParamSrc{
		if kernel == "" {
			continue
		}
		kernelList := strings.Split(kernel, "=")
		pattern := regexp.MustCompile(`[-+]?\d+`)
		kernelParams[strings.TrimSpace(kernelList[0])] = pattern.FindAllString(kernelList[1], -1)
	}
	retSuggests := make([]string, 0)
	for _, kernelParam := range kernelCheckItems {
		checkItem := kernelParam["check_item"]
		var checkItemStr string
		var ok bool
		if checkItemStr, ok = checkItem.(string); ok {
			if _, ok := kernelParams[checkItemStr]; !ok {
				continue
			}
		} else {
			return nil, fmt.Errorf("any type conversion to string type failed")
		}
		values := kernelParams[checkItemStr]
		needs := kernelParam["need"]
		recommends := kernelParam["recommend"]
		for i := 0; i < len(values); i++ {
			// This case is not handling the value of 'default'. Additional handling is required for 'default' in the future.
			itemValue, err := strconv.Atoi(values[i])
			if err != nil {
				return nil, fmt.Errorf("string type conversion to int type failed")
			}
			var need int
			var recommend int
			if isSlice(needs) {
				needList := needs.([]int)
				need = needList[i]
			}
			if isSlice(recommends) {
				recommendList := recommends.([]int)
				recommend = recommendList[i]
			} else {
				recommend = recommends.(int)
			}
			if itemValue != need && itemValue != recommend{
				retSuggests = append(retSuggests, fmt.Sprintf("%v -> %v current value: %v, recommend: %v", client.ip, checkItem, itemValue, recommend))
				return retSuggests, fmt.Errorf("please use the recommended parameter values")
			}
		}
	}
	return retSuggests,nil
}

func isSlice(v any) bool {
	rv := reflect.ValueOf(v)
	return rv.Kind() == reflect.Slice
}

func checkNodesUlimitParams(client *NodeClient, serverNum int) ([]string, error){
	INF := int(math.Inf(1))

	ulimitsMin := map[string]map[string]func(int) string{
        "open files":
			{"need":func(x int) string { return strconv.Itoa(20000 * x) }, 
			"recd": func(x int) string { return strconv.Itoa(655350) },
			"name": func(x int) string { return "nofile"} },
		"max user processes":
			{"need":func(x int) string { return strconv.Itoa(120000) }, 
			"recd": func(x int) string { return strconv.Itoa(655350) },
			"name": func(x int) string { return"nproc"}},
		"core file size":
			{"need":func(x int) string { return strconv.Itoa(0) }, 
			"recd": func(x int) string { return strconv.Itoa(INF) },
			"below_need_error": func(x int) string { return "false"},
            "below_recd_error_strict": func(x int) string { return "false"},
			"name": func(x int) string { return "core"}},
		"stack size":
			{"need":func(x int) string { return strconv.Itoa(1024) }, 
			"recd": func(x int) string { return strconv.Itoa(INF) },
			"below_recd_error_strict": func(x int) string { return "false"},
			"name": func(x int) string { return "stack"}},
	}

	ret := client.ExecuteCommand("ulimit -a")
	retList := strings.Split(ret.Stdout, "\n")

	ulimits :=  make(map[string]string)
	for _, value := range retList {
		if len(strings.Split(value, `(`)) > 1 {
			ulimits[strings.TrimSpace(strings.Split(value, `(`)[0])] = strings.Split(value, ")")[1]
		}
	}
	retSuggests := make([]string, 0)
	for key := range ulimitsMin{
		value := ulimits[key]
		if  strings.TrimSpace(value) == "unlimited"{
			continue
		}
		valueInt, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			fmt.Printf("转换出错: %v\n", err)
			return retSuggests, err
		}
		need := ulimitsMin[key]["need"](serverNum)
		needInt, err := strconv.Atoi(need)
		if err != nil {
			fmt.Printf("转换出错: %v\n", err)
			return retSuggests, err
		}
		if valueInt < needInt {
			retSuggests = append(retSuggests, fmt.Sprintf("%v -> %v{%v} current value: %v, recommend: %v", client.ip, key, ulimitsMin[key]["name"](serverNum), valueInt, needInt))
			return retSuggests, fmt.Errorf("please use the recommended parameter values")
		}
	}
	return retSuggests, nil
}

func checkRemoteDirEmpty(sshClient *ssh.Client, filePath string) (bool, error) {
	log.Infof("Check remote directory: %s\n", filePath)
	session, err := sshClient.NewSession()
	if err != nil {
		return false, err
	}
	defer session.Close()

	cmd := fmt.Sprintf("if [ -d %s ]; then ls -A %s; fi", filePath, filePath)
	output, err := session.Output(cmd)
	if err != nil {
		return false, errors.Wrap(err, "session run failed")
	}
	log.Debug("Output is ", string(output))
	return len(output) == 0, nil
}

const FLAG_ROOT_PWD = "password"

func startRemoteObshell(sshClient *ssh.Client, workDir, ip string, obshellPort int, password *string) error {
	log.Infof("Start obshell on %s\n", ip)
	cmd := fmt.Sprintf("cd %s; ./bin/obshell admin start --ip %s --port %d", workDir, ip, obshellPort)
	if password != nil {
		if ret := executeRemote(sshClient, fmt.Sprintf("%s/bin/obshell admin start -h | grep %s", workDir, FLAG_ROOT_PWD)); ret.Code != 0 {
			cmd = fmt.Sprintf("export OB_ROOT_PASSWORD='%s'; %s", strings.ReplaceAll(*password, "'", "'\"'\"'"), cmd)
		} else {
			cmd = fmt.Sprintf("%s --%s='%s'", cmd, FLAG_ROOT_PWD, strings.ReplaceAll(*password, "'", "'\"'\"'"))
		}
	}

	ret := executeRemote(sshClient, cmd)
	if ret.Code != 0 {
		return fmt.Errorf("failed to start obshell: %s", ret.Stderr)
	}
	return nil
}

func Takeover(password string, configs ...NodeConfig) error {
	var clientMap map[NodeConfig]*NodeClient
	defer func() {
		for _, client := range clientMap {
			client.Close()
		}
	}()
	clientMap, err := createClientMap(configs...)
	if err != nil {
		return err
	}

	// stop obshell
	for _, client := range clientMap {
		stopObshell(client)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(clientMap))
	for config, client := range clientMap {
		wg.Add(1)
		go func(client *ssh.Client, config NodeConfig) {
			defer wg.Done()
			if err := startRemoteObshell(client, config.workDir, config.ip, config.obshellPort, &password); err != nil {
				errChan <- err
			}
		}(client.Client, config)
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		return err
	}

	// wait for takeover
	for times := 0; times < 60; times++ {
		count := 0
		for config := range clientMap {
			client, err := services.NewClientWithPassword(config.ip, config.obshellPort, password)
			if err != nil {
				continue
			}
			if status, err := client.V1().GetStatus(); err != nil {
				continue
			} else if status.Agent.Identity == model.CLUSTER_AGENT {
				count += 1
			} else if status.Agent.Identity == model.TAKE_OVER_MASTER {
				if dag, err := client.V1().GetAgentLastMaintenanceDag(); err != nil || dag == nil {
					break
				} else {
					if dag, err = client.V1().WaitDagSucceed(dag.GenericID); dag != nil && dag.IsFailed() {
						return fmt.Errorf("takeover failed: %s", err)
					} else if err != nil {
						break // Request failed, try again
					}
					return nil
				}
			}
		}

		if count == len(clientMap) {
			return nil
		}
		time.Sleep(10 * time.Second)
	}
	return errors.New("takeover timeout")
}
