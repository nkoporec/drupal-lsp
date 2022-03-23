package parser

import (
	"io"

	"github.com/z7zmey/php-parser/pkg/ast"
)

type Dumper struct {
	writer        io.Writer
	indent        int
	withTokens    bool
	withPositions bool
	StaticCalls   []*ExprStaticCall
}

type ExprStaticCall struct {
	StartLine int
	EndLine   int
	StartPos  int
	EndPos    int
}

func NewServiceDumper(writer io.Writer) *Dumper {
	return &Dumper{writer: writer}
}

func (v *Dumper) WithTokens() *Dumper {
	v.withTokens = true
	return v
}

func (v *Dumper) WithPositions() *Dumper {
	v.withPositions = true
	return v
}

func (v *Dumper) Dump(n ast.Vertex) {
	n.Accept(v)
}

func (v *Dumper) dumpVertex(key string, node ast.Vertex) {
	if node == nil {
		return
	}

	node.Accept(v)
}

func (v *Dumper) dumpVertexList(key string, list []ast.Vertex) {
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

func (v *Dumper) Root(n *ast.Root) {
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *Dumper) Nullable(n *ast.Nullable) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) Parameter(n *ast.Parameter) {
	v.dumpVertex("Type", n.Type)
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("DefaultValue", n.DefaultValue)
}

func (v *Dumper) Identifier(n *ast.Identifier) {
}

func (v *Dumper) Argument(n *ast.Argument) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) StmtBreak(n *ast.StmtBreak) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) StmtCase(n *ast.StmtCase) {
	v.dumpVertex("Cond", n.Cond)
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *Dumper) StmtCatch(n *ast.StmtCatch) {
	v.dumpVertexList("Types", n.Types)
	v.dumpVertex("Var", n.Var)
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *Dumper) StmtClass(n *ast.StmtClass) {
	v.dumpVertexList("Modifiers", n.Modifiers)
	v.dumpVertex("Name", n.Name)
	v.dumpVertexList("Args", n.Args)
	v.dumpVertex("Extends", n.Extends)
	v.dumpVertexList("Implements", n.Implements)
	v.dumpVertexList("Stmts", n.Stmts)

}

func (v *Dumper) StmtClassConstList(n *ast.StmtClassConstList) {
	v.dumpVertexList("Modifiers", n.Modifiers)
	v.dumpVertexList("Consts", n.Consts)
}

func (v *Dumper) StmtClassMethod(n *ast.StmtClassMethod) {
	v.dumpVertexList("Modifiers", n.Modifiers)
	v.dumpVertex("Name", n.Name)
	v.dumpVertexList("Params", n.Params)
	v.dumpVertex("ReturnType", n.ReturnType)
	v.dumpVertex("Stmt", n.Stmt)
}

func (v *Dumper) StmtConstList(n *ast.StmtConstList) {
	v.dumpVertexList("Consts", n.Consts)
}

func (v *Dumper) StmtConstant(n *ast.StmtConstant) {
	v.dumpVertex("Name", n.Name)
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) StmtContinue(n *ast.StmtContinue) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) StmtDeclare(n *ast.StmtDeclare) {
	v.dumpVertexList("Consts", n.Consts)
	v.dumpVertex("Stmt", n.Stmt)
}

func (v *Dumper) StmtDefault(n *ast.StmtDefault) {
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *Dumper) StmtDo(n *ast.StmtDo) {
	v.dumpVertex("Stmt", n.Stmt)
	v.dumpVertex("Cond", n.Cond)
}

func (v *Dumper) StmtEcho(n *ast.StmtEcho) {
	v.dumpVertexList("Exprs", n.Exprs)
}

func (v *Dumper) StmtElse(n *ast.StmtElse) {
	v.dumpVertex("Stmt", n.Stmt)
}

func (v *Dumper) StmtElseIf(n *ast.StmtElseIf) {
	v.dumpVertex("Cond", n.Cond)
	v.dumpVertex("Stmt", n.Stmt)
}

