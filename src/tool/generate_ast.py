# Use this Python file to generate "expr.go" and "stmt.go"
#
# python3 generate_ast.py

# https://craftinginterpreters.com/appendix-ii.html

import os

from typing import List


def java_to_go(types: List[str]) -> List[List[str]]:
    go = []
    for kv in types:
        go.append([])

        key, vs = kv.split(":")
        go[-1].append(key.strip())

        values = vs.strip().split(",")
        for value in values:
            t, m = value.strip().split(" ")
            t = t.strip()  # type
            m = m.strip()  # member
            dot_idx = t.find(".")
            if dot_idx != -1:
                t = t[: dot_idx - 4] + t[dot_idx + 1 :]
            if t[:4] == "List":
                tmp_t = "[]"
                if t[5:9].lower() not in ("expr", "stmt"):
                    tmp_t += "*"
                tmp_t += t[5:-1]
                t = tmp_t
            elif t[:6] == "Object":
                t = "interface{}"
            go[-1].append([m, t])
    return go


def define_ast(base_name: str, types: List[str]):
    elements = java_to_go(types)
    # print("Decode:", elements)

    f = []  # file content

    def a(s: str = ""):
        for i in range(len(s)):
            if s[i] == " ":
                s = s[:i] + "\t" + s[i + 1 :]
            else:
                break
        f.append(s)

    a("package glox")
    a()
    a("// This code is generated by a Python script.")
    a()
    a("type {} interface".format(base_name[0].upper() + base_name[1:]) + " {")
    a(" accept(visitor {}Visitor) (interface{{}}, error) ".format(base_name))
    a("}")
    a()

    a("type {}Visitor interface {{".format(base_name))
    for el in elements:
        s = "\t" + "visit" + el[0] + base_name[0].upper() + base_name[1:]
        s += "(" + base_name + " *" + el[0] + ") (interface{}, error) "
        a(s)
    a("}")
    a()

    for el in elements:
        a("type {} struct".format(el[0]) + " {")
        for eel in el[1:]:
            star = (
                "*"
                if eel[1] != "interface{}"
                and eel[1].lower() != "expr"
                and eel[1].lower() != "stmt"
                else ""
            )
            a("\t{} {}{}".format(eel[0], star, eel[1]))
        a("}")
        a()

        constructor_arg_list = ""
        for eel in el[1:]:
            star = (
                "*"
                if eel[1] != "interface{}"
                and eel[1].lower() != "expr"
                and eel[1].lower() != "stmt"
                else ""
            )
            constructor_arg_list += "{} {}{}, ".format(eel[0], star, eel[1])
        constructor_arg_list = constructor_arg_list[:-2]
        a("func New{0}({1}) *{0}".format(el[0], constructor_arg_list) + " {")
        a("\t{} := new({})".format(base_name, el[0]))
        for eel in el[1:]:
            a("\t{0}.{1} = {1}".format(base_name, eel[0]))
        a("\treturn {}".format(base_name))
        a("}")
        a()

        a(
            "func ({} *{}) accept(visitor {}Visitor)".format(
                base_name, el[0], base_name
            )
            + " (interface{}, error) {"
        )
        a(
            "\treturn visitor.visit{}{}({})".format(
                el[0], base_name[0].upper() + base_name[1:], base_name
            )
        )
        a("}")
        a()

    # ========================

    # print("=====")
    # for line in f:
    #     print(line)
    # print("=====")

    print("======= WARNING =======")
    file_path = os.path.abspath("../" + base_name + ".go")
    print('The program will generate file in: "{}".'.format(file_path))
    input("Continue or Ctrl+C: ")
    with open(file_path, "w") as file:
        for line in f:
            file.write(line)
            file.write("\n")
    fmt_command = 'gofmt -w "{}"'.format(file_path)
    input("Run '{}' or Ctrl+C: ".format(fmt_command))
    os.system(fmt_command)
    print("Done!")


def main():
    define_ast(
        "expr",
        [
            "Assign   : Token name, Expr value",
            "Binary   : Expr left, Token operator, Expr right",
            "Call     : Expr callee, Token paren, List<Expr> arguments",
            "Get      : Expr object, Token name",
            "Grouping : Expr expression",
            "Literal  : Object value",
            "Logical  : Expr left, Token operator, Expr right",
            "Set      : Expr object, Token name, Expr value",
            "Super    : Token keyword, Token method",
            "This     : Token keyword",
            "Unary    : Token operator, Expr right",
            "Variable : Token name",
        ],
    )
    define_ast(
        "stmt",
        [
            "Block      : List<Stmt> statements",
            "Class      : Token name, Expr.Variable superclass, List<Stmt.Function> methods",
            "Expression : Expr expression",
            "Function   : Token name, List<Token> params, List<Stmt> body",
            "If         : Expr condition, Stmt thenBranch, Stmt elseBranch",
            "Print      : Expr expression",
            "Return     : Token keyword, Expr value",
            "Var        : Token name, Expr initializer",
            "While      : Expr condition, Stmt body",
        ],
    )


if __name__ == "__main__":
    main()
