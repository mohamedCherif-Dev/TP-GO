package main

import (
	"flag"
	"fmt"
	"os"

	"gopack/internal/backup"
	"gopack/internal/fetch"
	"gopack/internal/scan"
)

var version = "dev"

func printHelp() {
	fmt.Fprintf(os.Stderr, `gopack - Outil CLI de sauvegarde et transfert de fichiers (v%s)

Usage:
  gopack <commande> [arguments]

Commandes disponibles:
  scan    Analyse un répertoire et affiche un inventaire par extension
  backup  Copie un répertoire vers une destination avec préservation des droits
  fetch   Télécharge de manière concurrente une liste d'URLs
  version Affiche la version
  help    Affiche cet message d'aide

Utilisez "gopack <commande> -h" pour obtenir de l'aide sur une sous-commande.
`, version)
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "scan":
		scanFS := flag.NewFlagSet("scan", flag.ExitOnError)
		scanFS.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage: gopack scan <répertoire>\n\nInventaire du répertoire : parcours récursif et statistiques par extension.\n")
			scanFS.PrintDefaults()
		}
		_ = scanFS.Parse(os.Args[2:])
		if scanFS.NArg() < 1 {
			fmt.Fprintln(os.Stderr, "Erreur usage: chemin de répertoire manquant")
			scanFS.Usage()
			os.Exit(2)
		}
		if err := scan.Run(scanFS.Arg(0)); err != nil {
			fmt.Fprintf(os.Stderr, "gopack scan: %v\n", err)
			os.Exit(1)
		}

	case "backup":
		backupFS := flag.NewFlagSet("backup", flag.ExitOnError)
		backupFS.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage: gopack backup <src> <dst>\n\nCopie récursive d'un répertoire en préservant les permissions.\n")
			backupFS.PrintDefaults()
		}
		_ = backupFS.Parse(os.Args[2:])
		if backupFS.NArg() < 2 {
			fmt.Fprintln(os.Stderr, "Erreur usage: chemins source et destination requis")
			backupFS.Usage()
			os.Exit(2)
		}
		if err := backup.Run(backupFS.Arg(0), backupFS.Arg(1)); err != nil {
			fmt.Fprintf(os.Stderr, "gopack backup: %v\n", err)
			os.Exit(1)
		}

	case "fetch":
		fetchFS := flag.NewFlagSet("fetch", flag.ExitOnError)
		outDir := fetchFS.String("o", ".", "Répertoire de destination des téléchargements")
		workers := fetchFS.Int("n", 4, "Nombre de téléchargements concurents max")
		fetchFS.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage: gopack fetch [-o dir] [-n N] <url1> [url2...]\n\nTéléchargement concurrent d'URLs avec contrôle de la concurrence.\n")
			fetchFS.PrintDefaults()
		}
		_ = fetchFS.Parse(os.Args[2:])
		if fetchFS.NArg() < 1 {
			fmt.Fprintln(os.Stderr, "Erreur usage: au moins une URL est requise")
			fetchFS.Usage()
			os.Exit(2)
		}
		if err := fetch.Run(fetchFS.Args(), *outDir, *workers); err != nil {
			fmt.Fprintf(os.Stderr, "gopack fetch: %v\n", err)
			os.Exit(1)
		}

	case "version":
		fmt.Printf("gopack version %s\n", version)

	case "help", "-h", "--help":
		printHelp()

	default:
		fmt.Fprintf(os.Stderr, "gopack: commande inconnue %q\n\n", os.Args[1])
		printHelp()
		os.Exit(2)
	}
}
