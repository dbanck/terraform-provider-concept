package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func (p ConceptProvider) UpgradeResourceIdentity(ctx context.Context, request *tfprotov6.UpgradeResourceIdentityRequest) (*tfprotov6.UpgradeResourceIdentityResponse, error) {
	panic("not implemented")
}

func (p ConceptProvider) UpgradeResourceState(ctx context.Context, request *tfprotov6.UpgradeResourceStateRequest) (*tfprotov6.UpgradeResourceStateResponse, error) {
	panic("not implemented")
}

func (p ConceptProvider) ReadResource(ctx context.Context, request *tfprotov6.ReadResourceRequest) (*tfprotov6.ReadResourceResponse, error) {
	panic("not implemented")
}

func (p ConceptProvider) PlanResourceChange(ctx context.Context, request *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	panic("not implemented")
}

func (p ConceptProvider) ApplyResourceChange(ctx context.Context, request *tfprotov6.ApplyResourceChangeRequest) (*tfprotov6.ApplyResourceChangeResponse, error) {
	panic("not implemented")
}

func (p ConceptProvider) MoveResourceState(ctx context.Context, request *tfprotov6.MoveResourceStateRequest) (*tfprotov6.MoveResourceStateResponse, error) {
	panic("not implemented")
}

func (p ConceptProvider) ImportResourceState(ctx context.Context, request *tfprotov6.ImportResourceStateRequest) (*tfprotov6.ImportResourceStateResponse, error) {
	return &tfprotov6.ImportResourceStateResponse{}, nil
}
