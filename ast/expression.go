package ast

import (
	"fmt"
	ast_pb "github.com/unpackdev/protos/dist/go/ast"
	"github.com/unpackdev/solgo/parser"
	"go.uber.org/zap"
)

// Expression represents an AST node for an expression in Solidity.
type Expression struct {
	*ASTBuilder
}

// NewExpression creates a new Expression instance with the provided ASTBuilder.
// The ASTBuilder is used to facilitate the construction of the AST.
func NewExpression(b *ASTBuilder) *Expression {
	return &Expression{
		ASTBuilder: b,
	}
}

// Parse analyzes the provided parser.IExpressionContext and constructs the
// corresponding AST node. It supports various types of expressions in Solidity
// such as binary operations, assignments, function calls, member accesses, etc.
// If the expression type is not supported, a warning is logged.
//
// Parameters:
// - unit: The source unit node.
// - contractNode: The contract node within the source.
// - fnNode: The function node within the contract.
// - bodyNode: The body node of the function.
// - vDecar: The variable declaration node.
// - parentNode: The parent expression node.
// - ctx: The context representing the expression to be parsed.
//
// Returns:
//   - Node[NodeType]: The constructed AST node for the parsed expression. If the
//     expression type is not supported, it returns nil.
func (e *Expression) Parse(
	unit *SourceUnit[Node[ast_pb.SourceUnit]],
	contractNode Node[NodeType],
	fnNode Node[NodeType],
	bodyNode *BodyNode,
	vDecar *VariableDeclaration,
	parentNode Node[NodeType],
	parentNodeId int64,
	ctx parser.IExpressionContext,
) Node[NodeType] {
	switch childCtx := ctx.(type) {
	case *parser.AddSubOperationContext:
		binaryExp := NewBinaryOperationExpression(e.ASTBuilder)
		return binaryExp.ParseAddSub(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, parentNodeId, childCtx)
	case *parser.OrderComparisonContext:
		binaryExp := NewBinaryOperationExpression(e.ASTBuilder)
		return binaryExp.ParseOrderComparison(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, parentNodeId, childCtx)
	case *parser.MulDivModOperationContext:
		binaryExp := NewBinaryOperationExpression(e.ASTBuilder)
		return binaryExp.ParseMulDivMod(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, parentNodeId, childCtx)
	case *parser.EqualityComparisonContext:
		binaryExp := NewBinaryOperationExpression(e.ASTBuilder)
		return binaryExp.ParseEqualityComparison(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, parentNodeId, childCtx)
	case *parser.OrOperationContext:
		binaryExp := NewBinaryOperationExpression(e.ASTBuilder)
		return binaryExp.ParseOr(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, parentNodeId, childCtx)
	case *parser.AssignmentContext:
		assignment := NewAssignment(e.ASTBuilder)
		return assignment.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, parentNodeId, childCtx)
	case *parser.FunctionCallContext:
		statementNode := NewFunctionCall(e.ASTBuilder)
		return statementNode.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, childCtx)
	case *parser.MemberAccessContext:
		memberAccess := NewMemberAccessExpression(e.ASTBuilder)
		return memberAccess.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, childCtx)
	case *parser.PrimaryExpressionContext:
		primaryExp := NewPrimaryExpression(e.ASTBuilder)
		return primaryExp.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, parentNodeId, childCtx)
	case *parser.IndexAccessContext:
		indexAccess := NewIndexAccess(e.ASTBuilder)
		return indexAccess.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, childCtx)
	case *parser.MetaTypeContext:
		metaType := NewMetaTypeExpression(e.ASTBuilder)
		return metaType.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, childCtx)
	case *parser.PayableConversionContext:
		payableConversion := NewPayableConversionExpression(e.ASTBuilder)
		return payableConversion.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, childCtx)
	case *parser.UnarySuffixOperationContext:
		unarySuffixOperation := NewUnarySuffixExpression(e.ASTBuilder)
		return unarySuffixOperation.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, childCtx)
	case *parser.UnaryPrefixOperationContext:
		unaryPrefixOperation := NewUnaryPrefixExpression(e.ASTBuilder)
		return unaryPrefixOperation.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, childCtx)
	case *parser.NewExprContext:
		newExpr := NewExprExpression(e.ASTBuilder)
		return newExpr.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, childCtx)
	case *parser.TupleContext:
		tupleExpr := NewTupleExpression(e.ASTBuilder)
		return tupleExpr.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, childCtx)
	case *parser.FunctionCallOptionsContext:
		statementNode := NewFunctionCallOption(e.ASTBuilder)
		return statementNode.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, childCtx)
	case *parser.IndexRangeAccessContext:
		indexRangeAccess := NewIndexRangeAccessExpression(e.ASTBuilder)
		return indexRangeAccess.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, childCtx)
	case *parser.ExpOperationContext:
		expOperation := NewExprOperationExpression(e.ASTBuilder)
		return expOperation.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, childCtx)
	case *parser.ConditionalContext:
		conditional := NewConditionalExpression(e.ASTBuilder)
		return conditional.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, childCtx)
	case *parser.AndOperationContext:
		andOperation := NewAndOperationExpression(e.ASTBuilder)
		return andOperation.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, parentNodeId, childCtx)
	case *parser.BitAndOperationContext:
		bitAndOperation := NewBitAndOperationExpression(e.ASTBuilder)
		return bitAndOperation.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, parentNodeId, childCtx)
	case *parser.BitOrOperationContext:
		bitAndOperation := NewBitOrOperationExpression(e.ASTBuilder)
		return bitAndOperation.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, parentNodeId, childCtx)
	case *parser.BitXorOperationContext:
		bitXorOperation := NewBitXorOperationExpression(e.ASTBuilder)
		return bitXorOperation.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, parentNodeId, childCtx)
	case *parser.ShiftOperationContext:
		shiftOperation := NewShiftOperationExpression(e.ASTBuilder)
		return shiftOperation.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, parentNodeId, childCtx)
	case *parser.InlineArrayContext:
		inlineArray := NewInlineArrayExpression(e.ASTBuilder)
		return inlineArray.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, childCtx)
	case *parser.ExpressionContext:
		expressionCtx := NewExpressionContext(e.ASTBuilder)
		return expressionCtx.Parse(unit, contractNode, fnNode, bodyNode, vDecar, parentNode, parentNodeId, childCtx)
	default:
		zap.L().Warn(
			"Expression type not supported @ Expression.Parse",
			zap.String("type", fmt.Sprintf("%T", ctx)),
		)
	}

	return nil
}

