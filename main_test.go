package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/deepmap/oapi-codegen/pkg/testutil"
	"github.com/hmatsu47/list-settings-api/api"
	"github.com/hmatsu47/list-settings-api/testdouble"
	"github.com/stretchr/testify/assert"
)

func doGet(t *testing.T, handler http.Handler, url string) *httptest.ResponseRecorder {
	response := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, handler)
	return response.Recorder
}

// ファイルコピー
func fileCopy(srcPath string, dstPath string) (string, error) {
	src, err := os.Open(srcPath)
	if err != nil {
		return srcPath, err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return dstPath, err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return dstPath, err
	}
	return dstPath, err
}

// 設定をテンポラリディレクトリにコピー
func initConfig(templateConfigDir string) string {
	var err error
	tmpConfigDir, err := os.MkdirTemp("", "list-settings-test-config")
	if err != nil {
		panic(err)
	}
	fmt.Printf("テスト用のテンポラリディレクトリ（%s）を作成しました\n", tmpConfigDir)
	dirs, err := ioutil.ReadDir(templateConfigDir)
	if err != nil {
		panic(err)
	}
	// 設定ディレクトリをコピー
	for _, dir := range dirs {
		if dir.IsDir() {
			// 設定ディレクトリ作成
			err = os.Mkdir(filepath.Join(tmpConfigDir, dir.Name()), 0755)
			if err != nil {
				panic(err)
			}
			// 設定ディレクトリ内をコピー
			files, err := ioutil.ReadDir(filepath.Join(templateConfigDir, dir.Name()))
			if err != nil {
				panic(err)
			}
			for _, file := range files {
				_, err = fileCopy(filepath.Join(templateConfigDir, dir.Name(), file.Name()), filepath.Join(tmpConfigDir, dir.Name(), file.Name()))
				if err != nil {
					panic(err)
				}
			}
		}
	}
	return tmpConfigDir
}

// テンポラリディレクトリを削除
func clearTempDir(tmpDir string) {
	os.RemoveAll(tmpDir)
	fmt.Printf("テスト用のテンポラリディレクトリ（%s）を削除しました\n", tmpDir)
}

