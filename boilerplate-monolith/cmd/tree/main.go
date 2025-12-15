package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort" // <-- YENİ: Sıralama paketi eklendi
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

func printTree(w io.Writer, path string, prefix string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Fprintf(w, "%s[Erişim Engellendi]\n", prefix+"└── ")
		return nil
	}

	// 1. Filtreleme (Gizli dosyalar hariç)
	var visibleEntries []os.DirEntry
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), ".") && e.Name() != "tree_output_sorted.txt" {
			visibleEntries = append(visibleEntries, e)
		}
	}

	// 2. SIRALAMA MANTIĞI (Önce Klasörler, Sonra Dosyalar)
	sort.Slice(visibleEntries, func(i, j int) bool {
		entryI := visibleEntries[i]
		entryJ := visibleEntries[j]

		// İkisi de klasör veya ikisi de dosya ise -> Alfabetik sırala
		if entryI.IsDir() == entryJ.IsDir() {
			return entryI.Name() < entryJ.Name()
		}

		// Biri klasör biri dosya ise -> Klasör olanı (IsDir=true) öne al
		// true dönersek i önce gelir, false dönersek j önce gelir.
		return entryI.IsDir()
	})

	// 3. Yazdırma Döngüsü
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
