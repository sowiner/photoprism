package query

import (
	form "github.com/photoprism/photoprism/internal/form"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/photoprism/photoprism/internal/config"
)

func TestQuery_AlbumByUUID(t *testing.T) {
	conf := config.TestConfig()

	search := New(conf.Db())

	t.Run("existing uuid", func(t *testing.T) {
		albums, err := search.AlbumByUUID("3")
		assert.Nil(t, err)
		assert.Equal(t, "Christmas2030", albums.AlbumName)
	})

	t.Run("not existing uuid", func(t *testing.T) {
		albums, err := search.AlbumByUUID("3765")
		assert.Error(t, err, "record not found")
		t.Log(albums)
	})
}

func TestQuery_AlbumThumbByUUID(t *testing.T) {
	conf := config.TestConfig()

	search := New(conf.Db())

	t.Run("existing uuid", func(t *testing.T) {
		file, err := search.AlbumThumbByUUID("4")
		assert.Nil(t, err)
		assert.Equal(t, "exampleFileName.jpg", file.FileName)
	})

	t.Run("not existing uuid", func(t *testing.T) {
		file, err := search.AlbumThumbByUUID("3765")
		assert.Error(t, err, "record not found")
		t.Log(file)
	})
}

func TestQuery_Albums(t *testing.T) {
	conf := config.TestConfig()

	search := New(conf.Db())

	t.Run("search with string", func(t *testing.T) {
		query := form.NewAlbumSearch("chr")
		result, err := search.Albums(query)
		assert.Nil(t, err)
		assert.Equal(t, "Christmas2030", result[0].AlbumName)
	})

	t.Run("search with slug", func(t *testing.T) {
		query := form.NewAlbumSearch("slug:holiday count:10")
		result, err := search.Albums(query)
		assert.Nil(t, err)
		assert.Equal(t, "Holiday2030", result[0].AlbumName)
	})

	t.Run("favorites true", func(t *testing.T) {
		query := form.NewAlbumSearch("favorites:true count:10000")

		result, err := search.Albums(query)
		assert.Nil(t, err)
		assert.Equal(t, "Holiday2030", result[0].AlbumName)
	})
	t.Run("empty query", func(t *testing.T) {
		query := form.NewAlbumSearch("order:slug")

		result, err := search.Albums(query)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(result))
	})
	t.Run("search with invalid query string", func(t *testing.T) {
		query := form.NewAlbumSearch("xxx:bla")
		result, err := search.Albums(query)
		assert.Error(t, err, "unknown filter")
		t.Log(result)
	})
}
