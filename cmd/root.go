package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// flags
var (
	hasCopy      bool
	isSequential bool
	fileName     string
)

// configs
var (
	captureDir         string
	captureFilePattern string
	extension          string
	sequestialDigits   int
)

var RootCmd = &cobra.Command{
	// アプリケーションコマンド名
	Use: "mvsc",
	// mvsc --helpで表示される文言
	Short: "move screen capture",
	// 想定される引数
	Args: cobra.RangeArgs(1, 1),
	// 実行処理
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(args[0])
	},
}

func init() {
	readConfig()
	// copy指定オプション
	RootCmd.Flags().BoolVarP(&hasCopy, "copy", "c", false, "copy capture file")
	// ファイル名連番オプション
	RootCmd.Flags().BoolVarP(&isSequential, "seq", "s", false, "rename sequential numbering")
	// ファイル名指定オプション
	RootCmd.Flags().StringVarP(&fileName, "file", "f", "", "file name")
}

// 設定ファイルを読み込む
func readConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(home, ".app", "mvsc"))
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	captureDir = viper.GetString("captureDir")
	captureFilePattern = viper.GetString("captureFilePattern")
	extension = viper.GetString("extension")
	sequestialDigits = viper.GetInt("sequestialDigits")
}

// ファイルのコピー処理
func copy(srcPath, dstPath string) {
	src, err := os.Open(srcPath)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		panic(err)
	}
}

// 五十音順で末尾になる要素を返す、要素数が0の場合はnilを返す
// 引数を破壊的にソートする
func getEnd(s []string) *string {
	if len(s) == 0 {
		return nil
	}
	// 逆順でソートして0番目を返している
	sort.SliceStable(s, func(i, j int) bool { return s[i] > s[j] })
	return &s[0]
}

// 連番ファイル名を返す
func getSequentialFileName(dstDir string) string {
	var pattern = strings.Repeat("[0-9]", sequestialDigits) + extension
	// 移動先パスの連番ファイル一覧を取得
	dstFiles, err := filepath.Glob(filepath.Join(dstDir, pattern))
	if err != nil {
		panic(err)
	}
	// 最新のファイルパスを取得
	latest := getEnd(dstFiles)
	if latest == nil {
		// "00~~1.png"を返す
		return strings.Repeat("0", sequestialDigits-1) + "1" + extension
	}
	// ファイル名を取得
	fn := filepath.Base(*latest)
	// 末尾の".png"を削除して番号を取得
	rep := regexp.MustCompile(extension + "$")
	n := rep.ReplaceAllString(fn, "")
	// 数値に変換
	num, err := strconv.Atoi(n)
	if err != nil {
		panic(err)
	}
	// 次の番号のファイル名を返す
	num++
	f := "%0" + strconv.Itoa(sequestialDigits) + "d"
	return fmt.Sprintf(f, num) + extension
}

// メイン処理
func run(dstDir string) error {
	// -fオプションと-sオプションが同時に指定されていたらエラーとする
	if isSequential && fileName != "" {
		return errors.New("Do not specify A and B at the same time")
	}
	// コマンドに渡された移動先パスが存在するか判定
	if f, err := os.Stat(dstDir); os.IsNotExist(err) || !f.IsDir() {
		return fmt.Errorf("dir not found > %v\n", dstDir)
	}
	// 移動先パスを絶対パスに変換
	if !filepath.IsAbs(dstDir) {
		var err error
		dstDir, err = filepath.Abs(dstDir)
		if err != nil {
			return fmt.Errorf("dirpath cannot convert to abs > %v \n", dstDir)
		}
	}
	// homeディレクトリを取得
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	globPattern := filepath.Join(home, captureDir, captureFilePattern+extension)
	// 指定フォーマットの画像ファイル一覧を取得
	capturePaths, err := filepath.Glob(globPattern)
	if err != nil {
		panic(err)
	}
	// 最新のスクリーンショットのパスを取得
	latest := getEnd(capturePaths)
	if latest == nil {
		return fmt.Errorf("screenshot not found\n")
	}
	latestCapturePath := *latest

	// 移動先のファイル名
	var dstFileName string
	if isSequential {
		// 連番のファイル名を設定
		dstFileName = getSequentialFileName(dstDir)
	} else {
		if fileName == "" {
			// コピー元のファイル名を設定
			dstFileName = filepath.Base(latestCapturePath)
		} else {
			// 実行時に指定されたファイル名を設定
			dstFileName = fileName + extension
		}
	}

	// 移動先のフルパス
	dstPath := filepath.Join(dstDir, dstFileName)
	if hasCopy {
		// 画像ファイルをコピーする
		copy(latestCapturePath, dstPath)
	} else {
		// 画像ファイルを移動する
		if err := os.Rename(latestCapturePath, dstPath); err != nil {
			panic(err)
		}
	}
	return nil
}
