package parser

import (
	"fmt"
	"go-intepreter/ast"
	"go-intepreter/lexer"
	"go-intepreter/token"
	"strconv"
)

// Parser adalah suatu program atau fungsi dalam suatu sistem komputer yang bertugas untuk menganalisis dan memproses suatu inputan dalam bentuk teks atau data
// menghasilkan output yang dapat digunakan oleh komputer atau aplikasi lainnya.
// Parser berperan dalam mengurai atau memecah input menjadi struktur data yang dapat diakses atau dimanfaatkan oleh program.
// Tugas utama parser adalah mengonversi inputan dalam bentuk teks atau data menjadi struktur data yang dapat dimengerti oleh komputer atau aplikasi.
// Proses ini disebut parsing. Parser membaca input, mengenali pola-pola tertentu, dan kemudian membentuk representasi internal dari input tersebut.

const (
	_           int = iota
	LOWEST          //merupakan suatu konstanta atau nilai tertentu yang menunjukkan tingkat precedensi terendah.
	EQUALS          // ==
	LESSGREATER     // > or <
	SUM             // +
	PRODUCT         // *
	PREFIX          // -X or !X
	CALL            // myFunction(X)
)

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
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFn   map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:     l,
		erros: []string{},
	}

	p.nextToken()
	p.nextToken()

	// 	jika kita menemukan token bertipe token.IDENT, fungsi parsing yang akan dipanggil adalah
	// parseIdentifier, metode yang kita definisikan pada *Parser.
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)

	p.registerPrefix(token.INT, p.parseIntegralLiteral)

	// register prefix operator
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	// register infix operator
	p.infixParseFn = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	return p
}

// Metode ini hanya mengembalikan sebuah *ast.Identifier dengan token saat ini di bidang Token dan nilai literal token di Value. Metode ini tidak memajukan
// token, tidak memanggil nextToken. Ini sangat penting.
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curlToken, Value: p.curlToken.Literal}
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
		return p.parseExpressionStatement()
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

// mem-parse pernyataan ekspresi jika kita tidak menemukan salah satu dari dua(let & return):
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curlToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// mengoperkan precedence serendah mungkin ke parseExpression
// Yang dilakukannya adalah memeriksa apakah kita memiliki fungsi parsing yang terkait
// dengan p.curToken.Type di posisi awalan. Jika ada, ia akan memanggil fungsi parsing ini, jika tidak ada, ia akan
// mengembalikan nilai nol.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curlToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curlToken.Type)
		return nil
	}
	leftExp := prefix()

	// 	mencoba menemukan infixParseFns untuk token berikutnya. Jika ia menemukan fungsi tersebut, ia akan memanggilnya, melewatkan dalam ekspresi yang dikembalikan oleh prefixParseFn sebagai argumen.
	// Dan ia melakukan semua ini lagi dan lagi sampai ia menemukan token yang memiliki prioritas lebih tinggi.
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFn[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)

	}
	return leftExp
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
	// Kedua tipe fungsi  ini mengembalikan sebuah ast.Expression,
	prefixParseFn func() ast.Expression
	// infixParseFn mengambil argumen: ast.Expression lain.
	infixParseFn func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFn[tokenType] = fn
}

// INTEGRAL LITERAL
// Literal bilangan bulat adalah ekspresi. Nilai yang dihasilkannya adalah bilangan bulat itu sendiri. Sekali lagi,
// bayangkan di mana saja integer literal dapat muncul untuk memahami mengapa mereka adalah ekspresi
// let x = 5;
// add(5, 10);
// 5 + 5 + 5;
// Kita bisa menggunakan ekspresi lain selain literal bilangan bulat di sini dan ekspresi tersebut akan tetap valid

// Seperti parseIdentifier, metode ini sangat sederhana. Satu-satunya hal yang benar-benar berbeda adalah pemanggilan strconv.ParseInt, yang mengubah string di p.curToken.Literal menjadi
// int64. Int64 tersebut kemudian disimpan ke dalam bidang Value dan kita mengembalikan * simpul *ast.IntegerLiteral. Jika tidak berhasil, kita menambahkan kesalahan baru ke dalam kesalahan parser field
func (p *Parser) parseIntegralLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curlToken}

	value, err := strconv.ParseInt(p.curlToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curlToken.Literal)
		p.erros = append(p.erros, msg)
		return nil
	}

	lit.Value = value
	return lit
}