// go test -v で実行する
func TestListSettings1(t *testing.T) {
	var err error
	templateConfigDir := "./test/config1"
	workDir := initConfig(templateConfigDir)
	configPathPrefix := fmt.Sprintf("%s/select-repository-", workDir)
	tagRepositoryUri := ""
	var tagKeys []api.TagKey
	listSettings := api.NewListSettings(configPathPrefix, tagRepositoryUri, &tagKeys)
	var origins []string = []string{"http://example.com"}

	t.Cleanup(func() {
		clearTempDir(workDir)
	})

	ginListSettingsServer := NewGinListSettingsServer(listSettings, 28080, origins)
	r := ginListSettingsServer.Handler

	t.Run("uriSettingsのみ有効・URI形式での設定一覧取得", func(t *testing.T) {
		rr := doGet(t, r, "/uriSettings")

		var uriSettingList []api.UriSetting
		err = json.NewDecoder(rr.Body).Decode(&uriSettingList)
		assert.NoError(t, err, "error getting response")
		assert.Equal(t, 10, len(uriSettingList))
		assert.Equal(t, "env1", uriSettingList[0].EnvironmentName)
		assert.Equal(t, "000000000000.dkr.ecr.ap-northeast-1.amazonaws.com/repository1:20230501-release", uriSettingList[0].NextRelease.ImageUri)
		expectedTime111, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2023-04-30T19:05:00Z")
		assert.Equal(t, expectedTime111, uriSettingList[0].NextRelease.ReleaseAt)
		assert.Equal(t, "000000000000.dkr.ecr.ap-northeast-1.amazonaws.com/repository1:20230331-release", uriSettingList[0].LastReleased.ImageUri)
		expectedTime112, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2023-03-31T23:50:00+09:00")
		assert.Equal(t, expectedTime112, uriSettingList[0].LastReleased.ReleasedAt)
		assert.Equal(t, "test1", uriSettingList[0].ServiceName)
		assert.Equal(t, "env1", uriSettingList[1].EnvironmentName)
		assert.Equal(t, "000000000000.dkr.ecr.ap-northeast-1.amazonaws.com/repository2:20230502-release", uriSettingList[1].NextRelease.ImageUri)
		expectedTime121, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2023-05-02T04:05:00+09:00")
		assert.Equal(t, expectedTime121, uriSettingList[1].NextRelease.ReleaseAt)
		assert.Nil(t, uriSettingList[1].LastReleased)
		assert.Equal(t, "test2", uriSettingList[1].ServiceName)
		assert.Equal(t, "env1", uriSettingList[2].EnvironmentName)
		assert.Nil(t, uriSettingList[2].NextRelease)
		assert.Equal(t, "000000000000.dkr.ecr.ap-northeast-1.amazonaws.com/repository3:20230430-release", uriSettingList[2].LastReleased.ImageUri)
		expectedTime132, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2023-04-29T23:50:00Z")
		assert.Equal(t, expectedTime132, uriSettingList[2].LastReleased.ReleasedAt)
		assert.Equal(t, "test3", uriSettingList[2].ServiceName)
		assert.Equal(t, "env1", uriSettingList[3].EnvironmentName)
		assert.Nil(t, uriSettingList[3].NextRelease)
		assert.Nil(t, uriSettingList[3].LastReleased)
		assert.Equal(t, "test4", uriSettingList[3].ServiceName)
		assert.Equal(t, "env2", uriSettingList[4].EnvironmentName)
		assert.Nil(t, uriSettingList[4].NextRelease)
		assert.Nil(t, uriSettingList[4].LastReleased)
		assert.Equal(t, "test1", uriSettingList[4].ServiceName)
		assert.Equal(t, "env2", uriSettingList[5].EnvironmentName)
		assert.Equal(t, "000000000000.dkr.ecr.ap-northeast-1.amazonaws.com/repository3:20230503-release", uriSettingList[5].NextRelease.ImageUri)
		expectedTime231, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2023-05-02T19:05:00Z")
		assert.Equal(t, expectedTime231, uriSettingList[5].NextRelease.ReleaseAt)
		assert.Equal(t, "000000000000.dkr.ecr.ap-northeast-1.amazonaws.com/repository3:20230330-release", uriSettingList[5].LastReleased.ImageUri)
		expectedTime232, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2023-03-30T23:50:00+09:00")
		assert.Equal(t, expectedTime232, uriSettingList[5].LastReleased.ReleasedAt)
		assert.Equal(t, "test3", uriSettingList[5].ServiceName)
		assert.Equal(t, "env2", uriSettingList[6].EnvironmentName)
		assert.Nil(t, uriSettingList[6].NextRelease)
		assert.Nil(t, uriSettingList[6].LastReleased)
		assert.Equal(t, "test4", uriSettingList[6].ServiceName)
		assert.Equal(t, "env3", uriSettingList[7].EnvironmentName)
		assert.Nil(t, uriSettingList[7].NextRelease)
		assert.Nil(t, uriSettingList[7].LastReleased)
		assert.Equal(t, "test1", uriSettingList[7].ServiceName)
		assert.Equal(t, "env3", uriSettingList[8].EnvironmentName)
		assert.Nil(t, uriSettingList[8].NextRelease)
		assert.Nil(t, uriSettingList[8].LastReleased)
		assert.Equal(t, "test2", uriSettingList[8].ServiceName)
		assert.Equal(t, "env3", uriSettingList[9].EnvironmentName)
		assert.Nil(t, uriSettingList[9].NextRelease)
		assert.Nil(t, uriSettingList[9].LastReleased)
		assert.Equal(t, "test4", uriSettingList[9].ServiceName)
	})
}

