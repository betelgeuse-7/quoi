package parser

import (
	"fmt"
	"quoi/ast"
	"quoi/lexer"
	"quoi/token"
	"reflect"
	"testing"
)

func printErrs(t *testing.T, errs []lexer.Err) {
	if len(errs) > 0 {
		for i, e := range errs {
			t.Logf("err#%d: %+v\n", i, e)
		}
	}
}

func printTok(t *testing.T, tok token.Token) {
	t.Logf("Token_%s(Lit: %s, Line:Col(%d:%d)\n", tok.Type.String(), tok.Literal, tok.Line, tok.Col)
}

func TestParserAdvance(t *testing.T) {
	input := "hey "
	l := lexer.New(input)
	p := New(l)
	printErrs(t, p.lexerErrors)
	//fmt.Printf("%+v\n", p)
	fmt.Println("===========")
	printTok(t, p.tok)
	printTok(t, p.peek())
	fmt.Println("===========")
	p.move()
	fmt.Println("===========")
	printTok(t, p.tok)
	printTok(t, p.peek())
	fmt.Println("===========")
	p.move()
	fmt.Println("===========")
	printTok(t, p.tok)
	printTok(t, p.peek())
	fmt.Println("===========")
}

func _parse(input string) (*ast.Program, []Err, []lexer.Err) {
	l := lexer.New(input)
	p := New(l)
	return p.Parse(), p.Errs, p.lexerErrors
}

func check_stmt_count(t *testing.T, program *ast.Program, expectedNum int) {
	if lps := len(program.Stmts); lps != expectedNum {
		t.Errorf("ERROR <!!!> len(*ast.Program.Stmts) != %d, but '%d'\n\n", expectedNum, lps)
	}
}

/*
func check_lit(t *testing.T, node ast.Node, expectedLit string) {
	if ns := node.String(); ns != expectedLit {
		t.Errorf("node.String() != '%s', but '%s'\n", expectedLit, ns)
	}
}*/

func check_error_count(t *testing.T, errs []Err, expectedNum int) {
	if lpe := len(errs); lpe != expectedNum {
		t.Errorf("ERROR <!!!> len(*Parser.Errs) != %d, but '%d'\n\n", expectedNum, lpe)
	}
}

func print_errs(t *testing.T, errs []Err) {
	for _, v := range errs {
		t.Logf("___________________________________________\n")
		t.Logf("line:col %d:%d :: %s\n", v.Line, v.Column, v.Msg)
		t.Logf("___________________________________________\n")
	}
}

func print_stmts(t *testing.T, program *ast.Program) {
	for _, v := range program.Stmts {
		t.Logf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
		t.Logf("%s :: %s\n", reflect.TypeOf(v), v)
		t.Logf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
	}
}

