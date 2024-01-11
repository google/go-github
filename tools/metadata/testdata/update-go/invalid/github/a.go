package github

type AService struct{}

// NoOperation has no operation
func (*AService) NoOperation() {}

func (*AService) NoComment() {}

// Ambiguous has an operation that could resolve to multiple operations
//
//meta:operation GET /ambiguous/{}
func (*AService) Ambiguous() {}

// MissingOperation has an operation that is missing from the OpenAPI spec
//
//meta:operation GET /missing/{id}
func (*AService) MissingOperation() {}

// DuplicateOperations has duplicate operations
//
//meta:operation GET /a/{a_id}
//meta:operation POST /a/{a_id}
//meta:operation GET /a/{a_id}
func (*AService) DuplicateOperations() {}
