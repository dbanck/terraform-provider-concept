package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func (p ConceptProvider) CallFunction(ctx context.Context, request *tfprotov6.CallFunctionRequest) (*tfprotov6.CallFunctionResponse, error) {
	panic("not implemented")
}

func (p ConceptProvider) GetFunctions(ctx context.Context, request *tfprotov6.GetFunctionsRequest) (*tfprotov6.GetFunctionsResponse, error) {
	panic("not implemented")
}
