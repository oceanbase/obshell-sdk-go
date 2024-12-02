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
	"fmt"
	"io"
	"os"

	"github.com/cavaliergopher/cpio"
	"github.com/cavaliergopher/rpm"
	"github.com/pkg/errors"
	"github.com/ulikunitz/xz"
)

type rpmPackage struct {
	cpioReader   *cpio.Reader
	filePath     string
	rpmFile      *os.File
	architecture string
}

func newRpmPackage(filePath string) *rpmPackage {
	return &rpmPackage{
		filePath: filePath,
	}
}

func (rp *rpmPackage) open() (err error) {
	rp.close()
	rp.rpmFile, err = os.Open(rp.filePath)
	if err != nil {
		return errors.Wrap(err, "open rpm file failed")
	}

	defer func() {
		if err != nil {
			rp.close()
		}
	}()

	rpmData, err := rpm.Read(rp.rpmFile)
	if err != nil {
		return errors.Wrap(err, "rpm read failed")
	}
	if err = checkCompressAndFormat(rpmData); err != nil {
		return err
	}
	rp.architecture = rpmData.Architecture()

	xzReader, err := xz.NewReader(rp.rpmFile)
	if err != nil {
		return errors.Wrap(err, "new xz reader failed")
	}

	rp.cpioReader = cpio.NewReader(xzReader)
	return nil
}

func (rp *rpmPackage) close() {
	if rp.rpmFile != nil {
		rp.rpmFile.Close()
		rp.cpioReader = nil
	}
}

func (rp *rpmPackage) next() (*cpio.Header, error) {
	if rp.cpioReader == nil {
		return nil, errors.New("rpm package not opened")
	}

	header, err := rp.cpioReader.Next()
	if err == io.EOF {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "cpio read failed")
	}
	return header, nil
}

func checkCompressAndFormat(pkg *rpm.Package) error {
	if pkg.PayloadCompression() != "xz" {
		return fmt.Errorf("unsupported compression '%s', the supported compression is 'xz'", pkg.PayloadCompression())
	}
	if pkg.PayloadFormat() != "cpio" {
		return fmt.Errorf("unsupported payload format '%s', the supported payload format is 'cpio'", pkg.PayloadFormat())
	}
	return nil
}
