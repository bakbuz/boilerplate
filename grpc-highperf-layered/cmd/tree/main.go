package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	// Çıktı dosyasını oluştur
	fileName := "../../tree_output.txt"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Dosya oluşturulamadı: %v\n", err)
		return
	}
	defer file.Close()

	writer := io.MultiWriter(os.Stdout, file)

	// Mevcut dizini al ve 2 üst klasöre çık
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(writer, "Mevcut dizin okunamadı: %v\n", err)
		return
	}

	root := filepath.Clean(filepath.Join(currentDir, "../../"))

	fmt.Fprintf(writer, "Tarama Yolu: %s\n", root)
	fmt.Fprintln(writer, strings.Repeat("-", 30))
	fmt.Fprintln(writer, filepath.Base(root)+"/")

	if err := printTree(writer, root, ""); err != nil {
		fmt.Fprintf(writer, "Hata: %v\n", err)
	}
}

// Dosya/Klasörün sıralamadaki "Ağırlığını" belirleyen fonksiyon
func getSortWeight(entry os.DirEntry) int {
	name := entry.Name()

	// 1. Özel Klasörler (En Sona Atılacaklar)
	if name == "scripts" {
		return 2
	}
	if name == "deployment" {
		return 3
	}

	// 2. Normal Klasörler (En Başa)
	if entry.IsDir() {
		return 1
	}

	// 3. Dosyalar
	return 4
}

func printTree(w io.Writer, path string, prefix string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Fprintf(w, "%s[Erişim Engellendi]\n", prefix+"└── ")
		return nil
	}

	// Filtreleme
	var visibleEntries []os.DirEntry
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), ".") && e.Name() != "tree_custom_sort.txt" {
			visibleEntries = append(visibleEntries, e)
		}
	}

	// ÖZEL SIRALAMA MANTIĞI
	sort.Slice(visibleEntries, func(i, j int) bool {
		entryI := visibleEntries[i]
		entryJ := visibleEntries[j]

		weightI := getSortWeight(entryI)
		weightJ := getSortWeight(entryJ)

		// Eğer ağırlıkları farklıysa (örn: biri normal klasör, biri scripts), ağırlığa göre sırala
		if weightI != weightJ {
			return weightI < weightJ
		}

		// Eğer ağırlıkları aynıysa (ikisi de normal klasör), alfabetik sırala
		return entryI.Name() < entryJ.Name()
	})

	// Yazdırma Döngüsü
	for i, entry := range visibleEntries {
		isLast := i == len(visibleEntries)-1

		connector := "├── "
		if isLast {
			connector = "└── "
		}

		displayName := entry.Name()
		if entry.IsDir() {
			displayName += "/"
		}

		fmt.Fprintln(w, prefix+connector+displayName)

		if entry.IsDir() {
			newPrefix := prefix
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}

			subPath := filepath.Join(path, entry.Name())
			if err := printTree(w, subPath, newPrefix); err != nil {
				return err
			}
		}
	}
	return nil
}
