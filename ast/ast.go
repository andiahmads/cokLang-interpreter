package ast

import "go-intepreter/token"

// AST adalah singkatan dari "Abstract Syntax Tree" atau "Pohon Sintaksis Abstrak" dalam bahasa Indonesia.
// AST adalah struktur data pohon yang merepresentasikan struktur sintaksis dari suatu kode sumber setelah proses parsing.
// Setiap node dalam AST mewakili elemen-elemen sintaksis dari kode sumber, seperti pernyataan, ekspresi, dan deklarasi.

type Node interface {
	// membuat tiga bidang: satu untuk pengenal, satu untuk ekspresi yang menghasilkan nilai dalam pernyataan let dan satu untuk token.
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

// Node Program ini akan menjadi node akar dari setiap AST yang dihasilkan oleh pemilah kita.
// Setiap program COKLnag yang valid adalah serangkaian pernyataan.
// Pernyataan-pernyataan ini terkandung dalam Program.Statements,
// yang hanya merupakan potongan dari node AST yang mengimplementasikan antarmuka Statement.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// LetStatement memiliki bidang-bidang yang kita butuhkan: Nama untuk menyimpan pengenal pengikatan dan Nilai untuk yang menghasilkan nilai.
// Dua metode statementNode dan TokenLiteral memenuhi antarmuka Statement dan Node masing-masing.
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// Untuk menjaga jumlah tipe simpul yang berbeda tetap kecil,
// kita akan menggunakan Identifier di sini untuk merepresentasikan nama dalam pengikatan variabel dan kemudian menggunakannya kembali,
// untuk merepresentasikan sebuah pengenal sebagai bagian dari atau sebagai lengkap dari sebuah ekspresi.
type Identifier struct {
	Token token.Token
	Value string
}

// Untuk menyimpan pengenal dari pengikatan, x dalam let x = 5;, kita memiliki tipe pengenal struct, yang mengimplementasikan antarmuka Expression.
// Tetapi pengenal dalam pernyataan let tidak menghasilkan nilai, bukan? Jadi mengapa ini adalah sebuah Ekspresi? Ini untuk membuat segalanya tetap sederhana. Pengenal di bagian lain
// lain dari program Monkey memang menghasilkan nilai, contoh: let x = valueProducingIdentifier;.
func (i *Identifier) expressionNode() {
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// Dengan Program, LetStatement dan Identifier mendefinisikan bagian kode sumber COKLang ini;
// let x = 5;
