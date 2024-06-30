/*Palacios	Chavez	Roberto Paolo
Sulca	Mamani	Ivan Frank
Reyes	Rojas	Martin Abel
Vasquez	Castañeda	Jhonn Anderson*/

package main

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"pp/database"
)

// Definir la estructura para almacenar los datos del postulante
type Postulante struct {
	Cod_vcCodigo   string
	Ide_iIndice    int
	Esc_vcCodigo   string
	Esc_vcNombre   string
	Are_cCodigo    string
	Cal_fNotaFinal float64
	Merito         int
}

func main() {
	database.Connect()
	defer database.Close()
	start := time.Now()

	query := "SELECT postulante.cod_vcCodigo, calificacion.ide_iIndice, postulante.esc_vcCodigo, escuela.esc_vcNombre, escuela.are_cCodigo, calificacion.cal_fNotaFinal FROM postulante JOIN identificacion ON postulante.cod_vcCodigo = identificacion.cod_vcCodigo JOIN calificacion ON identificacion.ide_iIndice = calificacion.ide_iIndice JOIN escuela ON postulante.esc_vcCodigo = escuela.esc_vcCodigo"
	rows, err := database.Query(query)

	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	areaChannels := make(map[string]chan Postulante)
	areas := []string{"A", "B", "C", "D", "E"}

	for _, area := range areas {
		areaChannels[area] = make(chan Postulante)
	}

	go func() {
		defer func() {
			for _, ch := range areaChannels {
				close(ch)
			}
		}()
		for rows.Next() {
			var postulante Postulante
			err := rows.Scan(&postulante.Cod_vcCodigo, &postulante.Ide_iIndice, &postulante.Esc_vcCodigo, &postulante.Esc_vcNombre, &postulante.Are_cCodigo, &postulante.Cal_fNotaFinal)
			if err != nil {
				panic(err)
			}
			areaChannels[postulante.Are_cCodigo] <- postulante
		}
	}()

	processArea := func(area string, areaChannel chan Postulante, result *[]Postulante) {
		defer wg.Done()
		for postulante := range areaChannel {
			*result = append(*result, postulante)
		}
		sort.Slice(*result, func(i, j int) bool {
			return (*result)[i].Cal_fNotaFinal > (*result)[j].Cal_fNotaFinal
		})
		for i := range *result {
			(*result)[i].Merito = i + 1
			// Actualizar mérito en la base de datos
			_, err := database.Exec("UPDATE calificacion SET cal_iMeritoGeneral = ? WHERE ide_iIndice = ?", (*result)[i].Merito, (*result)[i].Ide_iIndice)
			if err != nil {
				panic(err)
			}
		}
	}

	var areasResults = make(map[string][]Postulante)
	for _, area := range areas {
		areasResults[area] = []Postulante{}
		wg.Add(1)
		go func(area string) {
			var tempResults []Postulante
			processArea(area, areaChannels[area], &tempResults)
			areasResults[area] = tempResults
		}(area)
	}

	wg.Wait()

	

	elapsed := time.Since(start).Seconds()
	
	fmt.Printf("Tiempo de ejecución: %.8f s\n", elapsed)
}