func (v *Dumper) StmtExpression(n *ast.StmtExpression) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) StmtFinally(n *ast.StmtFinally) {
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *Dumper) StmtFor(n *ast.StmtFor) {
	v.dumpVertexList("Init", n.Init)
	v.dumpVertexList("Cond", n.Cond)
	v.dumpVertexList("Loop", n.Loop)
	v.dumpVertex("Stmt", n.Stmt)
}

func (v *Dumper) StmtForeach(n *ast.StmtForeach) {
	v.dumpVertex("Expr", n.Expr)
	v.dumpVertex("Key", n.Key)
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Stmt", n.Stmt)
}

func (v *Dumper) StmtFunction(n *ast.StmtFunction) {
	v.dumpVertex("Name", n.Name)
	v.dumpVertexList("Params", n.Params)
	v.dumpVertex("ReturnType", n.ReturnType)
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *Dumper) StmtGlobal(n *ast.StmtGlobal) {
	v.dumpVertexList("Vars", n.Vars)
}

func (v *Dumper) StmtGoto(n *ast.StmtGoto) {
	v.dumpVertex("Label", n.Label)
}

func (v *Dumper) StmtHaltCompiler(n *ast.StmtHaltCompiler) {

}

func (v *Dumper) StmtIf(n *ast.StmtIf) {
	v.dumpVertex("Cond", n.Cond)
	v.dumpVertex("Stmt", n.Stmt)
	v.dumpVertexList("ElseIf", n.ElseIf)
	v.dumpVertex("Else", n.Else)
}

func (v *Dumper) StmtInlineHtml(n *ast.StmtInlineHtml) {

}

func (v *Dumper) StmtInterface(n *ast.StmtInterface) {
	v.dumpVertex("Name", n.Name)
	v.dumpVertexList("Extends", n.Extends)
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *Dumper) StmtLabel(n *ast.StmtLabel) {
	v.dumpVertex("Name", n.Name)
}

func (v *Dumper) StmtNamespace(n *ast.StmtNamespace) {
	v.dumpVertex("Name", n.Name)
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *Dumper) StmtNop(n *ast.StmtNop) {

}

func (v *Dumper) StmtProperty(n *ast.StmtProperty) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)

}

func (v *Dumper) StmtPropertyList(n *ast.StmtPropertyList) {
	v.dumpVertexList("Modifiers", n.Modifiers)
	v.dumpVertex("Type", n.Type)
	v.dumpVertexList("Props", n.Props)
}

func (v *Dumper) StmtReturn(n *ast.StmtReturn) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) StmtStatic(n *ast.StmtStatic) {
	v.dumpVertexList("Vars", n.Vars)
}

func (v *Dumper) StmtStaticVar(n *ast.StmtStaticVar) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) StmtStmtList(n *ast.StmtStmtList) {
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *Dumper) StmtSwitch(n *ast.StmtSwitch) {
	v.dumpVertex("Cond", n.Cond)
	v.dumpVertexList("Cases", n.Cases)
}

func (v *Dumper) StmtThrow(n *ast.StmtThrow) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) StmtTrait(n *ast.StmtTrait) {
	v.dumpVertex("Name", n.Name)
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *Dumper) StmtTraitUse(n *ast.StmtTraitUse) {
	v.dumpVertexList("Traits", n.Traits)
	v.dumpVertexList("Adaptations", n.Adaptations)
}

func (v *Dumper) StmtTraitUseAlias(n *ast.StmtTraitUseAlias) {
	v.dumpVertex("Trait", n.Trait)
	v.dumpVertex("Method", n.Method)
	v.dumpVertex("Modifier", n.Modifier)
	v.dumpVertex("Alias", n.Alias)

}

func (v *Dumper) StmtTraitUsePrecedence(n *ast.StmtTraitUsePrecedence) {
	v.dumpVertex("Trait", n.Trait)
	v.dumpVertex("Method", n.Method)
	v.dumpVertexList("Insteadof", n.Insteadof)
}

