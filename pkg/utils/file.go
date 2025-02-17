package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Project-IPCA/ipca-backend/pkg/constants"
	"github.com/google/uuid"
)

const (
	TIME_LIMIT      = 3 * time.Second
	MAX_OUTPUT_SIZE = 1024 * 1024
)

func GetPythonKeywordFromCode(sourceCode string) (*constants.ReceiveGetKeyWordPython, error) {
	sandboxPath, err := InitializeIsolate()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize isolate: %v", err)
	}
	defer CleanupIsolate()

	keywordListPath := "./python_file/keyword_list.py"
	sandboxKeywordListPath := filepath.Join(sandboxPath, "box", "keyword_list.py")
	err = MoveFile(keywordListPath, sandboxKeywordListPath)
	if err != nil {
		return nil, fmt.Errorf("failed to copy keyword_list.py to sandbox: %v", err)
	}

	sourceCodeFileName := fmt.Sprintf("%s.txt", uuid.New().String())
	sandboxSourceCodePath := filepath.Join(sandboxPath, "box", sourceCodeFileName)
	err = os.WriteFile(sandboxSourceCodePath, []byte(strings.TrimSpace(sourceCode)), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create sourcecode file in sandbox: %v", err)
	}

	command := fmt.Sprintf("/usr/bin/python3.12 keyword_list.py %s", sourceCodeFileName)
	output, err := ExecuteCommandWithIsolate(sandboxPath, command)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	var result constants.ReceiveGetKeyWordPython
	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON output: %w", err)
	}
	return &result, nil
}

func GetCKeywordFromCode(sourceCode string) (*constants.ReceiveGetKeyWordC, error) {
	sandboxPath, err := InitializeIsolate()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize isolate: %v", err)
	}
	defer CleanupIsolate()

	ckeywordListPath := "./c_file/program/keyword_list"
	sandboxCKeywordListPath := filepath.Join(sandboxPath, "box", "keyword_list")
	copyKeywordListCommand := fmt.Sprintf("cp %s %s", ckeywordListPath, sandboxCKeywordListPath)
	_, err = ExecuteCommand(copyKeywordListCommand)
	if err != nil {
		return nil, fmt.Errorf("failed to copy keyword_list to sandbox: %v", err)
	}

	sourceCodeFileName := fmt.Sprintf("%s.txt", uuid.New().String())
	sandboxSourceCodePath := filepath.Join(sandboxPath, "box", sourceCodeFileName)
	err = os.WriteFile(sandboxSourceCodePath, []byte(strings.TrimSpace(sourceCode)), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create sourcecode file in sandbox: %v", err)
	}

	command := fmt.Sprintf("./keyword_list %s", sourceCodeFileName)
	output, err := ExecuteCommandWithIsolate(sandboxPath, command)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	var result constants.ReceiveGetKeyWordC
	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON output: %w", err)
	}
	return &result, nil
}

func PythonKeywordCheck(sourceCode string, exerciseKwList constants.PythonCheckKeywordCategory) (*constants.PythonCheckKeywordResponse, error) {
	sandboxPath, err := InitializeIsolate()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize isolate: %v", err)
	}
	defer CleanupIsolate()

	keywordListPath := "./python_file/kw_checker.py"
	sandboxKeywordListPath := filepath.Join(sandboxPath, "box", "kw_checker.py")
	err = MoveFile(keywordListPath, sandboxKeywordListPath)
	if err != nil {
		return nil, fmt.Errorf("failed to copy kw_checker.py to sandbox: %v", err)
	}

	sourceCodeFileName := fmt.Sprintf("%s.txt", uuid.New().String())
	sandboxSourceCodePath := filepath.Join(sandboxPath, "box", sourceCodeFileName)
	err = os.WriteFile(sandboxSourceCodePath, []byte(strings.TrimSpace(sourceCode)), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create sourcecode file in sandbox: %v", err)
	}

	exerciseKwListByte, err := json.Marshal(exerciseKwList)
	if err != nil {
		return nil, fmt.Errorf("failed to convert object to JSON: %v", err)
	}

	// Trim spaces from the JSON string
	trimmedJSON := strings.TrimSpace(string(exerciseKwListByte))
	exerciseKwListFileName := fmt.Sprintf("%s.txt", uuid.New().String())
	sandboxexerciseKwListPath := filepath.Join(sandboxPath, "box", exerciseKwListFileName)
	err = os.WriteFile(sandboxexerciseKwListPath, []byte(trimmedJSON), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create exerciseKwList file in sandbox: %v", err)
	}

	command := fmt.Sprintf("/usr/bin/python3.12 kw_checker.py %s %s", sourceCodeFileName, exerciseKwListFileName)
	fmt.Println(command)
	output, err := ExecuteCommandWithIsolate(sandboxPath, command)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	var result constants.PythonCheckKeywordResponse
	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON output: %w", err)
	}
	return &result, nil
}

