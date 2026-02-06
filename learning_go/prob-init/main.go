package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// 引数チェック
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: prob-init <module-name>")
		os.Exit(1)
	}
	modName := os.Args[1]

	// 作成予定ディレクトリパス
	dirPath := filepath.Join(".", modName)

	// 既存チェック
	if _, err := os.Stat(dirPath); err == nil {
		fmt.Fprintf(os.Stderr, "directory already exists: %s\n", modName)
		os.Exit(1)
	} else if !os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// ディレクトリ作成
	if err := os.Mkdir(dirPath, 0755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// go mod init 実行
	cmd := exec.Command("go", "mod", "init", modName)
	cmd.Dir = dirPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// テンプレ生成
	mainFilePath := filepath.Join(dirPath, "main.go")
	template := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`

	if err := os.WriteFile(mainFilePath, []byte(template), 0644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("created:", dirPath)
}
