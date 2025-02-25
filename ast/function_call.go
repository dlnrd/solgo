package ast

import (
	"fmt"
	"github.com/goccy/go-json"
	"regexp"
	"strings"

	v3 "github.com/cncf/xds/go/xds/type/v3"
	ast_pb "github.com/unpackdev/protos/dist/go/ast"
	"github.com/unpackdev/solgo/parser"
)

// FunctionCall represents a function call node in the AST.
type FunctionCall struct {
	*ASTBuilder

	Id                    int64              `json:"id"`                               // Unique identifier for the node.
	NodeType              ast_pb.NodeType    `json:"node_type"`                        // Type of the node.
	Kind                  ast_pb.NodeType    `json:"kind"`                             // Kind of the node.
	Src                   SrcNode            `json:"src"`                              // Source location of the node.
	ArgumentTypes         []*TypeDescription `json:"argument_types"`                   // Types of the arguments.
	Arguments             []Node[NodeType]   `json:"arguments"`                        // Arguments of the function call.
	Expression            Node[NodeType]     `json:"expression"`                       // Expression of the function call.
	ReferencedDeclaration int64              `json:"referenced_declaration,omitempty"` // Referenced declaration of the function call.
	TypeDescription       *TypeDescription   `json:"type_description"`                 // Type description of the function call.
}

// NewFunctionCall creates a new FunctionCall node with a given ASTBuilder.
// It initializes the Arguments slice and sets the NodeType and Kind to FUNCTION_CALL.
func NewFunctionCall(b *ASTBuilder) *FunctionCall {
	return &FunctionCall{
		ASTBuilder:    b,
		Arguments:     make([]Node[NodeType], 0),
		ArgumentTypes: make([]*TypeDescription, 0),
		NodeType:      ast_pb.NodeType_FUNCTION_CALL,
		Kind:          ast_pb.NodeType_FUNCTION_CALL,
	}
}

// SetReferenceDescriptor sets the reference descriptions of the FunctionCall node.
func (f *FunctionCall) SetReferenceDescriptor(refId int64, refDesc *TypeDescription) bool {
	f.ReferencedDeclaration = refId
	f.TypeDescription = refDesc
	return false
}

// GetId returns the unique identifier of the FunctionCall node.
func (f *FunctionCall) GetId() int64 {
	return f.Id
}

// GetType returns the type of the FunctionCall node.
func (f *FunctionCall) GetType() ast_pb.NodeType {
	return f.NodeType
}

// GetSrc returns the source location of the FunctionCall node.
func (f *FunctionCall) GetSrc() SrcNode {
	return f.Src
}

// GetArguments returns the arguments of the FunctionCall node.
func (f *FunctionCall) GetArguments() []Node[NodeType] {
	return f.Arguments
}

// GetArgumentTypes returns the types of the arguments of the FunctionCall node.
func (f *FunctionCall) GetArgumentTypes() []*TypeDescription {
	return f.ArgumentTypes
}

// GetKind returns the kind of the FunctionCall node.
func (f *FunctionCall) GetKind() ast_pb.NodeType {
	return f.Kind
}

// GetExpression returns the expression of the FunctionCall node.
func (f *FunctionCall) GetExpression() Node[NodeType] {
	return f.Expression
}

// GetTypeDescription returns the type description of the FunctionCall node.
// Currently, it returns nil and needs to be implemented.
func (f *FunctionCall) GetTypeDescription() *TypeDescription {
	return f.TypeDescription
}

// GetNodes returns a slice of nodes that includes the expression of the FunctionCall node.
func (f *FunctionCall) GetNodes() []Node[NodeType] {
	toReturn := []Node[NodeType]{f.Expression}
	toReturn = append(toReturn, f.Arguments...)
	return toReturn
}

// GetReferenceDeclaration returns the referenced declaration of the FunctionCall node.
func (f *FunctionCall) GetReferenceDeclaration() int64 {
	return f.ReferencedDeclaration
}

// RebuildDescriptions rebuilds the type descriptions of the FunctionCall node. It is called after the AST is built.
func (f *FunctionCall) RebuildDescriptions() {
	var newArgs []*TypeDescription
	for _, arg := range f.GetArguments() {
		newArgs = append(newArgs, arg.GetTypeDescription())
	}
	f.ArgumentTypes = newArgs
	f.TypeDescription = f.buildTypeDescription()
}

