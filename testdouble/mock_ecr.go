package testdouble

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

// モックパラメーター
type ECRParams struct {
	RepositoryName  string
	RegistryId      string
	ImageDetails    []types.ImageDetail
	MaxResults      int32
}

// モック生成用
type MockECRParams struct {
	ECRParams ECRParams
}

// モック化
type MockECRAPI struct {
	DescribeImagesAPI MockECRDescribeImagesAPI
}

type MockECRDescribeImagesAPI func(ctx context.Context, params *ecr.DescribeImagesInput, optFns ...func(*ecr.Options)) (*ecr.DescribeImagesOutput, error)

func (m MockECRAPI) DescribeImages(ctx context.Context, params *ecr.DescribeImagesInput, optFns ...func(*ecr.Options)) (*ecr.DescribeImagesOutput, error) {
	return m.DescribeImagesAPI(ctx, params, optFns...)
}
