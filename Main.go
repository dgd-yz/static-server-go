package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const IndexBeginTemplate = `<!DOCTYPE HTML><html lang="en"><head><meta charset="utf-8"><style type="text/css">:root {color-scheme: light dark;} .bi-folder {color: #ffd54f;} .bi-file-earmark { color: #90caf9;} </style><title>Directory listing for {{FILE_PATH}}</title></head><body><h1>Directory listing for {{FILE_PATH}}</h1><hr><ul>`

const FileNameTemplate = `<li><a href="{{FILE_NAME_URL}}">{{FILE_NAME_TEXT}}</a></li>`

const IndexEndTemplate = `</ul><hr></body></html>`

const NotFoundTemplate = `<!DOCTYPE HTML><html lang="en"><head><meta charset="utf-8"><style type="text/css">:root {color-scheme: light dark;}</style><title>Error response</title></head><body><h1>Error response</h1><p>Error code: 404</p><p>Message: File not found.</p><p>Error code explanation: 404 - Nothing matches the given URI.</p></body></html>`

const fileIconSvg = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-file-earmark" viewBox="0 0 16 16"><path d="M14 4.5V14a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V2a2 2 0 0 1 2-2h5.5zm-3 0A1.5 1.5 0 0 1 9.5 3V1H4a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V4.5z"/></svg>`

const directoryIconSvg = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-folder" viewBox="0 0 16 16"><path d="M.54 3.87.5 3a2 2 0 0 1 2-2h3.672a2 2 0 0 1 1.414.586l.828.828A2 2 0 0 0 9.828 3h3.982a2 2 0 0 1 1.992 2.181l-.637 7A2 2 0 0 1 13.174 14H2.826a2 2 0 0 1-1.991-1.819l-.637-7a2 2 0 0 1 .342-1.31zM2.19 4a1 1 0 0 0-.996 1.09l.637 7a1 1 0 0 0 .995.91h10.348a1 1 0 0 0 .995-.91l.637-7A1 1 0 0 0 13.81 4zm4.69-1.707A1 1 0 0 0 6.172 2H2.5a1 1 0 0 0-1 .981l.006.139q.323-.119.684-.12h5.396z"/></svg>`

func main() {
	multiplexador := http.NewServeMux()
	multiplexador.HandleFunc("/", handleStatic)

	endereco := ":8080"
	log.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(endereco, multiplexador))
}

func handleStatic(resposta http.ResponseWriter, requisicao *http.Request) {

	exePath, err := os.Executable()
	if err != nil {
		http.Error(resposta, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}
	root := filepath.Dir(exePath)

	cleanPath := path.Clean("/" + requisicao.URL.Path)
	caminho := filepath.Join(root, cleanPath)

	info, err := os.Stat(caminho)
	if err == nil && info.IsDir() && !strings.HasSuffix(requisicao.URL.Path, "/") {
		http.Redirect(resposta, requisicao, requisicao.URL.Path+"/", http.StatusMovedPermanently)
		return
	}

	if !strings.HasPrefix(caminho, root) {
		resposta.WriteHeader(http.StatusForbidden)
		resposta.Write([]byte("<h1>403 Forbidden</h1><p>Você não pode acessar fora da pasta root.</p>"))
		return
	}

	fmt.Println(caminho)

	var construtor strings.Builder

	ext := path.Ext(caminho)

	contentType := mime.TypeByExtension(ext)
	if ext == ".js" {
		contentType = "text/javascript"
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	if ext == "" {
		contentType = "text/html; charset=utf-8"

		entradas, err := os.ReadDir(caminho)

		indexPath := filepath.Join(caminho, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			arquivo, err := os.Open(indexPath)
			if err == nil {
				defer arquivo.Close()
				resposta.Header().Set("Content-Type", "text/html; charset=utf-8")
				resposta.WriteHeader(http.StatusOK)
				io.Copy(resposta, arquivo)
				return
			}
		}

		if err != nil {
			resposta.Header().Set("Content-Type", "text/html; charset=utf-8")
			resposta.WriteHeader(http.StatusNotFound)
			resposta.Write([]byte(NotFoundTemplate))
			return
		}

		relPath, err := filepath.Rel(root, caminho)
		if err != nil {
			relPath = ""
		}

		displayPath := "/" + filepath.ToSlash(relPath)
		if displayPath == "/" {
			displayPath = "/"
		}

		construtor.WriteString(strings.Replace(IndexBeginTemplate, "{{FILE_PATH}}", displayPath, -1))

		if requisicao.URL.Path != "/" {
			voltarUrl := path.Join(requisicao.URL.Path, "..")
			voltar := strings.Replace(FileNameTemplate, "{{FILE_NAME_TEXT}}", "..", -1)
			voltar = strings.Replace(voltar, "{{FILE_NAME_URL}}", voltarUrl, -1)
			construtor.WriteString(voltar)
		}

		for _, entrada := range entradas {
			nome := entrada.Name()

			texto := ""
			url := nome

			if entrada.IsDir() {
				texto = directoryIconSvg + " " + nome
				url = nome + "/"
			} else {
				texto = fileIconSvg + " " + nome
			}

			linha := FileNameTemplate
			linha = strings.Replace(linha, "{{FILE_NAME_TEXT}}", texto, -1)
			linha = strings.Replace(linha, "{{FILE_NAME_URL}}", url, -1)

			construtor.WriteString(linha)
		}

		construtor.WriteString(IndexEndTemplate)
	} else {
		arquivo, err := os.Open(caminho)
		if err != nil {
			resposta.Header().Set("Content-Type", "text/html; charset=utf-8")
			resposta.WriteHeader(http.StatusNotFound)
			resposta.Write([]byte(NotFoundTemplate))
			return
		}
		defer arquivo.Close()

		io.Copy(&construtor, arquivo)
	}
	resposta.Header().Set("Content-Type", contentType)
	resposta.WriteHeader(http.StatusOK)

	_, _ = resposta.Write([]byte(construtor.String()))
}
