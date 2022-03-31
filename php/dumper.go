package php

import (
	"io"

	"github.com/z7zmey/php-parser/pkg/ast"
)

type PhpDumper struct {
	writer        io.Writer
	indent        int
	withTokens    bool
	withPositions bool
	Expressions   []*Expression
}

type Expression struct {
	Class ast.Vertex
	Call  ast.Vertex
	Args  []ast.Vertex
}

func NewPhpDumper(writer io.Writer) *PhpDumper {
	return &PhpDumper{writer: writer}
}

func (v *PhpDumper) Dump(n ast.Vertex) {
	n.Accept(v)
}

func (v *PhpDumper) dumpVertex(key string, node ast.Vertex) {
	if node == nil {
		return
	}

	node.Accept(v)
}

func (v *PhpDumper) dumpStaticCall(key string, node ast.Vertex) {
	if node == nil {
		return
	}
	node.Accept(v)
}

func (v *PhpDumper) dumpVertexList(key string, list []ast.Vertex) {
	if list == nil {
		return
	}

	if len(list) == 0 {
		return
	}

	for _, nn := range list {
		nn.Accept(v)
	}

}

func (v *PhpDumper) Root(n *ast.Root) {
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *PhpDumper) Nullable(n *ast.Nullable) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) Parameter(n *ast.Parameter) {
	v.dumpVertex("Type", n.Type)
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("DefaultValue", n.DefaultValue)
}

func (v *PhpDumper) Identifier(n *ast.Identifier) {
}

func (v *PhpDumper) Argument(n *ast.Argument) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) StmtBreak(n *ast.StmtBreak) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) StmtCase(n *ast.StmtCase) {
	v.dumpVertex("Cond", n.Cond)
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *PhpDumper) StmtCatch(n *ast.StmtCatch) {
	v.dumpVertexList("Types", n.Types)
	v.dumpVertex("Var", n.Var)
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *PhpDumper) StmtClass(n *ast.StmtClass) {
	v.dumpVertexList("Modifiers", n.Modifiers)
	v.dumpVertex("Name", n.Name)
	v.dumpVertexList("Args", n.Args)
	v.dumpVertex("Extends", n.Extends)
	v.dumpVertexList("Implements", n.Implements)
	v.dumpVertexList("Stmts", n.Stmts)

}

func (v *PhpDumper) StmtClassConstList(n *ast.StmtClassConstList) {
	v.dumpVertexList("Modifiers", n.Modifiers)
	v.dumpVertexList("Consts", n.Consts)
}

func (v *PhpDumper) StmtClassMethod(n *ast.StmtClassMethod) {
	v.dumpVertexList("Modifiers", n.Modifiers)
	v.dumpVertex("Name", n.Name)
	v.dumpVertexList("Params", n.Params)
	v.dumpVertex("ReturnType", n.ReturnType)
	v.dumpVertex("Stmt", n.Stmt)
}

func (v *PhpDumper) StmtConstList(n *ast.StmtConstList) {
	v.dumpVertexList("Consts", n.Consts)
}

func (v *PhpDumper) StmtConstant(n *ast.StmtConstant) {
	v.dumpVertex("Name", n.Name)
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) StmtContinue(n *ast.StmtContinue) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) StmtDeclare(n *ast.StmtDeclare) {
	v.dumpVertexList("Consts", n.Consts)
	v.dumpVertex("Stmt", n.Stmt)
}

func (v *PhpDumper) StmtDefault(n *ast.StmtDefault) {
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *PhpDumper) StmtDo(n *ast.StmtDo) {
	v.dumpVertex("Stmt", n.Stmt)
	v.dumpVertex("Cond", n.Cond)
}

func (v *PhpDumper) StmtEcho(n *ast.StmtEcho) {
	v.dumpVertexList("Exprs", n.Exprs)
}

func (v *PhpDumper) StmtElse(n *ast.StmtElse) {
	v.dumpVertex("Stmt", n.Stmt)
}

