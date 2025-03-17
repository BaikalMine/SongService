package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/BaikalMine/SongService/database"
	"github.com/BaikalMine/SongService/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetSongs godoc
// @Summary Получение списка песен
// @Description Получение списка песен с фильтрацией по группе и названию, а также с пагинацией.
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Фильтр по группе"
// @Param song query string false "Фильтр по названию песни"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество записей на странице" default(10)
// @Success 200 {array} models.Song
// @Failure 500 {object} map[string]string
// @Router /songs [get]
func GetSongs(c *gin.Context, db *sql.DB) {
	groupFilter := c.Query("group")
	songFilter := c.Query("song")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := "SELECT id, group_name, song_name, release_date, lyrics, link FROM songs WHERE 1=1"
	var args []interface{}
	argIndex := 1

	if groupFilter != "" {
		query += fmt.Sprintf(" AND group_name ILIKE $%d", argIndex)
		args = append(args, "%"+groupFilter+"%")
		argIndex++
	}
	if songFilter != "" {
		query += fmt.Sprintf(" AND song_name ILIKE $%d", argIndex)
		args = append(args, "%"+songFilter+"%")
		argIndex++
	}
	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	var songs []models.Song
	// Оборачиваем операцию выборки в транзакцию для обеспечения консистентности
	err = database.WithTransaction(db, func(tx *sql.Tx) error {
		rows, err := tx.Query(query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var s models.Song
			if err := rows.Scan(&s.ID, &s.Group, &s.Song, &s.ReleaseDate, &s.Lyrics, &s.Link); err != nil {
				logrus.Errorf("Ошибка чтения данных: %v", err)
				continue
			}
			songs = append(songs, s)
		}
		return nil
	})
	if err != nil {
		logrus.Errorf("Ошибка запроса песен: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
		return
	}
	c.JSON(http.StatusOK, songs)
}

// GetSongLyrics godoc
// @Summary Получение текста песни
// @Description Получение текста песни, разделённого на куплеты, с пагинацией.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество куплетов на страницу" default(1)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /songs/{id}/lyrics [get]
func GetSongLyrics(c *gin.Context, db *sql.DB) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID песни"})
		return
	}

	var lyrics string
	// Выполнение запроса чтения текста песни в транзакции
	err = database.WithTransaction(db, func(tx *sql.Tx) error {
		return tx.QueryRow("SELECT lyrics FROM songs WHERE id = $1", id).Scan(&lyrics)
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Песня не найдена"})
		return
	}

	verses := strings.Split(lyrics, "\n\n")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 1
	}
	start := (page - 1) * limit
	end := start + limit
	if start > len(verses) {
		c.JSON(http.StatusOK, gin.H{"verses": []string{}, "total": len(verses)})
		return
	}
	if end > len(verses) {
		end = len(verses)
	}
	c.JSON(http.StatusOK, gin.H{"verses": verses[start:end], "total": len(verses)})
}

// AddSong godoc
// @Summary Добавление новой песни
// @Description Добавление новой песни с обогащением через внешний API.
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Данные песни (обязательны поля group и song)"
// @Success 201 {object} map[string]int "ID добавленной песни"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs [post]
func AddSong(c *gin.Context, db *sql.DB, externalAPIUrl string) {
	var input struct {
		Group string `json:"group" binding:"required"`
		Song  string `json:"song" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}
	logrus.Infof("Добавление новой песни: %s - %s", input.Group, input.Song)

	// Запрос к внешнему API для получения обогащённой информации
	url := fmt.Sprintf("%s/info?group=%s&song=%s", externalAPIUrl, input.Group, input.Song)
	resp, err := http.Get(url)
	if err != nil {
		logrus.Errorf("Ошибка запроса к внешнему API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить информацию о песне"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Внешнее API вернуло статус: %d", resp.StatusCode)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных из внешнего API"})
		return
	}

	var songDetail struct {
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		logrus.Errorf("Ошибка декодирования ответа внешнего API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Некорректный ответ внешнего API"})
		return
	}

	var newID int
	err = database.WithTransaction(db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO songs (group_name, song_name, release_date, lyrics, link)
			VALUES ($1, $2, $3, $4, $5) RETURNING id
		`
		return tx.QueryRow(query, input.Group, input.Song, songDetail.ReleaseDate, songDetail.Text, songDetail.Link).Scan(&newID)
	})
	if err != nil {
		logrus.Errorf("Ошибка сохранения песни в БД: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось добавить песню"})
		return
	}
	logrus.Infof("Песня добавлена, ID: %d", newID)
	c.JSON(http.StatusCreated, gin.H{"id": newID})
}

// UpdateSong godoc
// @Summary Обновление песни
// @Description Обновление данных песни по ID.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param song body models.Song true "Обновлённые данные песни"
// @Success 200 {object} map[string]string "Песня обновлена"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs/{id} [put]
func UpdateSong(c *gin.Context, db *sql.DB) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID песни"})
		return
	}

	var input struct {
		Group       string `json:"group"`
		Song        string `json:"song"`
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	err = database.WithTransaction(db, func(tx *sql.Tx) error {
		query := `
			UPDATE songs
			SET group_name = $1, song_name = $2, release_date = $3, lyrics = $4, link = $5
			WHERE id = $6
		`
		_, err := tx.Exec(query, input.Group, input.Song, input.ReleaseDate, input.Text, input.Link, id)
		return err
	})
	if err != nil {
		logrus.Errorf("Ошибка обновления песни: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить песню"})
		return
	}
	logrus.Infof("Песня с ID %d обновлена", id)
	c.JSON(http.StatusOK, gin.H{"message": "Песня обновлена"})
}

// DeleteSong godoc
// @Summary Удаление песни
// @Description Удаление песни по ID.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {object} map[string]string "Песня удалена"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /songs/{id} [delete]
func DeleteSong(c *gin.Context, db *sql.DB) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID песни"})
		return
	}

	err = database.WithTransaction(db, func(tx *sql.Tx) error {
		query := "DELETE FROM songs WHERE id = $1"
		res, err := tx.Exec(query, id)
		if err != nil {
			return err
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return fmt.Errorf("Песня не найдена")
		}
		return nil
	})
	if err != nil {
		if err.Error() == "Песня не найдена" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Песня не найдена"})
		} else {
			logrus.Errorf("Ошибка удаления песни: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить песню"})
		}
		return
	}
	logrus.Infof("Песня с ID %d удалена", id)
	c.JSON(http.StatusOK, gin.H{"message": "Песня удалена"})
}