func (e *Expression) ParseInterface(
	unit *SourceUnit[Node[ast_pb.SourceUnit]],
	fnNode Node[NodeType],
	parentNodeId int64,
	ctx any,
) Node[NodeType] {
	switch childCtx := ctx.(type) {
	case *parser.AddSubOperationContext:
		binaryExp := NewBinaryOperationExpression(e.ASTBuilder)
		return binaryExp.ParseAddSub(unit, nil, fnNode, nil, nil, nil, parentNodeId, childCtx)
	case *parser.OrderComparisonContext:
		binaryExp := NewBinaryOperationExpression(e.ASTBuilder)
		return binaryExp.ParseOrderComparison(unit, nil, fnNode, nil, nil, nil, parentNodeId, childCtx)
	case *parser.MulDivModOperationContext:
		binaryExp := NewBinaryOperationExpression(e.ASTBuilder)
		return binaryExp.ParseMulDivMod(unit, nil, fnNode, nil, nil, nil, parentNodeId, childCtx)
	case *parser.EqualityComparisonContext:
		binaryExp := NewBinaryOperationExpression(e.ASTBuilder)
		return binaryExp.ParseEqualityComparison(unit, nil, fnNode, nil, nil, nil, parentNodeId, childCtx)
	case *parser.OrOperationContext:
		binaryExp := NewBinaryOperationExpression(e.ASTBuilder)
		return binaryExp.ParseOr(unit, nil, fnNode, nil, nil, nil, parentNodeId, childCtx)
	case *parser.AssignmentContext:
		assignment := NewAssignment(e.ASTBuilder)
		return assignment.Parse(unit, nil, fnNode, nil, nil, nil, parentNodeId, childCtx)
	case *parser.FunctionCallContext:
		statementNode := NewFunctionCall(e.ASTBuilder)
		return statementNode.Parse(unit, nil, fnNode, nil, nil, nil, childCtx)
	case *parser.MemberAccessContext:
		memberAccess := NewMemberAccessExpression(e.ASTBuilder)
		return memberAccess.Parse(unit, nil, fnNode, nil, nil, nil, childCtx)
	case *parser.PrimaryExpressionContext:
		primaryExp := NewPrimaryExpression(e.ASTBuilder)
		return primaryExp.Parse(unit, nil, fnNode, nil, nil, nil, parentNodeId, childCtx)
	case *parser.IndexAccessContext:
		indexAccess := NewIndexAccess(e.ASTBuilder)
		return indexAccess.Parse(unit, nil, fnNode, nil, nil, nil, childCtx)
	case *parser.MetaTypeContext:
		metaType := NewMetaTypeExpression(e.ASTBuilder)
		return metaType.Parse(unit, nil, fnNode, nil, nil, nil, childCtx)
	case *parser.PayableConversionContext:
		payableConversion := NewPayableConversionExpression(e.ASTBuilder)
		return payableConversion.Parse(unit, nil, fnNode, nil, nil, nil, childCtx)
	case *parser.UnarySuffixOperationContext:
		unarySuffixOperation := NewUnarySuffixExpression(e.ASTBuilder)
		return unarySuffixOperation.Parse(unit, nil, fnNode, nil, nil, nil, childCtx)
	case *parser.UnaryPrefixOperationContext:
		unaryPrefixOperation := NewUnaryPrefixExpression(e.ASTBuilder)
		return unaryPrefixOperation.Parse(unit, nil, fnNode, nil, nil, nil, childCtx)
	case *parser.NewExprContext:
		newExpr := NewExprExpression(e.ASTBuilder)
		return newExpr.Parse(unit, nil, fnNode, nil, nil, nil, childCtx)
	case *parser.TupleContext:
		tupleExpr := NewTupleExpression(e.ASTBuilder)
		return tupleExpr.Parse(unit, nil, fnNode, nil, nil, nil, childCtx)
	case *parser.FunctionCallOptionsContext:
		statementNode := NewFunctionCallOption(e.ASTBuilder)
		return statementNode.Parse(unit, nil, fnNode, nil, nil, nil, childCtx)
	case *parser.IndexRangeAccessContext:
		indexRangeAccess := NewIndexRangeAccessExpression(e.ASTBuilder)
		return indexRangeAccess.Parse(unit, nil, fnNode, nil, nil, nil, childCtx)
	case *parser.ExpOperationContext:
		expOperation := NewExprOperationExpression(e.ASTBuilder)
		return expOperation.Parse(unit, nil, fnNode, nil, nil, nil, childCtx)
	case *parser.ConditionalContext:
		conditional := NewConditionalExpression(e.ASTBuilder)
		return conditional.Parse(unit, nil, fnNode, nil, nil, nil, childCtx)
	case *parser.AndOperationContext:
		andOperation := NewAndOperationExpression(e.ASTBuilder)
		return andOperation.Parse(unit, nil, fnNode, nil, nil, nil, parentNodeId, childCtx)
	case *parser.BitAndOperationContext:
		bitAndOperation := NewBitAndOperationExpression(e.ASTBuilder)
		return bitAndOperation.Parse(unit, nil, fnNode, nil, nil, nil, parentNodeId, childCtx)
	case *parser.BitOrOperationContext:
		bitAndOperation := NewBitOrOperationExpression(e.ASTBuilder)
		return bitAndOperation.Parse(unit, nil, fnNode, nil, nil, nil, parentNodeId, childCtx)
	case *parser.BitXorOperationContext:
		bitXorOperation := NewBitXorOperationExpression(e.ASTBuilder)
		return bitXorOperation.Parse(unit, nil, fnNode, nil, nil, nil, parentNodeId, childCtx)
	case *parser.ShiftOperationContext:
		shiftOperation := NewShiftOperationExpression(e.ASTBuilder)
		return shiftOperation.Parse(unit, nil, fnNode, nil, nil, nil, parentNodeId, childCtx)
	case *parser.InlineArrayContext:
		inlineArray := NewInlineArrayExpression(e.ASTBuilder)
		return inlineArray.Parse(unit, nil, fnNode, nil, nil, nil, childCtx)
	default:
		zap.L().Warn(
			"Expression type not supported @ Expression.Parse",
			zap.String("type", fmt.Sprintf("%T", ctx)),
		)
	}

	return nil
}
