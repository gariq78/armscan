package updater

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	agent "ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/usecases"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs"
)

type Restarter struct {
	logg         agent.Logger
	currFileName string
}

func New(logger agent.Logger, currFileName string) *Restarter {
	return &Restarter{
		logg:         logger,
		currFileName: currFileName,
	}
}

func (r *Restarter) Update(archive structs.AgentArchive) error {
	if archive.Zip != nil {
		rdr := bytes.NewReader(archive.Zip)
		files, err := r.unzipBytes(rdr, int64(len(archive.Zip)), ".")
		if err != nil {
			return fmt.Errorf("unzipBytes %d err: %w", len(archive.Zip), err)
		}

		err = os.Rename(r.currFileName, "~"+r.currFileName)
		if err != nil {
			return fmt.Errorf("rename currFileName err: %w", err)
		}

		err = os.Rename(files[0], r.currFileName)
		if err != nil {
			_ = os.Rename("~"+r.currFileName, r.currFileName)
			return fmt.Errorf("rename err: %w", err)
		}

		var procname = getOsDepCmd(r.currFileName)
		cmd := exec.Command(procname)
		var out, eee bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &eee
		err = cmd.Start()
		if err != nil {
			return fmt.Errorf("start %q err: %w", procname, err)
		}

		// Непонятно как это победить, без задержки дочерний процесс убивается этим при Interrupt
		time.Sleep(100 * time.Millisecond)

		log.Printf("Output: %s", out.String())
		log.Printf("Eee: %s", eee.String())

		p, err := os.FindProcess(os.Getpid())
		if err != nil {
			return fmt.Errorf("findprocess pid=%d err: %w", os.Getpid(), err)
		}

		err = p.Signal(os.Interrupt)
		if err != nil {
			return fmt.Errorf("self interrupt err: %w", err)
		}
	}

	return nil
}

func (r *Restarter) unzipBytes(b io.ReaderAt, size int64, targetPath string) ([]string, error) {
	reader, err := zip.NewReader(b, size)
	if err != nil {
		return nil, fmt.Errorf("zip.NewReader err: %w", err)
	}
	return r.unzip(reader, targetPath)
}

func (r *Restarter) unzip(reader *zip.Reader, targetPath string) ([]string, error) {
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return nil, fmt.Errorf("os.MkdirAll(%q) err: %w", targetPath, err)
	}

	var res = make([]string, 0, len(reader.File))
	for _, file := range reader.File {
		path := filepath.Join(targetPath, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("file.Open %q err: %w", file.Name, err)
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return nil, fmt.Errorf("os.OpenFile(%q) err: %w", path, err)
		}
		defer targetFile.Close()

		n, err := io.Copy(targetFile, fileReader)
		if err != nil {
			return nil, fmt.Errorf("io.Copy %q err: %w", file.Name, err)
		}

		r.logg.Inf("%d bytes extract to file %s\n", n, file.Name)
		res = append(res, file.Name)
	}

	return res, nil
}
