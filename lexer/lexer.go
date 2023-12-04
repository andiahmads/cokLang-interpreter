package lexer

import (
	"go-intepreter/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

// postion & readPosition
// Keduanya akan digunakan untuk mengakses karakter dalam input dengan menggunakannya sebagai indeks

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// Tujuan dari readChar adalah untuk memberi kita karakter berikutnya dan memajukan posisi kita dalam input string.
// Hal pertama yang dilakukannya adalah memeriksa apakah kita telah mencapai akhir input.
// Jika sudah maka ia akan mengatur l.ch ke 0, yang merupakan kode ASCII untuk karakter "NUL" dan menandakan "kita belum membaca apapun" atau "akhir file" untuk kita.
// Tetapi jika kita belum mencapai akhir dari input, maka ia akan mengeset l.ch ke karakter berikutnya dengan mengakses l.input[l.readPosition].
// Setelah itu l.position diperbarui ke l.readPosition yang baru saja digunakan dan l.readPosition bertambah satu.
// Dengan begitu, l.readPosition selalu menunjuk ke posisi berikutnya di mana kita akan untuk membaca dari berikutnya dan l.position selalu menunjuk ke posisi di mana kita terakhir kali membaca
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition

	l.readPosition += 1

}

// Yang perlu dilakukan oleh lexer kita adalah mengenali apakah karakter saat ini adalah huruf,
// dan jika ya, ia perlu membaca sisa pengenal/kata kunci hingga menemukan karakter yang bukan huruf.
// Setelah membaca pengenal/kata kunci tersebut, kita perlu untuk mengetahui apakah itu adalah pengenal atau kata kunci, sehingga kita dapat menggunakan token.TokenType yang benar.
// Langkah Langkah pertama adalah memperluas pernyataan switch kita:
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		tok = newToken(token.BANG, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default: // memeriksa pengenal kapan pun l.ch bukan salah satu karakter yang dikenali
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal) //lexing pengenal dan kata kunci
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILEGAL, l.ch) //menangani karakter saat ini dan mendeklarasikan menyatakannya sebagai token.ILLEGAL.
		}
	}

	l.readChar()
	return tok
}

// membaca sebuah identifier dan memajukan posisi lexer kita sampai bertemu dengan karakter non-huruf
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// memeriksa apakah argumen yang diberikan adalah sebuah huruf.
// tetapi yang perlu diperhatikan tentang isLetter adalah bahwa mengubah fungsi ini memiliki besar pada bahasa yang dapat diurai oleh interpreter kita daripada yang diharapkan dari
// fungsi kecil. Seperti yang Anda lihat, dalam kasus kita, fungsi ini berisi pemeriksaan ch == '_', yang berarti bahwa
// kita akan memperlakukan _ sebagai huruf dan mengizinkannya dalam pengenal dan kata kunci. Ini berarti kita dapat menggunakan variabel
// nama variabel seperti foo_bar. Bahasa pemrograman lain bahkan mengizinkan ! dan ? dalam pengenal. Jika Anda
// ingin mengizinkannya juga, ini adalah tempat untuk menyelipkannya.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// skipWhitespace(), lexer akan melewatkan angka 5 pada bagian let five = 5; dari pengujian kita masukan
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// currently only supports integers, for other data types such as float, double, for now it is not supported.
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// Fungsi ini hanya mengembalikan apakah byte yang dimasukkan adalah sebuah Angka Latin antara 0 dan 9.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
