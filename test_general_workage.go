// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-27-08 22:34<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

package depta

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

func TestRecordMiningShouldWorkAsExpected(t *testing.T) {
	data, _ := ioutil.ReadFile("test/1.html")
	regions := simplified_api.findDataRegions(data)
	assert.Equal(t, len(regions), 28)
	records := simplified_api.findDataRecords(regions[2])
	assert.Equal(t, len(records), 13)

	getTextOfBlock := func(item *DataRecord) string {
		return strings.TrimSpace(item.convertToBase()[0].Children[0].Children[2].Children[3].Data)
	}

	assert.Equal(t, getTextOfBlock(records[0]), "Энциклопедия на французском языке. Содержит 70500 статей, 5150 иллюстраций, 245 карт.")
	assert.Equal(t, getTextOfBlock(records[1]), "Каталог выставки. На титульном листе автограф Аникушина М.К.")
	assert.Equal(t, getTextOfBlock(records[2]), "Ожерелье королевы (фр. Le Collier de la Reine) — вторая часть тетралогии Александра Дюма-отца, объединённой похождениями известного мага, предсказателя судеб и вечного человека Джузеппе Бальзамо (Калиостро). Издание полностью на французском языке.")
	assert.Equal(t, getTextOfBlock(records[3]), "Курс лекций для студентов-заочников по специальности Экономическая география")
	assert.Equal(t, getTextOfBlock(records[4]), "Том 1 - Рассказы 1906-1910 Том 2 - Рассказы 1910-1914 Том 3 - Алые паруса, Феерия; Блистающий мир, роман; Рассказы 1914-1916 Том 4 - Золотая цепь, роман; Рассказы 1916 - 1923 Том 5 - Бегущая по волнам, роман; Рассказы 1923 - 1929 Том 6 - Дорога никуда, роман; Автобиографическая повесть;Рассказы 1929 - 1930. Все тома в картонных футлярах.")
	assert.Equal(t, getTextOfBlock(records[12]), "Изложены практически все вопросы, связанные с устройством, эксплуатацией и совершенствованием яхт, содержатся разделы по основам парусного спорта, подготовке яхтенных рулевых, крейсерским плаваниям, навигации, метеорологии, уходу за яхтой и безопасности плавания. В приложении даны основы аэрогидродинамики и такелажных работ. Большой иллюстративный материал: рисунки, схемы, фотографии.")
}
