package googlenl

import (
	"context"

	"cloud.google.com/go/language/apiv1/languagepb"
)

func AnalyzeEntities(ctx context.Context, text string) (*languagepb.AnalyzeEntitiesResponse, error) {
	return NLUClient.AnalyzeEntities(ctx, &languagepb.AnalyzeEntitiesRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
}

func AnalyzeSentiment(ctx context.Context, text string) (*languagepb.AnalyzeSentimentResponse, error) {
	return NLUClient.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
	})
}

func AnalyzeEntitySentiment(ctx context.Context, text string) (*languagepb.AnalyzeEntitySentimentResponse, error) {
	return NLUClient.AnalyzeEntitySentiment(ctx, &languagepb.AnalyzeEntitySentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
	})
}

func AnalyzeSyntax(ctx context.Context, text string) (*languagepb.AnnotateTextResponse, error) {
	return NLUClient.AnnotateText(ctx, &languagepb.AnnotateTextRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		Features: &languagepb.AnnotateTextRequest_Features{
			ExtractSyntax: true,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
}

func ClassifyText(ctx context.Context, text string) (*languagepb.ClassifyTextResponse, error) {
	return NLUClient.ClassifyText(ctx, &languagepb.ClassifyTextRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		ClassificationModelOptions: &languagepb.ClassificationModelOptions{
			ModelType: &languagepb.ClassificationModelOptions_V2Model_{
				V2Model: &languagepb.ClassificationModelOptions_V2Model{
					ContentCategoriesVersion: languagepb.ClassificationModelOptions_V2Model_V2,
				},
			},
		},
	})
}