// UnmarshalJSON unmarshals a given JSON byte array into a FunctionCall node.
func (f *FunctionCall) UnmarshalJSON(data []byte) error {
	var tempMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}

	if id, ok := tempMap["id"]; ok {
		if err := json.Unmarshal(id, &f.Id); err != nil {
			return err
		}
	}

	if nodeType, ok := tempMap["node_type"]; ok {
		if err := json.Unmarshal(nodeType, &f.NodeType); err != nil {
			return err
		}
	}

	if kind, ok := tempMap["kind"]; ok {
		if err := json.Unmarshal(kind, &f.Kind); err != nil {
			return err
		}
	}

	if src, ok := tempMap["src"]; ok {
		if err := json.Unmarshal(src, &f.Src); err != nil {
			return err
		}
	}

	if referencedDeclaration, ok := tempMap["referenced_declaration"]; ok {
		if err := json.Unmarshal(referencedDeclaration, &f.ReferencedDeclaration); err != nil {
			return err
		}
	}

	if typeDescription, ok := tempMap["type_description"]; ok {
		if err := json.Unmarshal(typeDescription, &f.TypeDescription); err != nil {
			return err
		}
	}

	if argTypes, ok := tempMap["argument_types"]; ok {
		if err := json.Unmarshal(argTypes, &f.ArgumentTypes); err != nil {
			return err
		}
	}

	if arguments, ok := tempMap["arguments"]; ok {
		f.Arguments = make([]Node[NodeType], 0)
		var nodes []json.RawMessage
		if err := json.Unmarshal(arguments, &nodes); err != nil {
			return err
		}

		for _, tempNode := range nodes {
			var tempNodeMap map[string]json.RawMessage
			if err := json.Unmarshal(tempNode, &tempNodeMap); err != nil {
				return err
			}

			var tempNodeType ast_pb.NodeType
			if err := json.Unmarshal(tempNodeMap["node_type"], &tempNodeType); err != nil {
				return err
			}

			node, err := unmarshalNode(tempNode, tempNodeType)
			if err != nil {
				return err
			}
			f.Arguments = append(f.Arguments, node)
		}
	}

	if expression, ok := tempMap["expression"]; ok {
		if err := json.Unmarshal(expression, &f.Expression); err != nil {
			var tempNodeMap map[string]json.RawMessage
			if err := json.Unmarshal(expression, &tempNodeMap); err != nil {
				return err
			}

			var tempNodeType ast_pb.NodeType
			if err := json.Unmarshal(tempNodeMap["node_type"], &tempNodeType); err != nil {
				return err
			}

			node, err := unmarshalNode(expression, tempNodeType)
			if err != nil {
				return err
			}
			f.Expression = node
		}
	}

	return nil
}

// ToProto returns a protobuf representation of the FunctionCall node.
// Currently, it returns an empty Statement and needs to be implemented.
func (f *FunctionCall) ToProto() NodeType {
	proto := ast_pb.FunctionCall{
		Id:                    f.GetId(),
		NodeType:              f.GetType(),
		Kind:                  f.GetKind(),
		Src:                   f.Src.ToProto(),
		ReferencedDeclaration: f.GetReferenceDeclaration(),
		TypeDescription:       f.GetTypeDescription().ToProto(),
		ArgumentTypes:         make([]*ast_pb.TypeDescription, 0),
	}

	if f.GetExpression() != nil {
		proto.Expression = f.GetExpression().ToProto().(*v3.TypedStruct)
	}

	for _, arg := range f.GetArguments() {
		proto.Arguments = append(proto.Arguments, arg.ToProto().(*v3.TypedStruct))
	}

	for _, argType := range f.GetArgumentTypes() {
		if argType == nil {
			continue
		}

		proto.ArgumentTypes = append(proto.ArgumentTypes, argType.ToProto())
	}

	return NewTypedStruct(&proto, "FunctionCall")
}

// Parse takes a parser.FunctionCallContext and parses it into a FunctionCall node.
// It sets the Id, Src, Arguments, ArgumentTypes, and Expression of the FunctionCall node.
// It returns the created FunctionCall node.
func (f *FunctionCall) Parse(
	unit *SourceUnit[Node[ast_pb.SourceUnit]],
	contractNode Node[NodeType],
	fnNode Node[NodeType],
	bodyNode *BodyNode,
	vDeclar *VariableDeclaration,
	expNode Node[NodeType],
	ctx *parser.FunctionCallContext,
) Node[NodeType] {
	f.Id = f.GetNextID()
	f.Src = SrcNode{
		Line:   int64(ctx.GetStart().GetLine()),
		Column: int64(ctx.GetStart().GetColumn()),
		Start:  int64(ctx.GetStart().GetStart()),
		End:    int64(ctx.GetStop().GetStop()),
		Length: int64(ctx.GetStop().GetStop() - ctx.GetStart().GetStart() + 1),
		ParentIndex: func() int64 {
			if vDeclar != nil {
				return vDeclar.GetId()
			}

			if expNode != nil {
				return expNode.GetId()
			}

			if bodyNode != nil {
				return bodyNode.GetId()
			}

			if fnNode != nil {
				return fnNode.GetId()
			}

			return contractNode.GetId()
		}(),
	}

	expression := NewExpression(f.ASTBuilder)

	if ctx.Expression() != nil {
		f.Expression = expression.Parse(
			unit, contractNode, fnNode, bodyNode, nil, f, f.GetId(), ctx.Expression(),
		)
	}

	if ctx.CallArgumentList() != nil {
		for _, expressionCtx := range ctx.CallArgumentList().AllExpression() {
			expr := expression.Parse(unit, contractNode, fnNode, bodyNode, nil, f, f.GetId(), expressionCtx)
			f.Arguments = append(
				f.Arguments,
				expr,
			)

			f.ArgumentTypes = append(
				f.ArgumentTypes,
				expr.GetTypeDescription(),
			)
		}
	}

	f.TypeDescription = f.buildTypeDescription()
	return f
}