// Prefix Operators
// Ada dua operator awalan dalam bahasa pemrograman CokLang: ! dan -.
// Penggunaan mereka adalah hampir sama dengan apa yang Anda harapkan dari bahasa-bahasa lain:
// -5;
// !foobar;
// 5 + -10;
// Struktur penggunaannya adalah sebagai berikut:
// <operator awalan <ekspresi>;
// !isGreaterThanZero(2);
// 5 + -add(5, 5);
// Ini berarti bahwa simpul AST untuk ekspresi operator awalan harus cukup fleksibel untuk menunjuk ke ekspresi apa pun sebagai operan.

// noPrefixParseFnError hanya menambahkan pesan kesalahan yang diformat ke
// bidang kesalahan pada parser kita. Tetapi itu cukup untuk mendapatkan pesan kesalahan yang lebih baik dalam pengujian yang gagal
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.erros = append(p.erros, msg)
}

// Untuk token.BANG dan token.MINUS kita mendaftarkan metode yang sama dengan prefixParseFn:
// metode yang baru saja dibuat, yaitu parsePrefixExpression.
// Metode ini membangun simpul AST, dalam hal ini *ast.PrefixExpression, sama seperti fungsi penguraian yang telah kita lihat sebelumnya.
// Tetapi kemudian  melakukan sesuatu yang berbeda: metode ini benar-benar memajukan token kita dengan memanggil p.nextToken()!
// Ketika parsePrefixExpression dipanggil, p.curToken akan bertipe token.BANG atau token.MINUS, karena jika tidak, fungsi ini tidak akan dipanggil.
// Tetapi untuk mengurai ekspresi awalan dengan benar seperti -5, lebih dari satu token harus "dikonsumsi".
// Jadi setelah menggunakan p.curToken untuk membangun sebuah *ast.PrefixExpression, metode ini memajukan token-token dan memanggil parseExpression lagi.
// ketika parseExpression dipanggil oleh parsePrefixExpression, token-token telah dimajukan dan token saat ini adalah token setelah operator awalan. Dalam kasus -5, ketika parseExpression dipanggil, p.curToken.Type adalah token.INT.
// parseExpression kemudian memeriksa yang terdaftar dan menemukan parseIntegerLiteral, yang membangun sebuah *ast.IntegerLiteral dan mengembalikannya.
// parseExpression mengembalikan simpul yang baru.
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curlToken,
		Operator: p.curlToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

// Infix Operators
// Infix operators adalah operator-operasi matematika dan logika yang ditempatkan di antara dua operand
// Selanjutnya kita akan menguraikan kedelapan operator infiks ini:
// 5 + 5;
// 5 - 5;
// 5 * 5;
// 5 / 5;
// 5 > 5;
// 5 < 5;
// 5 == 5;
// 5 != 5;

// Tingkat precedensi (precedence level) adalah konsep yang digunakan dalam pemrograman untuk menentukan urutan evaluasi atau eksekusi operasi dalam suatu ekspresi.
// Ini mengindikasikan seberapa kuat suatu operator memengaruhi operan-operandnya.
// Contoh umumnya dapat ditemukan dalam ekspresi matematika. Misalnya, dalam ekspresi 3 + 5 * 2, kita tahu bahwa perkalian (*) memiliki tingkat precedensi yang lebih tinggi daripada penambahan (+).
// Oleh karena itu, ekspresi tersebut akan dievaluasi sebagai 3 + (5 * 2), karena perkalian memiliki prioritas lebih tinggi.
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

// Metode peekPrecedence mengembalikan prioritas yang terkait dengan tipe token p.peekToken.
// Jika tidak menemukan prioritas untuk p.peekToken, maka akan menjadi LOWEST, yaitu prioritas terendah yang mungkin dimiliki oleh operator mana pun.
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Metode curPrecedence melakukan hal yang sama, tetapi untuk p.curToken.
func (p *Parser) curlPrecedence() int {
	if p, ok := precedences[p.curlToken.Type]; ok {
		return p
	}
	return LOWEST
}

//	metode baru ini mengambil sebuah argumen, sebuah ast.Expression bernama left.
//
// Metode ini menggunakan argumen ini untuk membuat simpul *ast.InfixExpression, dengan left berada di bidang Left.
// Kemudian metode ini menetapkan prioritas token saat ini (yang merupakan operator dari ekspresi infix) ke prioritas variabel lokal.
// sebelum memajukan token dengan memanggil nextToken dan mengisi bidang Kanan node dengan panggilan lain untuk menguraiEkspresi - kali ini dengan melewatkan prioritas token operator.
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curlToken,
		Operator: p.curlToken.Literal,
		Left:     left,
	}

	precedence := p.curlPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}
