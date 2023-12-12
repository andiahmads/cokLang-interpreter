package parser

import (
	"fmt"
	"go-intepreter/ast"
	"go-intepreter/lexer"
	"go-intepreter/token"
)

// Parser adalah suatu program atau fungsi dalam suatu sistem komputer yang bertugas untuk menganalisis dan memproses suatu inputan dalam bentuk teks atau data
// menghasilkan output yang dapat digunakan oleh komputer atau aplikasi lainnya.
// Parser berperan dalam mengurai atau memecah input menjadi struktur data yang dapat diakses atau dimanfaatkan oleh program.
// Tugas utama parser adalah mengonversi inputan dalam bentuk teks atau data menjadi struktur data yang dapat dimengerti oleh komputer atau aplikasi.
// Proses ini disebut parsing. Parser membaca input, mengenali pola-pola tertentu, dan kemudian membentuk representasi internal dari input tersebut.

// Parser memiliki tiga bidang: l, curToken dan peekToken.
// l adalah sebuah penunjuk ke sebuah instance dari lexer, di mana kita berulang kali memanggil NextToken() untuk mendapatkan token berikutnya dalam input.
type Parser struct {
	l *lexer.Lexer
	// 	curToken dan peekToken bertindak persis seperti dua "penunjuk" yang dimiliki oleh lexer kita: position dan peekPosition.
	// Namun, alih-alih menunjuk pada karakter dalam input, mereka menunjuk pada token saat ini dan token berikutnya.
	curlToken token.Token // curToken, yang merupakan token saat ini yang sedang diperiksa
	peekToken token.Token //untuk memutuskan apa yang harus dilakukan selanjutnya, dan kita juga membutuhkan peekToken, untuk memutuskan apakah kita berada di akhir baris atau apakah kita apakah kita berada di awal ekspresi aritmatika.

	erros []string

	// 	Dengan adanya peta-peta ini, kita tinggal memeriksa apakah peta yang sesuai (infiks atau awalan) memiliki penguraian
	// yang terkait dengan curToken.Type.
	prefixParseFn map[token.TokenType]prefixParseFn
	infixParseFn  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:     l,
		erros: []string{},
	}

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curlToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram adalah membuat simpul akar dari AST, sebuah *ast.Program.
// kemudian mengulang setiap token dalam masukan sampai menemukan token.EOF
// Ini dilakukan dengan berulang kali memanggil nextToken yang memajukan p.curToken dan p.peekToken.
// Dalam setiap iteras ia memanggil parseStatement, yang tugasnya adalah mengurai pernyataan.
// Jika parseStatement mengembalikan sesuatu selain nil, sebuah ast.Statement, nilai kembaliannya ditambahkan ke irisan Pernyataan dari simpul akar AST. Jika tidak ada yang tersisa untuk diurai, simpul akar *ast.Program akan dikembalikan.
func (p *Parser) ParseProgram() *ast.Program {
	// logic ini didapat file pseudocode
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curlToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()

	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curlToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()

	default:
		return nil
	}
}

// membangun simpul *ast.LetStatement dengan token yang saat ini berada (token token.LET)
// kemudian memajukan token sambil membuat pernyataan tentang token berikutnya dengan panggilan ke expectPeek
// Pertama, ia mengharapkan token token.IDENT, yang kemudian kemudian digunakan untuk membuat simpul *ast.Identifier
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curlToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curlToken, Value: p.curlToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: Kita melewatkan ekspresi sampai kita
	// menemukan titik koma
	for !p.curlTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curlToken}

	p.nextToken()

	for !p.curlTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// menemukan titik koma/SEMICOLON
func (p *Parser) curlTokenIs(t token.TokenType) bool {
	return p.curlToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// Metode expectPeek adalah salah satu dari "fungsi pernyataan" yang dimiliki oleh hampir semua pengurai.
// Tujuan utama mereka adalah untuk menegakkan kebenaran urutan token dengan memeriksa jenis token berikutnya.
// memeriksa jenis dari peekToken dan hanya jika jenisnya benar, maka ia akan memajukan token dengan memanggil nextToken.
// Seperti yang akan Anda lihat, ini adalah sesuatu yang sering dilakukan oleh pengurai.
// Namun, apa yang terjadi jika kita menemukan token di expectPeek yang tidak sesuai dengan jenis yang diharapkan?
// Pada saat ini, kita hanya mengembalikan nilai nol, yang akan diabaikan dalam ParseProgram, yang mengakibatkan seluruh
// yang diabaikan karena ada kesalahan pada input.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.erros
}

// memeriksa apakah parser menemukan kesalahan apa pun.
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.erros = append(p.erros, msg)
}

// Parsing Expressions
// 5 * 5 + 10 -> AST = (5*5) + 10.
// parser harus mengetahui tentang prioritas operator.
// di mana prioritas * lebih tinggi dari +. Itu adalah contoh yang paling umum untuk operator.
// kasus lain: 5 * (5 + 10), penjumlahan sekarang harus dievaluasi sebelum perkalian. Itu karena tanda kurung memiliki prioritas yang lebih tinggi daripada operator *
// kasus lain: -5 - 10, Di sini operator - muncul di awal ekspresi, sebagai operator awalan, dan kemudian sebagai operator infiks di tengah
// kasus lain: 5 * (add(2, 3) + 10)
// Validitas dari sebuah posisi token sekarang tergantung pada konteks, token yang datang sebelum dan sesudahnya, dan yang didahulukan.

// terminology
// Operator awalan adalah operator "di depan" operan. Contoh:
// --5
// Di sini operatornya adalah -- (pengurangan), operannya adalah bilangan bulat literal 5 dan operatornya berada di
// berada di posisi awalan.
// Operator postfix adalah operator "setelah" operan. Contoh:
// foobar++
//  operator infix adalah sesuatu yang telah kita lihat sebelumnya. Operator infix berada di antara
// operan, seperti ini:
// 5 * 8

// 5 + 5 * 10
// Hasil dari ekspresi ini adalah 55, bukan 100. Dan itu karena operator * memiliki prioritas yang lebih tinggi
// lebih tinggi, sebuah "peringkat yang lebih tinggi". Ini "lebih penting" daripada operator +

// kesimpulan:
// parsing operator memiliki beberapa isitalah
// prefix, postfix, infix operator and precedence.

// Implementing the Pratt Parser
// Ide utama parser Pratt adalah asosiasi fungsi penguraian (yang disebut Pratt sebagai "semantik semantik") dengan tipe-tipe token
// Kapan pun jenis token ini ditemukan, fungsi penguraian akan dipanggil untuk mem-parsing ekspresi yang sesuai dan mengembalikan simpul AST yang mewakilinya.
// Setiap tipe token dapat memiliki hingga dua fungsi penguraian yang terkait dengannya, tergantung pada apakah
// token ditemukan dalam posisi awalan atau akhiran.
// Hal pertama yang perlu kita lakukan adalah mengatur asosiasi ini. Kami mendefinisikan dua jenis fungsi:
// fungsi penguraian awalan dan fungsi penguraian infiks.
type (
	// 	prefixParseFns akan dipanggil ketika kita menemukan token yang terkait
	// di posisi awalan dan infixParseFn dipanggil ketika kita menemukan tipe token di posisi infix
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFn[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFn[tokenType] = fn
}
