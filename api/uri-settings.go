package api

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type SettingItems struct {
	ImageUri  string
	ReleaseAt time.Time
}

// 指定ファイルから設定を取得
func ReadSettingFromFile(settingFile string) (*SettingItems, error) {
	f, err := os.Open(settingFile)
	if err != nil {
		return nil, fmt.Errorf("ファイルがありません : %s", settingFile)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	imageUri := scanner.Text()
	scanner.Scan()
	tmpReleaseAt := scanner.Text()
	releaseAt, err := time.Parse("2006-01-02T15:04:05Z07:00", tmpReleaseAt)
	if err != nil {
		return nil, fmt.Errorf("リリース日時の形式が誤っています : %s", tmpReleaseAt)
	}
	return &SettingItems{
		ImageUri:  imageUri,
		ReleaseAt: releaseAt,
	}, nil
}

// 次回リリース設定読み込み
func ReadNextRelease(settingPath string, serviceName string) *NextRelease {
	var err error
	// リリース設定ファイルがあればその情報を返す
	settingFile := fmt.Sprintf("%s/%s-release-setting", settingPath, serviceName)
	settingItems, err := ReadSettingFromFile(settingFile)
	if err == nil {
		return &NextRelease{
			ImageUri:   settingItems.ImageUri,
			ReleaseAt:  settingItems.ReleaseAt,
		}
	}
	return nil
}

// 前回リリース時設定読み込み
func ReadLastReleased(settingPath string, serviceName string) *LastReleased {
	var err error
	// リリース済みの設定ファイルがあればその情報を返す
	settingFile := fmt.Sprintf("%s/%s-released", settingPath, serviceName)
	settingItems, err := ReadSettingFromFile(settingFile)
	if err == nil {
		return &LastReleased{
			ImageUri:   settingItems.ImageUri,
			ReleasedAt: settingItems.ReleaseAt,
		}
	}
	return nil
}

// ファイルシステムから必要情報を取得
func ReadSettings(pathPrefix string) (*[]UriSetting, error) {
	var err error
	var repositoryItems []UriSetting
	// ワイルドカードに一致するディレクトリ名をすべて取得
    fl, err := filepath.Glob(fmt.Sprintf("%s*", pathPrefix))
	if err != nil {
		return nil, err
	}
    for _, f := range fl {
        fInfo, err := os.Stat(f)
		if err != nil {
			return nil, err
		}
		if fInfo.IsDir() {
			// ディレクトリならディレクトリ内の services ファイルからサービス一覧を取得
			environmentName := f[len(pathPrefix):]
			servicesFile := fmt.Sprintf("%s/services", f)
			fs, err := os.Open(servicesFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "サービス一覧ファイル（%s）の読み取りに失敗しました\n: %s", servicesFile, err)
				return nil, err
			}
			defer fs.Close()

			scanner := bufio.NewScanner(fs)
			for scanner.Scan() {
				serviceName := scanner.Text()
				// 環境名・サービス名から設定を取得
				nextRelease := ReadNextRelease(f, serviceName)
				lastReleased := ReadLastReleased(f, serviceName)
				repositoryItems = append(repositoryItems, UriSetting{
					EnvironmentName: environmentName,
					ServiceName: serviceName,
					NextRelease: nextRelease,
					LastReleased: lastReleased,
				})
			}
		}
    }
    // 環境名・サービス名をソート
    sort.Slice(repositoryItems, func(i, j int) bool {
        return fmt.Sprintf("%s %s", repositoryItems[i].EnvironmentName, repositoryItems[i].ServiceName) < fmt.Sprintf("%s %s", repositoryItems[j].EnvironmentName, repositoryItems[j].ServiceName)
    })

	return &repositoryItems, err
}