func (v *Dumper) StmtTry(n *ast.StmtTry) {
	v.dumpVertexList("Stmts", n.Stmts)
	v.dumpVertexList("Catches", n.Catches)
	v.dumpVertex("Finally", n.Finally)
}

func (v *Dumper) StmtUnset(n *ast.StmtUnset) {
	v.dumpVertexList("Vars", n.Vars)
}

func (v *Dumper) StmtUse(n *ast.StmtUseList) {
	v.dumpVertex("Type", n.Type)
	v.dumpVertexList("Uses", n.Uses)
}

func (v *Dumper) StmtGroupUse(n *ast.StmtGroupUseList) {
	v.dumpVertex("Type", n.Type)
	v.dumpVertex("Prefix", n.Prefix)
	v.dumpVertexList("Uses", n.Uses)
}

func (v *Dumper) StmtUseDeclaration(n *ast.StmtUse) {
	v.dumpVertex("Type", n.Type)
	v.dumpVertex("Uses", n.Use)
	v.dumpVertex("Alias", n.Alias)
}

func (v *Dumper) StmtWhile(n *ast.StmtWhile) {
	v.dumpVertex("Cond", n.Cond)
	v.dumpVertex("Stmt", n.Stmt)
}

func (v *Dumper) ExprArray(n *ast.ExprArray) {
	v.dumpVertexList("Items", n.Items)
}

func (v *Dumper) ExprArrayDimFetch(n *ast.ExprArrayDimFetch) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Dim", n.Dim)
}

func (v *Dumper) ExprArrayItem(n *ast.ExprArrayItem) {
	v.dumpVertex("Key", n.Key)
	v.dumpVertex("Val", n.Val)
}

func (v *Dumper) ExprArrowFunction(n *ast.ExprArrowFunction) {
	v.dumpVertexList("Params", n.Params)
	v.dumpVertex("ReturnType", n.ReturnType)
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprBitwiseNot(n *ast.ExprBitwiseNot) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprBooleanNot(n *ast.ExprBooleanNot) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprBrackets(n *ast.ExprBrackets) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprClassConstFetch(n *ast.ExprClassConstFetch) {

	v.dumpVertex("Class", n.Class)

	v.dumpVertex("Const", n.Const)

}

func (v *Dumper) ExprClone(n *ast.ExprClone) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprClosure(n *ast.ExprClosure) {
	v.dumpVertexList("Params", n.Params)
	v.dumpVertexList("Uses", n.Uses)
	v.dumpVertex("ReturnType", n.ReturnType)
	v.dumpVertexList("Stmts", n.Stmts)
}

func (v *Dumper) ExprClosureUse(n *ast.ExprClosureUse) {
	v.dumpVertex("Var", n.Var)
}

func (v *Dumper) ExprConstFetch(n *ast.ExprConstFetch) {
	v.dumpVertex("Const", n.Const)
}

func (v *Dumper) ExprEmpty(n *ast.ExprEmpty) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprErrorSuppress(n *ast.ExprErrorSuppress) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprEval(n *ast.ExprEval) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprExit(n *ast.ExprExit) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprFunctionCall(n *ast.ExprFunctionCall) {
	v.dumpVertex("Function", n.Function)
	v.dumpVertexList("Args", n.Args)
}

func (v *Dumper) ExprInclude(n *ast.ExprInclude) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprIncludeOnce(n *ast.ExprIncludeOnce) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprInstanceOf(n *ast.ExprInstanceOf) {
	v.dumpVertex("Expr", n.Expr)
	v.dumpVertex("Class", n.Class)
}

func (v *Dumper) ExprIsset(n *ast.ExprIsset) {
	v.dumpVertexList("Vars", n.Vars)
}

func (v *Dumper) ExprList(n *ast.ExprList) {
	v.dumpVertexList("Items", n.Items)
}

func (v *Dumper) ExprMethodCall(n *ast.ExprMethodCall) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Method", n.Method)
	v.dumpVertexList("Args", n.Args)
}

