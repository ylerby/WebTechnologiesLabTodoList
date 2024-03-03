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

const (
	responseErrorKey      = "Error"
	successfulValueCreate = "запись успешно создана"
	successfulValueUpdate = "запись успешно обновлена"
	successfulValueDelete = "запись успешно удалена"
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
		responseMessage = fmt.Sprintf("%s: ошибка при чтении тела запроса - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	currentRequestBody := model.TodoListModel{}

	err = json.Unmarshal(requestBody, &currentRequestBody)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при десериализации объекта - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	err = h.cache.SetValue(currentRequestBody)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при создании записи - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusServiceUnavailable)
		return
	}

	responseData := schemas.CorrectResponse{
		Data: successfulValueCreate,
	}

	result, err := json.Marshal(&responseData)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при сериализации объекта - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	err = response.CorrectResponseWriter(w, result, http.StatusOK)
	if err != nil {
		h.logger.Errorf("ошибка при получении ответа -%s", err)
		return
	}
}

func (h *Handler) GetTodoListByTitle(w http.ResponseWriter, r *http.Request) {
	var responseMessage string

	reader := io.Reader(r.Body)
	requestBody, err := io.ReadAll(reader)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при чтении тела запроса - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	currentRequestBody := schemas.GetTodoListByTitle{}
	err = json.Unmarshal(requestBody, &currentRequestBody)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при десериализации объекта - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	todoLists, err := h.cache.GetValueByTitle(currentRequestBody.Title)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при получении записей - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusServiceUnavailable)
		return
	}

	responseData := schemas.CorrectResponse{
		Data: todoLists,
	}

	result, err := json.Marshal(&responseData)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при сериализации объекта - %s", responseErrorKey, err)
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
		responseMessage = fmt.Sprintf("%s: ошибка при получении всех записей - %s", responseErrorKey, err)
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
		responseMessage = fmt.Sprintf("%s: ошибка при сериализации объекта - %s", responseErrorKey, err)
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

func (h *Handler) UpdateTodoList(w http.ResponseWriter, r *http.Request) {
	var responseMessage string

	reader := io.Reader(r.Body)
	requestBody, err := io.ReadAll(reader)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при чтении тела запроса - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	currentRequestBody := schemas.UpdateTodoList{}

	err = json.Unmarshal(requestBody, &currentRequestBody)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при десериализации объекта - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	err = h.cache.UpdateValue(currentRequestBody)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при обновлении записи - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusServiceUnavailable)
		return
	}

	responseData := schemas.CorrectResponse{
		Data: successfulValueUpdate,
	}

	result, err := json.Marshal(&responseData)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при сериализации объекта - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	err = response.CorrectResponseWriter(w, result, http.StatusOK)
	if err != nil {
		h.logger.Errorf("ошибка при получении ответа -%s", err)
		return
	}
}

func (h *Handler) DeleteTodoList(w http.ResponseWriter, r *http.Request) {
	var responseMessage string

	reader := io.Reader(r.Body)
	requestBody, err := io.ReadAll(reader)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при чтении тела запроса - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	currentRequestBody := model.TodoListModel{}

	err = json.Unmarshal(requestBody, &currentRequestBody)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при десериализации объекта - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	//fixme
	err = h.cache.DeleteValue(currentRequestBody)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при удалении записи - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusServiceUnavailable)
		return
	}

	responseData := schemas.CorrectResponse{
		Data: successfulValueDelete,
	}

	result, err := json.Marshal(&responseData)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при сериализации объекта - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	err = response.CorrectResponseWriter(w, result, http.StatusOK)
	if err != nil {
		h.logger.Errorf("ошибка при получении ответа -%s", err)
		return
	}
}
