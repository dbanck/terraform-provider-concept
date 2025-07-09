package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func (p ConceptProvider) OpenEphemeralResource(ctx context.Context, request *tfprotov6.OpenEphemeralResourceRequest) (*tfprotov6.OpenEphemeralResourceResponse, error) {
	panic("not implemented")
}

func (p ConceptProvider) RenewEphemeralResource(ctx context.Context, request *tfprotov6.RenewEphemeralResourceRequest) (*tfprotov6.RenewEphemeralResourceResponse, error) {
	panic("not implemented")
}

func (p ConceptProvider) CloseEphemeralResource(ctx context.Context, request *tfprotov6.CloseEphemeralResourceRequest) (*tfprotov6.CloseEphemeralResourceResponse, error) {
	panic("not implemented")
}