func CKeywordCheck(sourceCode string, exerciseKwList constants.CCheckKeywordCategory) (*constants.CCheckKeywordResponse, error) {
	sandboxPath, err := InitializeIsolate()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize isolate: %v", err)
	}
	defer CleanupIsolate()

	ckeywordListPath := "./c_file/program/kw_checker"
	sandboxCKeywordListPath := filepath.Join(sandboxPath, "box", "kw_checker")
	copyKeywordListCommand := fmt.Sprintf("cp %s %s", ckeywordListPath, sandboxCKeywordListPath)
	_, err = ExecuteCommand(copyKeywordListCommand)
	if err != nil {
		return nil, fmt.Errorf("failed to copy kw_checker to sandbox: %v", err)
	}

	sourceCodeFileName := fmt.Sprintf("%s.txt", uuid.New().String())
	sandboxSourceCodePath := filepath.Join(sandboxPath, "box", sourceCodeFileName)
	err = os.WriteFile(sandboxSourceCodePath, []byte(strings.TrimSpace(sourceCode)), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create sourcecode file in sandbox: %v", err)
	}

	exerciseKwListByte, err := json.Marshal(exerciseKwList)
	if err != nil {
		return nil, fmt.Errorf("failed to convert object to JSON: %v", err)
	}

	// Trim spaces from the JSON string
	trimmedJSON := strings.TrimSpace(string(exerciseKwListByte))
	exerciseKwListFileName := fmt.Sprintf("%s.txt", uuid.New().String())
	sandboxexerciseKwListPath := filepath.Join(sandboxPath, "box", exerciseKwListFileName)
	err = os.WriteFile(sandboxexerciseKwListPath, []byte(trimmedJSON), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create exerciseKwList file in sandbox: %v", err)
	}

	command := fmt.Sprintf("./kw_checker %s %s", sourceCodeFileName, exerciseKwListFileName)
	fmt.Println(command)
	output, err := ExecuteCommandWithIsolate(sandboxPath, command)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	var result constants.CCheckKeywordResponse
	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON output: %w", err)
	}
	return &result, nil
}

func ExecuteCommandWithIsolate(sandboxPath, command string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIME_LIMIT)
	defer cancel()

	isolateCmd := fmt.Sprintf("isolate --run --time=%d --wall-time=%d --extra-time=1 --mem=128000 -- %s",
		int(TIME_LIMIT.Seconds()),
		int(TIME_LIMIT.Seconds()),
		command)

	cmd := exec.CommandContext(ctx, "bash", "-c", isolateCmd)
	cmd.Dir = filepath.Join(sandboxPath, "box")

	output, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("time limit exceeded: %w", ctx.Err())
		}
		fmt.Println(err)
		return nil, fmt.Errorf("error while running command: %w", err)
	}
	return output, nil
}

func ExecuteCommand(command string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIME_LIMIT)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(cmdOutput), nil
}

func MoveFile(sourcePath, destinationPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("not have source file : %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	fmt.Printf("File copied successfully from %s to %s\n", sourcePath, destinationPath)
	return nil
}

func InitializeIsolate() (string, error) {
	cmd := exec.Command("isolate", "--init")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to initialize isolate: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func CleanupIsolate() error {
	cmd := exec.Command("isolate", "--cleanup")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to cleanup isolate: %v", err)
	}
	return nil
}

func CreateTempFile(fileName string, sourceCode string) (*os.File, error) {
	tempFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return nil, fmt.Errorf("error while creating temp file: %v", err)
	}

	_, err = tempFile.Write([]byte(sourceCode))
	if err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, fmt.Errorf("error while writing to temp file: %v", err)
	}

	_, err = tempFile.Seek(0, 0)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, fmt.Errorf("error while seeking temp file: %v", err)
	}
	return tempFile, nil
}
