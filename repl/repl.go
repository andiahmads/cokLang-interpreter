package repl

import (
	"bufio"
	"fmt"
	"go-intepreter/lexer"
	"go-intepreter/token"
	"io"
)

// Bahasa COK membutuhkan REPL. REPL adalah singkatan dari "Read Eval Print Loop"
// Terkadang REPL disebut "konsol", terkadang "mode interaktif".
// Konsepnya adalah sama: REPL REPL membaca input, mengirimkannya ke interpreter untuk dievaluasi, mencetak hasil/keluaran dari penerjemah dan memulai lagi. Baca, Evaluasi, Cetak, Ulangi.
// Kita belum tahu bagaimana cara "mengevaluasi" kode sumber COK sepenuhnya.
// Kita hanya memiliki satu bagian dari proses yang bersembunyi di balik "Eval":
const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf("%s", PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