// buildTypeDescription constructs and returns the TypeDescription of the FunctionCall.
func (f *FunctionCall) buildTypeDescription() *TypeDescription {
	typeString := "function("
	typeIdentifier := "t_function_"
	typeStrings := make([]string, 0)
	typeIdentifiers := make([]string, 0)

	for _, paramType := range f.GetArgumentTypes() {
		if paramType == nil {
			typeStrings = append(typeStrings, fmt.Sprintf("unknown_%d", f.GetId()))
			typeIdentifiers = append(typeIdentifiers, fmt.Sprintf("$_t_unknown_%d", f.GetId()))
			continue
		} else if strings.Contains(paramType.TypeString, "literal_string") {
			typeStrings = append(typeStrings, "string memory")
			typeIdentifiers = append(typeIdentifiers, "_"+paramType.TypeIdentifier)
			continue
		} else if strings.Contains(paramType.TypeString, "contract") {
			typeStrings = append(typeStrings, "address")
			typeIdentifiers = append(typeIdentifiers, "$_t_address")
			continue
		}

		typeStrings = append(typeStrings, paramType.TypeString)
		typeIdentifiers = append(typeIdentifiers, "$_"+paramType.TypeIdentifier)
	}

	typeString += strings.Join(typeStrings, ",") + ")"
	typeIdentifier += strings.Join(typeIdentifiers, "$")

	if !strings.HasSuffix(typeIdentifier, "$") {
		typeIdentifier += "$"
	}

	re := regexp.MustCompile(`\${2,}`)
	typeIdentifier = re.ReplaceAllString(typeIdentifier, "$")

	return &TypeDescription{
		TypeString:     typeString,
		TypeIdentifier: typeIdentifier,
	}
}

// FunctionCallOption represents a function call node in the AST.
type FunctionCallOption struct {
	*ASTBuilder

	Id                    int64            `json:"id"`                               // Unique identifier for the node.
	NodeType              ast_pb.NodeType  `json:"node_type"`                        // Type of the node.
	Kind                  ast_pb.NodeType  `json:"kind"`                             // Kind of the node.
	Src                   SrcNode          `json:"src"`                              // Source location of the node.
	Expression            Node[NodeType]   `json:"expression"`                       // Expression of the function call.
	ReferencedDeclaration int64            `json:"referenced_declaration,omitempty"` // Referenced declaration of the function call.
	TypeDescription       *TypeDescription `json:"type_description"`                 // Type description of the function call.
}

// NewFunctionCall creates a new FunctionCallOption node with a given ASTBuilder.
// It initializes the Arguments slice and sets the NodeType and Kind to FUNCTION_CALL.
func NewFunctionCallOption(b *ASTBuilder) *FunctionCallOption {
	return &FunctionCallOption{
		ASTBuilder: b,
		NodeType:   ast_pb.NodeType_FUNCTION_CALL_OPTION,
		Kind:       ast_pb.NodeType_FUNCTION_CALL_OPTION,
	}
}

// SetReferenceDescriptor sets the reference descriptions of the FunctionCallOption node.
func (f *FunctionCallOption) SetReferenceDescriptor(refId int64, refDesc *TypeDescription) bool {
	f.ReferencedDeclaration = refId
	f.TypeDescription = refDesc
	return false
}

// GetId returns the unique identifier of the FunctionCallOption node.
func (f *FunctionCallOption) GetId() int64 {
	return f.Id
}

// GetType returns the type of the FunctionCallOption node.
func (f *FunctionCallOption) GetType() ast_pb.NodeType {
	return f.NodeType
}

// GetSrc returns the source location of the FunctionCallOption node.
func (f *FunctionCallOption) GetSrc() SrcNode {
	return f.Src
}

// GetKind returns the kind of the FunctionCallOption node.
func (f *FunctionCallOption) GetKind() ast_pb.NodeType {
	return f.Kind
}

// GetExpression returns the expression of the FunctionCallOption node.
func (f *FunctionCallOption) GetExpression() Node[NodeType] {
	return f.Expression
}