func (v *PhpDumper) StmtElseIf(n *ast.StmtElseIf) {
	v.dumpVertex("Cond", n.Cond)
	v.dumpVertex("Stmt", n.Stmt)
}

func (v *PhpDumper) StmtExpression(n *ast.StmtExpression) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) StmtFinally(n *ast.StmtFinally) {
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *PhpDumper) StmtFor(n *ast.StmtFor) {
	v.dumpVertexList("Init", n.Init)
	v.dumpVertexList("Cond", n.Cond)
	v.dumpVertexList("Loop", n.Loop)
	v.dumpVertex("Stmt", n.Stmt)
}

func (v *PhpDumper) StmtForeach(n *ast.StmtForeach) {
	v.dumpVertex("Expr", n.Expr)
	v.dumpVertex("Key", n.Key)
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Stmt", n.Stmt)
}

func (v *PhpDumper) StmtFunction(n *ast.StmtFunction) {
	v.dumpVertex("Name", n.Name)
	v.dumpVertexList("Params", n.Params)
	v.dumpVertex("ReturnType", n.ReturnType)
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *PhpDumper) StmtGlobal(n *ast.StmtGlobal) {
	v.dumpVertexList("Vars", n.Vars)
}

func (v *PhpDumper) StmtGoto(n *ast.StmtGoto) {
	v.dumpVertex("Label", n.Label)
}

func (v *PhpDumper) StmtHaltCompiler(n *ast.StmtHaltCompiler) {

}

func (v *PhpDumper) StmtIf(n *ast.StmtIf) {
	v.dumpVertex("Cond", n.Cond)
	v.dumpVertex("Stmt", n.Stmt)
	v.dumpVertexList("ElseIf", n.ElseIf)
	v.dumpVertex("Else", n.Else)
}

func (v *PhpDumper) StmtInlineHtml(n *ast.StmtInlineHtml) {

}

func (v *PhpDumper) StmtInterface(n *ast.StmtInterface) {
	v.dumpVertex("Name", n.Name)
	v.dumpVertexList("Extends", n.Extends)
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *PhpDumper) StmtLabel(n *ast.StmtLabel) {
	v.dumpVertex("Name", n.Name)
}

func (v *PhpDumper) StmtNamespace(n *ast.StmtNamespace) {
	v.dumpVertex("Name", n.Name)
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *PhpDumper) StmtNop(n *ast.StmtNop) {

}

func (v *PhpDumper) StmtProperty(n *ast.StmtProperty) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)

}

func (v *PhpDumper) StmtPropertyList(n *ast.StmtPropertyList) {
	v.dumpVertexList("Modifiers", n.Modifiers)
	v.dumpVertex("Type", n.Type)
	v.dumpVertexList("Props", n.Props)
}

func (v *PhpDumper) StmtReturn(n *ast.StmtReturn) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) StmtStatic(n *ast.StmtStatic) {
	v.dumpVertexList("Vars", n.Vars)
}

func (v *PhpDumper) StmtStaticVar(n *ast.StmtStaticVar) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) StmtStmtList(n *ast.StmtStmtList) {
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *PhpDumper) StmtSwitch(n *ast.StmtSwitch) {
	v.dumpVertex("Cond", n.Cond)
	v.dumpVertexList("Cases", n.Cases)
}

func (v *PhpDumper) StmtThrow(n *ast.StmtThrow) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) StmtTrait(n *ast.StmtTrait) {
	v.dumpVertex("Name", n.Name)
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *PhpDumper) StmtTraitUse(n *ast.StmtTraitUse) {
	v.dumpVertexList("Traits", n.Traits)
	v.dumpVertexList("Adaptations", n.Adaptations)
}

func (v *PhpDumper) StmtTraitUseAlias(n *ast.StmtTraitUseAlias) {
	v.dumpVertex("Trait", n.Trait)
	v.dumpVertex("Method", n.Method)
	v.dumpVertex("Modifier", n.Modifier)
	v.dumpVertex("Alias", n.Alias)

}