func TestListSettings2(t *testing.T) {
	var err error
	templateConfigDir := "./test/config2"
	workDir := initConfig(templateConfigDir)
	configPathPrefix := fmt.Sprintf("%s/select-repository-", workDir)
	tagRepositoryUri := "000000000000.dkr.ecr.ap-northeast-1.amazonaws.com/repository1"
	var tagKeys []api.TagKey = []api.TagKey{
		{
			TagName:         "latest",
			EnvironmentName: "dev",
		},
		{
			TagName:         "release",
			EnvironmentName: "prod",
		},
	}
	listSettings := api.NewListSettings(configPathPrefix, tagRepositoryUri, &tagKeys)
	var origins []string = []string{"http://example.com"}

	t.Cleanup(func() {
		clearTempDir(workDir)
	})

	ginListSettingsServer := NewGinListSettingsServer(listSettings, 28080, origins)
	r := ginListSettingsServer.Handler

	t.Run("uriSettings・tagSettings有効・URI形式での設定一覧取得", func(t *testing.T) {
		rr := doGet(t, r, "/uriSettings")

		var uriSettingList []api.UriSetting
		err = json.NewDecoder(rr.Body).Decode(&uriSettingList)
		assert.NoError(t, err, "error getting response")
		assert.Equal(t, 3, len(uriSettingList))
		assert.Equal(t, "env1", uriSettingList[0].EnvironmentName)
		assert.Equal(t, "000000000000.dkr.ecr.ap-northeast-1.amazonaws.com/repository1:20230501-release", uriSettingList[0].NextRelease.ImageUri)
		expectedTime111, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2023-04-30T19:05:00Z")
		assert.Equal(t, expectedTime111, uriSettingList[0].NextRelease.ReleaseAt)
		assert.Equal(t, "000000000000.dkr.ecr.ap-northeast-1.amazonaws.com/repository1:20230331-release", uriSettingList[0].LastReleased.ImageUri)
		expectedTime112, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2023-03-31T23:50:00+09:00")
		assert.Equal(t, expectedTime112, uriSettingList[0].LastReleased.ReleasedAt)
		assert.Equal(t, "test1", uriSettingList[0].ServiceName)
		assert.Equal(t, "env1", uriSettingList[1].EnvironmentName)
		assert.Nil(t, uriSettingList[1].NextRelease)
		assert.Nil(t, uriSettingList[1].LastReleased)
		assert.Equal(t, "test2", uriSettingList[1].ServiceName)
		assert.Equal(t, "env2", uriSettingList[2].EnvironmentName)
		assert.Nil(t, uriSettingList[2].NextRelease)
		assert.Nil(t, uriSettingList[2].LastReleased)
		assert.Equal(t, "test1", uriSettingList[2].ServiceName)
	})
}