func (v *Dumper) ExprNew(n *ast.ExprNew) {
	v.dumpVertex("Class", n.Class)
	v.dumpVertexList("Args", n.Args)
}

func (v *Dumper) ExprPostDec(n *ast.ExprPostDec) {
	v.dumpVertex("Var", n.Var)
}

func (v *Dumper) ExprPostInc(n *ast.ExprPostInc) {
	v.dumpVertex("Var", n.Var)
}

func (v *Dumper) ExprPreDec(n *ast.ExprPreDec) {
	v.dumpVertex("Var", n.Var)
}

func (v *Dumper) ExprPreInc(n *ast.ExprPreInc) {
	v.dumpVertex("Var", n.Var)
}

func (v *Dumper) ExprPrint(n *ast.ExprPrint) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprPropertyFetch(n *ast.ExprPropertyFetch) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Prop", n.Prop)
}

func (v *Dumper) ExprRequire(n *ast.ExprRequire) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprRequireOnce(n *ast.ExprRequireOnce) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprShellExec(n *ast.ExprShellExec) {
	v.dumpVertexList("Parts", n.Parts)
}

func (v *Dumper) ExprStaticCall(n *ast.ExprStaticCall) {
	staticCall := &ExprStaticCall{
		StartLine: n.Position.StartLine,
		EndLine:   n.Position.EndLine,
		StartPos:  n.Position.StartPos,
		EndPos:    n.Position.EndPos,
	}

	v.StaticCalls = append(v.StaticCalls, staticCall)
}

func (v *Dumper) ExprStaticPropertyFetch(n *ast.ExprStaticPropertyFetch) {
	v.dumpVertex("Class", n.Class)
	v.dumpVertex("Prop", n.Prop)
}

func (v *Dumper) ExprTernary(n *ast.ExprTernary) {
	v.dumpVertex("Cond", n.Cond)
	v.dumpVertex("IfTrue", n.IfTrue)
	v.dumpVertex("IfFalse", n.IfFalse)
}

func (v *Dumper) ExprUnaryMinus(n *ast.ExprUnaryMinus) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprUnaryPlus(n *ast.ExprUnaryPlus) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprVariable(n *ast.ExprVariable) {
	v.dumpVertex("Name", n.Name)
}

func (v *Dumper) ExprYield(n *ast.ExprYield) {
	v.dumpVertex("Key", n.Key)
	v.dumpVertex("Val", n.Val)
}

func (v *Dumper) ExprYieldFrom(n *ast.ExprYieldFrom) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprAssign(n *ast.ExprAssign) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprAssignReference(n *ast.ExprAssignReference) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprAssignBitwiseAnd(n *ast.ExprAssignBitwiseAnd) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprAssignBitwiseOr(n *ast.ExprAssignBitwiseOr) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprAssignBitwiseXor(n *ast.ExprAssignBitwiseXor) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprAssignCoalesce(n *ast.ExprAssignCoalesce) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprAssignConcat(n *ast.ExprAssignConcat) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprAssignDiv(n *ast.ExprAssignDiv) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprAssignMinus(n *ast.ExprAssignMinus) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprAssignMod(n *ast.ExprAssignMod) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprAssignMul(n *ast.ExprAssignMul) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprAssignPlus(n *ast.ExprAssignPlus) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)

}