func (v *PhpDumper) StmtTraitUsePrecedence(n *ast.StmtTraitUsePrecedence) {
	v.dumpVertex("Trait", n.Trait)
	v.dumpVertex("Method", n.Method)
	v.dumpVertexList("Insteadof", n.Insteadof)
}

func (v *PhpDumper) StmtTry(n *ast.StmtTry) {
	v.dumpVertexList("Stmts", n.Stmts)
	v.dumpVertexList("Catches", n.Catches)
	v.dumpVertex("Finally", n.Finally)
}

func (v *PhpDumper) StmtUnset(n *ast.StmtUnset) {
	v.dumpVertexList("Vars", n.Vars)
}

func (v *PhpDumper) StmtUse(n *ast.StmtUseList) {
	v.dumpVertex("Type", n.Type)
	v.dumpVertexList("Uses", n.Uses)
}

func (v *PhpDumper) StmtGroupUse(n *ast.StmtGroupUseList) {
	v.dumpVertex("Type", n.Type)
	v.dumpVertex("Prefix", n.Prefix)
	v.dumpVertexList("Uses", n.Uses)
}

func (v *PhpDumper) StmtUseDeclaration(n *ast.StmtUse) {
	v.dumpVertex("Type", n.Type)
	v.dumpVertex("Uses", n.Use)
	v.dumpVertex("Alias", n.Alias)
}

func (v *PhpDumper) StmtWhile(n *ast.StmtWhile) {
	v.dumpVertex("Cond", n.Cond)
	v.dumpVertex("Stmt", n.Stmt)
}

func (v *PhpDumper) ExprArray(n *ast.ExprArray) {
	v.dumpVertexList("Items", n.Items)
}

func (v *PhpDumper) ExprArrayDimFetch(n *ast.ExprArrayDimFetch) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Dim", n.Dim)
}

func (v *PhpDumper) ExprArrayItem(n *ast.ExprArrayItem) {
	v.dumpVertex("Key", n.Key)
	v.dumpVertex("Val", n.Val)
}

func (v *PhpDumper) ExprArrowFunction(n *ast.ExprArrowFunction) {
	v.dumpVertexList("Params", n.Params)
	v.dumpVertex("ReturnType", n.ReturnType)
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprBitwiseNot(n *ast.ExprBitwiseNot) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprBooleanNot(n *ast.ExprBooleanNot) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprBrackets(n *ast.ExprBrackets) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprClassConstFetch(n *ast.ExprClassConstFetch) {

	v.dumpVertex("Class", n.Class)

	v.dumpVertex("Const", n.Const)

}

func (v *PhpDumper) ExprClone(n *ast.ExprClone) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprClosure(n *ast.ExprClosure) {
	v.dumpVertexList("Params", n.Params)
	v.dumpVertexList("Uses", n.Uses)
	v.dumpVertex("ReturnType", n.ReturnType)
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *PhpDumper) ExprClosureUse(n *ast.ExprClosureUse) {
	v.dumpVertex("Var", n.Var)
}

func (v *PhpDumper) ExprConstFetch(n *ast.ExprConstFetch) {
	v.dumpVertex("Const", n.Const)
}

func (v *PhpDumper) ExprEmpty(n *ast.ExprEmpty) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprErrorSuppress(n *ast.ExprErrorSuppress) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprEval(n *ast.ExprEval) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprExit(n *ast.ExprExit) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprFunctionCall(n *ast.ExprFunctionCall) {
	v.dumpVertex("Function", n.Function)
	v.dumpVertexList("Args", n.Args)
}

func (v *PhpDumper) ExprInclude(n *ast.ExprInclude) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprIncludeOnce(n *ast.ExprIncludeOnce) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprInstanceOf(n *ast.ExprInstanceOf) {
	v.dumpVertex("Expr", n.Expr)
	v.dumpVertex("Class", n.Class)
}

func (v *PhpDumper) ExprIsset(n *ast.ExprIsset) {
	v.dumpVertexList("Vars", n.Vars)
}

func (v *PhpDumper) ExprList(n *ast.ExprList) {
	v.dumpVertexList("Items", n.Items)
}

func (v *PhpDumper) ExprMethodCall(n *ast.ExprMethodCall) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Method", n.Method)
	v.dumpVertexList("Args", n.Args)
}

