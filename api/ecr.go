package api

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

// ECR API interface
type ECRAPI interface {
	EcrDescribeImagesAPI
}

// ECR クライアント生成
func EcrClient(region string) (*ecr.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("AWS（API）の認証に失敗しました : %s", err)
	}
	return ecr.NewFromConfig(cfg), nil
}

// ECR DescribeImages
type EcrDescribeImagesAPI interface {
	DescribeImages(ctx context.Context, params *ecr.DescribeImagesInput, optFns ...func(*ecr.Options)) (*ecr.DescribeImagesOutput, error)
}

func EcrDescribeImages(ctx context.Context, api EcrDescribeImagesAPI, repositoryName string, registryId string) ([]types.ImageDetail, error) {
	// ページネーションさせないために最大件数を 1,000 に（実際には数十個程度の想定）
	maxResults := int32(1000)

	ecrImages, err := api.DescribeImages(ctx, &ecr.DescribeImagesInput{
		RepositoryName: aws.String(repositoryName),
		RegistryId:     aws.String(registryId),
		MaxResults:     aws.Int32(maxResults),
	})
	if err != nil {
		return nil, fmt.Errorf("リポジトリ（%s）のイメージ詳細一覧の取得に失敗しました : %s", repositoryName, err)
	}
	return ecrImages.ImageDetails, nil
}

func contains(elems []string, v string) bool {
    for _, s := range elems {
        if v == s {
            return true
        }
    }
    return false
}

// TagSettingList を取得
func GetTagSettingList(imageDetails []types.ImageDetail, tagKeys []TagKey) []TagSetting {
	var settingList []TagSetting
	var envList []string
	for _, v := range imageDetails {
		tags := v.ImageTags

		if len(tags) > 0 {
			// タグがあるイメージのみ検索対象に
			pushedAt := v.ImagePushedAt
			for _, st := range tagKeys {
				if contains(tags, st.TagName) {
					// 対象イメージ
					setting := TagSetting{
						Tags:            tags,
						EnvironmentName: st.EnvironmentName,
						PushedAt:        pushedAt,
					}
					settingList = append(settingList, setting)
					envList = append(envList, st.EnvironmentName)
				}
			}
		}
	}
	// 対象イメージがなかったタグの補完
	for _, sl := range tagKeys {
		if !contains(envList, sl.EnvironmentName) {
			var dummyTags []string = []string{"（未指定）"}
			setting := TagSetting{
				Tags:            dummyTags,
				EnvironmentName: sl.EnvironmentName,
			}
			settingList = append(settingList, setting)
		}
	}
	// 結果を環境名でソート
	sort.Slice(settingList, func(i, j int) bool {
		return settingList[i].EnvironmentName < settingList[j].EnvironmentName
	})
	return settingList
}

// ECR リポジトリ内イメージ一覧取得
func TagSettingList(ctx context.Context, api ECRAPI, repositoryUri string, tagKeys []TagKey) ([]TagSetting, error) {
	repositoryName := strings.Split(repositoryUri, "/")[1]
	registryId := strings.Split(repositoryUri, ".")[0]

	imageDetails, err := EcrDescribeImages(ctx, api, repositoryName, registryId)
	if err != nil {
		return nil, err
	}

	imageList := GetTagSettingList(imageDetails, tagKeys)
	return imageList, nil
}