// GetTypeDescription returns the type description of the FunctionCallOption node.
// Currently, it returns nil and needs to be implemented.
func (f *FunctionCallOption) GetTypeDescription() *TypeDescription {
	return f.TypeDescription
}

// GetNodes returns a slice of nodes that includes the expression of the FunctionCallOption node.
func (f *FunctionCallOption) GetNodes() []Node[NodeType] {
	return []Node[NodeType]{f.Expression}
}

// GetReferenceDeclaration returns the referenced declaration of the FunctionCallOption node.
func (f *FunctionCallOption) GetReferenceDeclaration() int64 {
	return f.ReferencedDeclaration
}

// UnmarshalJSON unmarshals a given JSON byte array into a FunctionCallOption node.
func (f *FunctionCallOption) UnmarshalJSON(data []byte) error {
	var tempMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}

	if id, ok := tempMap["id"]; ok {
		if err := json.Unmarshal(id, &f.Id); err != nil {
			return err
		}
	}

	if nodeType, ok := tempMap["node_type"]; ok {
		if err := json.Unmarshal(nodeType, &f.NodeType); err != nil {
			return err
		}
	}

	if kind, ok := tempMap["kind"]; ok {
		if err := json.Unmarshal(kind, &f.Kind); err != nil {
			return err
		}
	}

	if src, ok := tempMap["src"]; ok {
		if err := json.Unmarshal(src, &f.Src); err != nil {
			return err
		}
	}

	if referencedDeclaration, ok := tempMap["referenced_declaration"]; ok {
		if err := json.Unmarshal(referencedDeclaration, &f.ReferencedDeclaration); err != nil {
			return err
		}
	}

	if typeDescription, ok := tempMap["type_description"]; ok {
		if err := json.Unmarshal(typeDescription, &f.TypeDescription); err != nil {
			return err
		}
	}

	if expression, ok := tempMap["expression"]; ok {
		if err := json.Unmarshal(expression, &f.Expression); err != nil {
			var tempNodeMap map[string]json.RawMessage
			if err := json.Unmarshal(expression, &tempNodeMap); err != nil {
				return err
			}

			var tempNodeType ast_pb.NodeType
			if err := json.Unmarshal(tempNodeMap["node_type"], &tempNodeType); err != nil {
				return err
			}

			node, err := unmarshalNode(expression, tempNodeType)
			if err != nil {
				return err
			}
			f.Expression = node
		}
	}

	return nil
}

// ToProto returns a protobuf representation of the FunctionCallOption node.
// Currently, it returns an empty Statement and needs to be implemented.
func (f *FunctionCallOption) ToProto() NodeType {
	proto := ast_pb.FunctionCallOption{
		Id:                    f.GetId(),
		NodeType:              f.GetType(),
		Kind:                  f.GetKind(),
		Src:                   f.Src.ToProto(),
		ReferencedDeclaration: f.GetReferenceDeclaration(),
		TypeDescription:       f.GetTypeDescription().ToProto(),
	}

	if f.GetExpression() != nil {
		proto.Expression = f.GetExpression().ToProto().(*v3.TypedStruct)
	}

	return NewTypedStruct(&proto, "FunctionCallOption")
}

// Parse takes a parser.FunctionCallOptionsContext and parses it into a FunctionCallOption node.
// It sets the Id, Src, Expression, and TypeDescription of the FunctionCallOption node.
// It returns the created FunctionCallOption node.
func (f *FunctionCallOption) Parse(
	unit *SourceUnit[Node[ast_pb.SourceUnit]],
	contractNode Node[NodeType],
	fnNode Node[NodeType],
	bodyNode *BodyNode,
	vDeclar *VariableDeclaration,
	expNode Node[NodeType],
	ctx *parser.FunctionCallOptionsContext,
) Node[NodeType] {
	f.Id = f.GetNextID()
	f.Src = SrcNode{
		Line:   int64(ctx.GetStart().GetLine()),
		Column: int64(ctx.GetStart().GetColumn()),
		Start:  int64(ctx.GetStart().GetStart()),
		End:    int64(ctx.GetStop().GetStop()),
		Length: int64(ctx.GetStop().GetStop() - ctx.GetStart().GetStart() + 1),
		ParentIndex: func() int64 {
			if vDeclar != nil {
				return vDeclar.GetId()
			}

			if expNode != nil {
				return expNode.GetId()
			}

			if bodyNode != nil {
				return bodyNode.GetId()
			}

			if fnNode != nil {
				return fnNode.GetId()
			}

			return contractNode.GetId()
		}(),
	}

	expression := NewExpression(f.ASTBuilder)

	if ctx.Expression() != nil {
		f.Expression = expression.Parse(
			unit, contractNode, fnNode, bodyNode, nil, f, f.GetId(), ctx.Expression(),
		)
		f.TypeDescription = f.Expression.GetTypeDescription()
	}

	return f
}