func TestLit1(t *testing.T) {
	input := `
			"Hello".
			1316.
			-5471.
			true.
			false.
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	check_stmt_count(t, program, 5)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestVarDecl1(t *testing.T) {
	input := `
		int n = 4.
		string x = "Hello".
		bool y = true.
		User u = "User#1".
		User u = "hey".
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	check_stmt_count(t, program, 5)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestVarDecl2(t *testing.T) {
	input := `
		int a = 1
		User u = "hey".
		bool x true.
		listof } a = [].
		string
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 4)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestVarDecl3(t *testing.T) {
	input := `
		; listof string names = ["Jennifer", "Hasan", "Ali", "Ayşe",].
		; listof string a = ["He",]
		; listof string a = ["He",
		listof string names = ["Jennifer", "Hasan", "Ali", "Ayşe"].
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestVarDecl4(t *testing.T) {
	input := `
		;User u = User {name="Jenny"}.
		User u, int x = User {name="Jenny"}, 5.
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestSubseqVar1(t *testing.T) {
	input := `
		;int x, int y = 1, 2.
		
		;int 
		;
		;x, int y = 1, 2.
	
		;int x
		;, int y = 1, 2.
		
		;int x, int
		;y, string z = -161, 1236, 
		;"HEllo".

		;int x, int y = 
		;5, 6.
		`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestSubseqVar2(t *testing.T) {
	input := `
		;listof int x, listof string y = [], [].
		int x, listof User ux, bool y = 5, [User {name="Hello"}], true. 
		`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestDT1(t *testing.T) {
	input := `
		datatype City {
			string name
			int x
			int y 
			bool z
			User u

			; hello
		}
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestDT2(t *testing.T) {
	input := `
		; datatype City {}
		; datatype {}
		; datatype X {
		; datatype X }
		; datatype X { int x }
		; datatype X { 
		;	int x string name
		;}
		datatype X { int x
		}
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestBlock1(t *testing.T) {
	input := `
		block 
			Stdout::println("Hey").
			print_it(1416).
		end
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestBlock2(t *testing.T) {
	input := `
;		block 
;			Stdout::println("Hey").
;			print_it(1416).
;	block
;	end
	block "Hey". end
`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestPrefExpr1(t *testing.T) {
	input := `
		(+ 1 4).								; 5
		(+ (* Int::from_string("5") 5) 2). 		; 27
		(' ["Hey", "Hello"] 0).					; "Hey"									


	(* 2 Int::from_string(String::from_int(
		(+ 3 5 18925
			Int::from_string("-1516")
		),
	))). 										; 34834

	Stdout::println((* 4 Math::pow(2, 2))).		; 16
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestPrefExpr2(t *testing.T) {
	input := `
		;(m 4 5 67).
		;().
		;(+ 5 6 7 8 9
		;(' [0, 1, 2, 3, 4] 2)
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 1)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestPrefExpr3(t *testing.T) {
	input := `
		User u = User{ name="User 1" }.
		u = (set u name "Jenny"). 
		Stdout::println( (get u name) ).
		(gte 5 6 6 ).
		(get u name).
		(gte x y).
		(+ a b).
		(+ aa bb).
		`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestRA1(t *testing.T) {
	input := `
		name = "Hey".
		age = 51.
		u = "User".
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestRA2(t *testing.T) {
	input := `
		;age =.
		;name = "Hey"
		;u = 
		a=    "716".
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestReturn1(t *testing.T) {
	input := `
		; return
		; return.
		; return .
		; return 5
		return Math::sqrt(64).
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestReturn2(t *testing.T) {
	input := `
		;return "Hey", a, true, b.
		;return 1,.
		;return .
		;return 1 2.
		return "Hey", "a", true.
		return 1.
		`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestBreakAndContinue(t *testing.T) {
	input := `
		;break
		;continue
		block
			continue.
			break.

		end
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestFC1(t *testing.T) {
	input := `
		string_concat("Hello", "World").
		Os::read_file("hello.txt").
		Math::pow(
			2, 2,
		).
		Stdout::println(
			1, 2,  3,
			5, "Hello", "Yay",
		).
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestFC2(t *testing.T) {
	input := `
		;string_concat("Hello" "World")
		;Os::read_file ("hello.txt" ).
		;Math::pow(2, 2,).
		;Math::pow(2, 2)
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestLoop1(t *testing.T) {
	input := `
		;loop true {}
		;loop (lte 4 5) {
		;	Stdout::println("4 is less than or equal to 5").
		;}
		;loop (and 
		;	(not false) 
		;	(lt -1 0)) {
		;	Stdout::println("wow much complex").
		;}
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestLoop2(t *testing.T) {
	input := `
		;loop {}
		;loop true {
		;loop true }
		;loop true {{}
		;loop fun {}
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 1)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestIf1(t *testing.T) {
	input := `
		if true {} elseif true {} else {}
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestIf2(t *testing.T) {
	input := `
		if (gt 6 -1) {
			Stdout::println("6 is greater than -1").
		} elseif (not (lt 5 6)) {
			Stdout::println("::D").
		} else {
			Stdout::println( (+ 2 2) ).			
		}
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestIf3(t *testing.T) {
	input := `
		if (gt 6 -1) {
			Stdout::println("6 is greater than -1").
		} elseif false { Stderr::println("Errororroorr"). }
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestIf4(t *testing.T) {
	input := `
		if (gt 6 -1) {
			Stdout::println("6 is greater than -1").
		} else { Stdout::println("else ran").  }
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestIf5(t *testing.T) {
	input := `
	elseif true { Stdout::println("else ran").  }
;	else { Stdout::println("else ran").  }
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 1)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestParseFunParams1(t *testing.T) {
	input := `
		(int a, int b)
	`
	l := lexer.New(input)
	p := New(l)
	p.move()
	p.move()
	fmt.Println(p.tok)
	params := p.parseFunctionParams("fake_fn")
	fmt.Println(params)
	for i, v := range params {
		t.Logf("param#%d: (%s %s)\n", i, v.Tok.Literal, v.Name.String())
	}
}

func TestFD1(t *testing.T) {
	input := `
		;fun a() {}
		;fun
		;fun a
		;fun a(
		;fun a()
		;fun a() {
		;fun a() -> {}
		fun a() -> {
			Stdout::println("Yes").
			Stdout::println( Math::avg(1, 2, 3) ).
		}
		`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestFD2(t *testing.T) {
	input := `
		;fun a(int x) -> { }
		;fun a(int x, int
		;	) {} 
		;fun a(int x, int y,)  {}
		;fun a(int x,
		;int y) -> { Stdout::println("Hello"). Stdout::println("Hey"). }
		
		;fun a(int
		;x) -> {}
		
		;fun a(int y, 
		;	string x) -> {}
		
		;fun datatype() -> {}
		
		;fun a(User u, -> {}
		;fun a(User u,) -> {}

		;fun a() {}
		;fun a( listof int x, listof string y ) -> {}
		;fun a( listof listof string yy ) -> {} 
		;fun a(listof string y, City city) -> {}
		;fun a(listof string y, City city
		;	) -> {}

;		fun a(
;			listof string y, City city
;			) -> {}
		`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestFD3(t *testing.T) {
	input := `
			;fun a() -> string {}
			
			;fun a() -> 
			;string {}
		
			;fun a(int a) -> int, {}

			;fun a(User u) -> int,
			;{}

			;fun a(listof User ux) -> int
			;{}

			;fun a(listof int nx) -> int, string {}
			
			;fun a() -> int string int {}
			
;			fun a() -> string,
;			int, 
;			bool {
;				return "Hey", 1, true.
;			}

			;fun a() -> string
			;bool {}
			`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestDTL1(t *testing.T) {
	input := `
		;User{name="Jennifer" age=34}.	
		;City {name="X" founded_in = 1935}.
		;Person{}.
		
		;Person 
		;{}.

		;Person {
		;
		;}.

		;Lexer 
		;{
		;
		;}.

		Monster
		{
			name = "Very Scary Monster"
			power = 1620
			lives_in = Swamp {
				swamp_number=12
				kmsqr = 10
			}
			nested = Nested{ 
				Datatype=Jenny {
					na="me"
				}
			}
		}.
		`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestDTL2(t *testing.T) {
	input := `
		;User {
		;User {}
			`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestErrorRecovery1(t *testing.T) {
	input := `
		int x, { y = 6, 7.
	`
	_, errs, _ := _parse(input)
	print_errs(t, errs)
}

func TestDatatypeListField(t *testing.T) {
	input := `
		datatype X {
			listof string y
		}
	`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}

func TestListDecl1(t *testing.T) {
	input := `
		listof int nx = [1, 2, 3].
		;listof int nx2 = .;nx.
		listof int nx3 = nx.
s		listof int nxx = "hey".
		`
	program, errs, _ := _parse(input)
	check_error_count(t, errs, 0)
	print_stmts(t, program)
	print_errs(t, errs)
}
