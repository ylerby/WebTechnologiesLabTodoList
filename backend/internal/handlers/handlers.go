package handlers

import (
	"backend/internal/model"
	"backend/internal/response"
	"backend/internal/schemas"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (h *Handler) Main(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte(""))
	if err != nil {
	}
}

func (h *Handler) CreateTodoList(w http.ResponseWriter, r *http.Request) {
	var responseMessage string

	reader := io.Reader(r.Body)
	requestBody, err := io.ReadAll(reader)
	if err != nil {
		responseMessage = fmt.Sprintf("ошибка при чтении тела запроса - %s", err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	currentRequestBody := model.TodoListModel{}

	err = json.Unmarshal(requestBody, &currentRequestBody)
	if err != nil {
		responseMessage = fmt.Sprintf("ошибка при десериализации объекта - %s", err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	err = h.cache.SetValue(currentRequestBody)
	if err != nil {
		responseMessage = fmt.Sprintf("ошибка при создании записи - %s", err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusServiceUnavailable)
		return
	}

	responseData := schemas.CorrectResponse{
		Data: "запись успешно создана",
	}

	result, err := json.Marshal(&responseData)
	if err != nil {
		responseMessage = fmt.Sprintf("ошибка при сериализации объекта - %s", err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	err = response.CorrectResponseWriter(w, result, http.StatusOK)
	if err != nil {
		h.logger.Errorf("ошибка при получении ответа -%s", err)
		return
	}
}

func (h *Handler) GetAllTodoLists(w http.ResponseWriter, _ *http.Request) {
	var responseMessage string

	todoLists, err := h.cache.GetAllValues()
	if err != nil {
		responseMessage = fmt.Sprintf("ошибка при получении всех записей - %s", err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusServiceUnavailable)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}
	}

	responseData := schemas.CorrectResponse{
		Data: todoLists,
	}

	result, err := json.Marshal(&responseData)
	if err != nil {
		responseMessage = fmt.Sprintf("ошибка при сериализации объекта - %s", err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}
	}

	err = response.CorrectResponseWriter(w, result, http.StatusOK)
	if err != nil {
		h.logger.Errorf("ошибка при получении ответа -%s", err)
		return
	}
}
