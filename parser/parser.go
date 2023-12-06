package parser

import (
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
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curlToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