func (v *PhpDumper) ExprNew(n *ast.ExprNew) {
	v.dumpVertex("Class", n.Class)
	v.dumpVertexList("Args", n.Args)
}

func (v *PhpDumper) ExprPostDec(n *ast.ExprPostDec) {
	v.dumpVertex("Var", n.Var)
}

func (v *PhpDumper) ExprPostInc(n *ast.ExprPostInc) {
	v.dumpVertex("Var", n.Var)
}

func (v *PhpDumper) ExprPreDec(n *ast.ExprPreDec) {
	v.dumpVertex("Var", n.Var)
}

func (v *PhpDumper) ExprPreInc(n *ast.ExprPreInc) {
	v.dumpVertex("Var", n.Var)
}

func (v *PhpDumper) ExprPrint(n *ast.ExprPrint) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprPropertyFetch(n *ast.ExprPropertyFetch) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Prop", n.Prop)
}

func (v *PhpDumper) ExprRequire(n *ast.ExprRequire) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprRequireOnce(n *ast.ExprRequireOnce) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprShellExec(n *ast.ExprShellExec) {
	v.dumpVertexList("Parts", n.Parts)
}

func (v *PhpDumper) ExprStaticCall(n *ast.ExprStaticCall) {
	// Drupal::service('foo');
	// Drupal -> ast.Class
	expr := &Expression{
		Class: n.Class,
		Call:  n.Call,
		Args:  n.Args,
	}
	v.Expressions = append(v.Expressions, expr)
}

func (v *PhpDumper) ExprStaticPropertyFetch(n *ast.ExprStaticPropertyFetch) {
	v.dumpVertex("Class", n.Class)
	v.dumpVertex("Prop", n.Prop)
}

func (v *PhpDumper) ExprTernary(n *ast.ExprTernary) {
	v.dumpVertex("Cond", n.Cond)
	v.dumpVertex("IfTrue", n.IfTrue)
	v.dumpVertex("IfFalse", n.IfFalse)
}

func (v *PhpDumper) ExprUnaryMinus(n *ast.ExprUnaryMinus) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprUnaryPlus(n *ast.ExprUnaryPlus) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprVariable(n *ast.ExprVariable) {
	v.dumpVertex("Name", n.Name)
}

func (v *PhpDumper) ExprYield(n *ast.ExprYield) {
	v.dumpVertex("Key", n.Key)
	v.dumpVertex("Val", n.Val)
}

func (v *PhpDumper) ExprYieldFrom(n *ast.ExprYieldFrom) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprAssign(n *ast.ExprAssign) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprAssignReference(n *ast.ExprAssignReference) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprAssignBitwiseAnd(n *ast.ExprAssignBitwiseAnd) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprAssignBitwiseOr(n *ast.ExprAssignBitwiseOr) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprAssignBitwiseXor(n *ast.ExprAssignBitwiseXor) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprAssignCoalesce(n *ast.ExprAssignCoalesce) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprAssignConcat(n *ast.ExprAssignConcat) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprAssignDiv(n *ast.ExprAssignDiv) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprAssignMinus(n *ast.ExprAssignMinus) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprAssignMod(n *ast.ExprAssignMod) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprAssignMul(n *ast.ExprAssignMul) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprAssignPlus(n *ast.ExprAssignPlus) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)

}

