// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-27-08 22:34<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

package depta

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestNodeCopyingInvestigation(t *testing.T) {
	data, _ := ioutil.ReadFile("test/1.html")
	dtree := ReadHtml(ParseHTML(bytes.NewReader(data)))
	regions := DeptaExtract(dtree, 5, 0.75)

	first_record := regions[2].Items[0].ConvertToBase()
	second_record := regions[2].Items[1].ConvertToBase()
	assert.Equal(t, len(first_record), 13)
	assert.Equal(t, len(second_record), 13)

	assert.NotEqual(t, first_record[12].Data, second_record[12].Data)
}

func TestRecordMiningShouldWorkAsExpected(t *testing.T) {
	data, _ := ioutil.ReadFile("test/1.html")
	dtree := ReadHtml(ParseHTML(bytes.NewReader(data)))
	region_finder := MiningDataRegion{}
	region_finder.init(dtree, 5, 0.75)
	regions := region_finder.find_regions(dtree)
	assert.Equal(t, len(regions), 28)
	record_finder := MiningDataRecord{}
	record_finder.init(0.75)
	records := record_finder.find_records(regions[2])

	assert.Equal(t, len(records), 13)
	getTextOfBlock := func(item *Record) string {
		return strings.TrimSpace(item.ConvertToBase()[0].Children[0].Children[2].Children[3].Data)
	}
	assert.Equal(t, getTextOfBlock(records[0]), "Энциклопедия на французском языке. Содержит 70500 статей, 5150 иллюстраций, 245 карт.")
	assert.Equal(t, getTextOfBlock(records[1]), "Каталог выставки. На титульном листе автограф Аникушина М.К.")
	assert.Equal(t, getTextOfBlock(records[2]), "Ожерелье королевы (фр. Le Collier de la Reine) — вторая часть тетралогии Александра Дюма-отца, объединённой похождениями известного мага, предсказателя судеб и вечного человека Джузеппе Бальзамо (Калиостро). Издание полностью на французском языке.")
	assert.Equal(t, getTextOfBlock(records[3]), "Курс лекций для студентов-заочников по специальности Экономическая география")
	assert.Equal(t, getTextOfBlock(records[4]), "Том 1 - Рассказы 1906-1910 Том 2 - Рассказы 1910-1914 Том 3 - Алые паруса, Феерия; Блистающий мир, роман; Рассказы 1914-1916 Том 4 - Золотая цепь, роман; Рассказы 1916 - 1923 Том 5 - Бегущая по волнам, роман; Рассказы 1923 - 1929 Том 6 - Дорога никуда, роман; Автобиографическая повесть;Рассказы 1929 - 1930. Все тома в картонных футлярах.")
	assert.Equal(t, getTextOfBlock(records[12]), "Изложены практически все вопросы, связанные с устройством, эксплуатацией и совершенствованием яхт, содержатся разделы по основам парусного спорта, подготовке яхтенных рулевых, крейсерским плаваниям, навигации, метеорологии, уходу за яхтой и безопасности плавания. В приложении даны основы аэрогидродинамики и такелажных работ. Большой иллюстративный материал: рисунки, схемы, фотографии.")

}

func TestExamlesParsing(t *testing.T) {
	testFiles := map[string]string{
		"test/1.html": "test/result/1.html",
		"test/2.html": "test/result/2.html",
		"test/3.html": "test/result/3.html",
		"test/4.html": "test/result/4.ru",
	}
	for k, v := range testFiles {
		os.RemoveAll(v)
		file, err := ioutil.ReadFile(k)
		if err != nil {
			log.Fatal(err)
		}
		ioutil.WriteFile(v, ExtractToHTML(file), 0777)
	}

}
