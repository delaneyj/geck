package generator

import (
	"context"
	"testing"

	geckpb "github.com/delaneyj/geck/pb/gen/geck/v1"
	"github.com/stretchr/testify/assert"
)

func TestGenerator(t *testing.T) {
	testBundle := &geckpb.BundleDefinition{
		Name:        "TestBundle",
		Description: "A test bundle",
		States: []*geckpb.StateDefinition{
			{
				Name: "Menu",
				States: []string{
					"Main",
					"Options",
					"Quit",
				},
			},
		},
		Components: []*geckpb.ComponentDefinition{
			{
				Name: "Position",
				Fields: []*geckpb.FieldDefinition{
					{
						Name:       "X",
						ResetValue: &geckpb.FieldDefinition_F32{F32: 0},
					},
					{
						Name:       "Y",
						ResetValue: &geckpb.FieldDefinition_F32{F32: 0},
					},
					{
						Name:       "Z",
						ResetValue: &geckpb.FieldDefinition_F32{F32: 0},
					},
				},
			},
			{
				Name: "Velocity",
				Fields: []*geckpb.FieldDefinition{
					{
						Name:       "X",
						ResetValue: &geckpb.FieldDefinition_F32{F32: 0},
					},
					{
						Name:       "Y",
						ResetValue: &geckpb.FieldDefinition_F32{F32: 0},
					},
					{
						Name:       "Z",
						ResetValue: &geckpb.FieldDefinition_F32{F32: 0},
					},
				},
			},
			{
				Name: "Rotation",
				Fields: []*geckpb.FieldDefinition{
					{
						Name:       "X",
						ResetValue: &geckpb.FieldDefinition_F32{F32: 0},
					},
					{
						Name:       "Y",
						ResetValue: &geckpb.FieldDefinition_F32{F32: 0},
					},
					{
						Name:       "Z",
						ResetValue: &geckpb.FieldDefinition_F32{F32: 0},
					},
					{
						Name:       "W",
						ResetValue: &geckpb.FieldDefinition_F32{F32: 1},
					},
				},
			},
		},
		Queries: []*geckpb.QueryDefinition{
			{
				Name: "UpdatePosition",
				Group: &geckpb.QueryDefinition_Group{
					Op: geckpb.QueryDefinition_AND,
					Terms: []*geckpb.QueryDefinition_Term{
						{
							EntityName:   "Position",
							VariableName: "p",
							IsMutable:    true,
						},
						{
							EntityName:   "Velocity",
							VariableName: "v",
						},
					},
				},
			},
		},
	}

	opts := OutputOptions{
		Path:        "out",
		PackageName: "out",
	}
	assert.NoError(t, GenerateECS(context.Background(), opts, testBundle))
}