func (v *PhpDumper) ExprAssignPow(n *ast.ExprAssignPow) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprAssignShiftLeft(n *ast.ExprAssignShiftLeft) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprAssignShiftRight(n *ast.ExprAssignShiftRight) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprBinaryBitwiseAnd(n *ast.ExprBinaryBitwiseAnd) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryBitwiseOr(n *ast.ExprBinaryBitwiseOr) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryBitwiseXor(n *ast.ExprBinaryBitwiseXor) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryBooleanAnd(n *ast.ExprBinaryBooleanAnd) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryBooleanOr(n *ast.ExprBinaryBooleanOr) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryCoalesce(n *ast.ExprBinaryCoalesce) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryConcat(n *ast.ExprBinaryConcat) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryDiv(n *ast.ExprBinaryDiv) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryEqual(n *ast.ExprBinaryEqual) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryGreater(n *ast.ExprBinaryGreater) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryGreaterOrEqual(n *ast.ExprBinaryGreaterOrEqual) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryIdentical(n *ast.ExprBinaryIdentical) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryLogicalAnd(n *ast.ExprBinaryLogicalAnd) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryLogicalOr(n *ast.ExprBinaryLogicalOr) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryLogicalXor(n *ast.ExprBinaryLogicalXor) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryMinus(n *ast.ExprBinaryMinus) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryMod(n *ast.ExprBinaryMod) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryMul(n *ast.ExprBinaryMul) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryNotEqual(n *ast.ExprBinaryNotEqual) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryNotIdentical(n *ast.ExprBinaryNotIdentical) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryPlus(n *ast.ExprBinaryPlus) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryPow(n *ast.ExprBinaryPow) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryShiftLeft(n *ast.ExprBinaryShiftLeft) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinaryShiftRight(n *ast.ExprBinaryShiftRight) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinarySmaller(n *ast.ExprBinarySmaller) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinarySmallerOrEqual(n *ast.ExprBinarySmallerOrEqual) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprBinarySpaceship(n *ast.ExprBinarySpaceship) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *PhpDumper) ExprCastArray(n *ast.ExprCastArray) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprCastBool(n *ast.ExprCastBool) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprCastDouble(n *ast.ExprCastDouble) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprCastInt(n *ast.ExprCastInt) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprCastObject(n *ast.ExprCastObject) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ExprCastString(n *ast.ExprCastString) {
	v.dumpVertex("Expr", n.Expr)

}

func (v *PhpDumper) ExprCastUnset(n *ast.ExprCastUnset) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *PhpDumper) ScalarDnumber(n *ast.ScalarDnumber) {

}

func (v *PhpDumper) ScalarEncapsed(n *ast.ScalarEncapsed) {
	v.dumpVertexList("Parts", n.Parts)
}

func (v *PhpDumper) ScalarEncapsedStringPart(n *ast.ScalarEncapsedStringPart) {

}

func (v *PhpDumper) ScalarEncapsedStringVar(n *ast.ScalarEncapsedStringVar) {
	v.dumpVertex("Name", n.Name)
	v.dumpVertex("Dim", n.Dim)
}

func (v *PhpDumper) ScalarEncapsedStringBrackets(n *ast.ScalarEncapsedStringBrackets) {
	v.dumpVertex("Var", n.Var)
}

func (v *PhpDumper) ScalarHeredoc(n *ast.ScalarHeredoc) {
	v.dumpVertexList("Parts", n.Parts)
}

func (v *PhpDumper) ScalarLnumber(n *ast.ScalarLnumber) {

}

func (v *PhpDumper) ScalarMagicConstant(n *ast.ScalarMagicConstant) {

}

func (v *PhpDumper) ScalarString(n *ast.ScalarString) {

}

func (v *PhpDumper) NameName(n *ast.Name) {
	v.dumpVertexList("Parts", n.Parts)
}

func (v *PhpDumper) NameFullyQualified(n *ast.NameFullyQualified) {
	v.dumpVertexList("Parts", n.Parts)
}

func (v *PhpDumper) NameRelative(n *ast.NameRelative) {
	v.dumpVertexList("Parts", n.Parts)
}

func (v *PhpDumper) NameNamePart(n *ast.NamePart) {

}