func TestListSettings3(t *testing.T) {
	// テスト用のパラメーターを生成
	repositoryUri := "000000000000.dkr.ecr.ap-northeast-1.amazonaws.com/repository1"
	repositoryName := "repository1"
	registryId := "000000000000"
	maxResults := int32(1000)

	// テスト用の ImageIds を生成
	digest1 := "sha256:4d2653f861f1c4cb187f1a61f97b9af7adec9ec1986d8e253052cfa60fd7372f"
	env1 := "prod"
	tag1 := "release"
	imageId1 :=
		types.ImageIdentifier{
			ImageDigest: aws.String(digest1),
			ImageTag:    aws.String(tag1),
		}
	digest2 := "sha256:20b39162cb057eab7168652ab012ae3712f164bf2b4ef09e6541fca4ead3df62"
	env2 := "dev"
	tag2 := "latest"
	imageId2 :=
		types.ImageIdentifier{
			ImageDigest: aws.String(digest2),
		}
	digest3 := "sha256:0000000000000000000000000000000000000000000000000000000000000000"
	imageId3 :=
		types.ImageIdentifier{
			ImageDigest: aws.String(digest3),
		}
	var imageIds []types.ImageIdentifier
	imageIds = append(imageIds, imageId1)
	imageIds = append(imageIds, imageId2)
	imageIds = append(imageIds, imageId3)
	var tagKeys []api.TagKey = []api.TagKey{
		{
			TagName:         tag2,
			EnvironmentName: env2,
		},
		{
			TagName:         tag1,
			EnvironmentName: env1,
		},
	}

	// テスト用の ImageDetails を生成
	time1, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2023-05-10T05:27:02Z")
	time2, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2022-05-01T05:07:10Z")
	time3, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2022-05-11T04:18:10Z")
	size1Int64 := int64(10017365)
	var tags1 []string = []string{tag1, "hogefuga"}
	imageDetail1 :=
		types.ImageDetail{
			ImageDigest:      aws.String(digest1),
			ImagePushedAt:    aws.Time(time1),
			ImageSizeInBytes: aws.Int64(size1Int64),
			ImageTags:        tags1,
			RegistryId:       aws.String(registryId),
			RepositoryName:   aws.String(repositoryName),
		}
	size2Int64 := int64(10017367)
	var tags2 []string = []string{"testtest", tag2}
	imageDetail2 :=
		types.ImageDetail{
			ImageDigest:      aws.String(digest2),
			ImagePushedAt:    aws.Time(time2),
			ImageSizeInBytes: aws.Int64(size2Int64),
			ImageTags:        tags2,
			RegistryId:       aws.String(registryId),
			RepositoryName:   aws.String(repositoryName),
		}
		size3Int64 := int64(10000)
		var tags3 []string = []string{"foobar"}
		imageDetail3 :=
			types.ImageDetail{
				ImageDigest:      aws.String(digest3),
				ImagePushedAt:    aws.Time(time3),
				ImageSizeInBytes: aws.Int64(size3Int64),
				ImageTags:        tags3,
				RegistryId:       aws.String(registryId),
				RepositoryName:   aws.String(repositoryName),
			}
	var imageDetails []types.ImageDetail
	imageDetails = append(imageDetails, imageDetail1)
	imageDetails = append(imageDetails, imageDetail2)
	imageDetails = append(imageDetails, imageDetail3)

	// テスト用の Images（ECR）を生成
	manifest := "{\"test\":\"testtext\"}"
	image1 := types.Image{
		ImageId:        &imageIds[0],
		ImageManifest:  aws.String(manifest),
		RegistryId:     aws.String(registryId),
		RepositoryName: aws.String(repositoryName),
	}
	var images []types.Image
	images = append(images, image1)

	// テストケース
	testParams := testdouble.ECRParams{
		RepositoryName:  repositoryName,
		RegistryId:      registryId,
		ImageDetails:    imageDetails,
		MaxResults:      maxResults,
	}
	mockParams := testdouble.MockECRParams{
		ECRParams: testParams,
	}

	t.Run("タグ形式での設定一覧取得（モック利用）", func(t *testing.T) {
		ecrClient := func(t *testing.T) testdouble.MockECRAPI {
			return testdouble.GenerateMockECRAPI(mockParams)
		}
		ctx := context.TODO()
		// TagSettingList のテスト
		settingList, err := api.TagSettingList(ctx, ecrClient(t), repositoryUri, tagKeys)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(settingList))
		assert.Equal(t, "dev", settingList[0].EnvironmentName)
		assert.Equal(t, tags2, settingList[0].Tags)
		expectedTime1, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2022-05-01T05:07:10Z")
		assert.Equal(t, expectedTime1, *settingList[0].PushedAt)
		assert.Equal(t, "prod", settingList[1].EnvironmentName)
		assert.Equal(t, tags1, settingList[1].Tags)
		expectedTime2, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2023-05-10T05:27:02Z")
		assert.Equal(t, expectedTime2, *settingList[1].PushedAt)
	})
	// このモック利用テスト実施との差分を考えると GetTagSettingList を単独でテストする意味はないので実施しない
}