func (v *Dumper) ExprAssignPow(n *ast.ExprAssignPow) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprAssignShiftLeft(n *ast.ExprAssignShiftLeft) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprAssignShiftRight(n *ast.ExprAssignShiftRight) {
	v.dumpVertex("Var", n.Var)
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprBinaryBitwiseAnd(n *ast.ExprBinaryBitwiseAnd) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryBitwiseOr(n *ast.ExprBinaryBitwiseOr) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryBitwiseXor(n *ast.ExprBinaryBitwiseXor) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryBooleanAnd(n *ast.ExprBinaryBooleanAnd) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryBooleanOr(n *ast.ExprBinaryBooleanOr) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryCoalesce(n *ast.ExprBinaryCoalesce) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryConcat(n *ast.ExprBinaryConcat) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryDiv(n *ast.ExprBinaryDiv) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryEqual(n *ast.ExprBinaryEqual) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryGreater(n *ast.ExprBinaryGreater) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryGreaterOrEqual(n *ast.ExprBinaryGreaterOrEqual) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryIdentical(n *ast.ExprBinaryIdentical) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryLogicalAnd(n *ast.ExprBinaryLogicalAnd) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryLogicalOr(n *ast.ExprBinaryLogicalOr) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryLogicalXor(n *ast.ExprBinaryLogicalXor) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryMinus(n *ast.ExprBinaryMinus) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryMod(n *ast.ExprBinaryMod) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryMul(n *ast.ExprBinaryMul) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryNotEqual(n *ast.ExprBinaryNotEqual) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryNotIdentical(n *ast.ExprBinaryNotIdentical) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryPlus(n *ast.ExprBinaryPlus) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryPow(n *ast.ExprBinaryPow) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryShiftLeft(n *ast.ExprBinaryShiftLeft) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinaryShiftRight(n *ast.ExprBinaryShiftRight) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinarySmaller(n *ast.ExprBinarySmaller) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinarySmallerOrEqual(n *ast.ExprBinarySmallerOrEqual) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprBinarySpaceship(n *ast.ExprBinarySpaceship) {
	v.dumpVertex("Left", n.Left)
	v.dumpVertex("Right", n.Right)
}

func (v *Dumper) ExprCastArray(n *ast.ExprCastArray) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprCastBool(n *ast.ExprCastBool) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprCastDouble(n *ast.ExprCastDouble) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprCastInt(n *ast.ExprCastInt) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprCastObject(n *ast.ExprCastObject) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ExprCastString(n *ast.ExprCastString) {
	v.dumpVertex("Expr", n.Expr)

}

func (v *Dumper) ExprCastUnset(n *ast.ExprCastUnset) {
	v.dumpVertex("Expr", n.Expr)
}

func (v *Dumper) ScalarDnumber(n *ast.ScalarDnumber) {

}

func (v *Dumper) ScalarEncapsed(n *ast.ScalarEncapsed) {
	v.dumpVertexList("Parts", n.Parts)
}

func (v *Dumper) ScalarEncapsedStringPart(n *ast.ScalarEncapsedStringPart) {

}

func (v *Dumper) ScalarEncapsedStringVar(n *ast.ScalarEncapsedStringVar) {
	v.dumpVertex("Name", n.Name)
	v.dumpVertex("Dim", n.Dim)
}

func (v *Dumper) ScalarEncapsedStringBrackets(n *ast.ScalarEncapsedStringBrackets) {
	v.dumpVertex("Var", n.Var)
}

func (v *Dumper) ScalarHeredoc(n *ast.ScalarHeredoc) {
	v.dumpVertexList("Parts", n.Parts)
}

func (v *Dumper) ScalarLnumber(n *ast.ScalarLnumber) {

}

func (v *Dumper) ScalarMagicConstant(n *ast.ScalarMagicConstant) {

}

func (v *Dumper) ScalarString(n *ast.ScalarString) {

}

func (v *Dumper) NameName(n *ast.Name) {
	v.dumpVertexList("Parts", n.Parts)
}

func (v *Dumper) NameFullyQualified(n *ast.NameFullyQualified) {
	v.dumpVertexList("Parts", n.Parts)
}

func (v *Dumper) NameRelative(n *ast.NameRelative) {
	v.dumpVertexList("Parts", n.Parts)
}

func (v *Dumper) NameNamePart(n *ast.NamePart) {

}
