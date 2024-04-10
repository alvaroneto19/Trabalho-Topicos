package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var FiltarCaracteresIndesejados = func(itens string) string {

	itens = strings.ReplaceAll(itens, `\[\d+:\d+:\d+,\d+\]`, "")

	RemoverCaracteresIndejados := func(caracteres rune) bool {
		caracteresParaRemover := `?.,!@♪-:1234567890><()/\`
		return strings.ContainsRune(caracteresParaRemover, caracteres)
	}

	return strings.Map(func(caracteresRemovidos rune) rune {
		if RemoverCaracteresIndejados(caracteresRemovidos) {
			return -1
		}
		return caracteresRemovidos
	}, itens)

}

var DeixarPalavrasMinuculas = func(palavras string) string {
	return strings.ToLower(palavras)
}

func main() {
	os.MkdirAll("Resultados", os.ModePerm)

	PastaDeLegendas := "./legendas"

	ListagemLegendas, err := filepath.Glob(filepath.Join(PastaDeLegendas, "*.srt"))
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(ListagemLegendas); i++ {

		file := ListagemLegendas[i]

		f, err := os.Open(file)

		if err != nil {
			panic(err)
		}
		defer f.Close()

		ConteudoLegendas, err := io.ReadAll(f)
		if err != nil {
			panic(err)
		}

		textoFiltrado := FiltarCaracteresIndesejados(string(ConteudoLegendas))
		palavras := strings.Fields(textoFiltrado)

		contagem := make(map[string]int)

		for j := 0; j < len(palavras); j++ {
			palavra := palavras[j]
			contagem[palavra]++
		}

		NomeDoArquivoJSON := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)) + "_.json"
		CaminhoPastaLegendas := filepath.Join("Resultados", NomeDoArquivoJSON)
		arquivo, err := os.Create(CaminhoPastaLegendas)
		if err != nil {
			panic(err)
		}
		defer arquivo.Close()

		ConversãoParaJSON := json.NewEncoder(arquivo)
		ConversãoParaJSON.SetIndent("  ", "\t")
		ConversãoParaJSON.SetEscapeHTML(true)

		var palavrasOrdenadas []struct {
			Palavra    string `json:"palavra"`
			Frequencia int    `json:"frequencia"`
		}
		for palavra, frequencia := range contagem {
			palavrasOrdenadas = append(palavrasOrdenadas, struct {
				Palavra    string `json:"palavra"`
				Frequencia int    `json:"frequencia"`
			}{
				Palavra:    palavra,
				Frequencia: frequencia,
			})
		}

		sort.Slice(palavrasOrdenadas, func(i, j int) bool {
			return palavrasOrdenadas[i].Frequencia > palavrasOrdenadas[j].Frequencia

		})
		ConversãoParaJSON.Encode(palavrasOrdenadas)

	}

}
