// +build windows

package clipboard

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func write(file string) error {
	cmd := exec.Command("PowerShell", "-Command", "Add-Type", "-AssemblyName",
		fmt.Sprintf("System.Windows.Forms;[Windows.Forms.Clipboard]::SetImage([System.Drawing.Image]::FromFile('%s'));", file))
	b, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(b))
	}
	return nil
}

func read() (io.Reader, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}
	f.Close()
	defer os.Remove(f.Name())

	cmd := exec.Command("PowerShell", "-Command", "Add-Type", "-AssemblyName",
		fmt.Sprintf("System.Windows.Forms;$clip=[Windows.Forms.Clipboard]::GetImage();if ($clip -ne $null) { $clip.Save('%s') };", f.Name()))
	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err, string(b))
	}

	r := new(bytes.Buffer)
	f, err = os.Open(f.Name())
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.Size() == 0 {
		return nil, fmt.Errorf("no clipboard image")
	}

	if _, err := io.Copy(r, f); err != nil {
		return nil, err
	}

	return r, nil
}
