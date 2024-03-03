package handlers

import (
	"backend/internal/auth"
	"backend/internal/model"
	"backend/internal/response"
	"backend/internal/schemas"
	"backend/internal/service"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	minPasswordSize       = 12
	maxPasswordSize       = 255
	responseErrorKey      = "Error"
	successfulValueCreate = "запись успешно создана"
	successfulValueUpdate = "запись успешно обновлена"
	successfulValueDelete = "запись успешно удалена"
	successfulUserCreate  = "пользователь успешно создан"
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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var responseMessage string

	reader := io.Reader(r.Body)
	requestBody, err := io.ReadAll(reader)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при чтении тела запроса - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	currentRequestBody := auth.Credentials{}

	err = json.Unmarshal(requestBody, &currentRequestBody)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при десериализации объекта - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	user, isFound, err := h.database.GetUser(currentRequestBody.Login)
	if err != nil {
		h.logger.Errorf("ошибка при получении значения - %s", err)
		responseMessage = fmt.Sprintf("%s: ошибка при авторизации", responseErrorKey)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusServiceUnavailable)
		return
	}

	if !isFound {
		responseMessage = fmt.Sprintf("%s: пользователь с таким логином не найден", responseErrorKey)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusBadRequest)
		return
	}

	err = service.CompareHashAndPassword(currentRequestBody.Password, user.Password)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при авторизации", responseErrorKey)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusUnauthorized)
		return
	} else {
		token, err := auth.GenerateToken(currentRequestBody.Login)
		if err != nil {
			responseMessage = fmt.Sprintf("%s: ошибка при генерации токена - %s", responseErrorKey, err)
			err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
			return
		}

		responseData := schemas.CorrectResponse{
			Data: token,
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
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var responseMessage string

	reader := io.Reader(r.Body)
	requestBody, err := io.ReadAll(reader)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при чтении тела запроса - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	currentRequestBody := schemas.RegisterUserRequestBody{}

	err = json.Unmarshal(requestBody, &currentRequestBody)
	if err != nil {
		responseMessage = fmt.Sprintf("%s: ошибка при десериализации объекта - %s", responseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	if len(currentRequestBody.Password) < minPasswordSize || len(currentRequestBody.Password) > maxPasswordSize {
		responseMessage = fmt.Sprintf("%s: некорректная длина пароля. "+
			"(мин. длина пароля - 12 символов, макс. - 255)", responseErrorKey)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusBadRequest)
		return
	}

	if currentRequestBody.Password != currentRequestBody.AgainPassword {
		responseMessage = fmt.Sprintf("%s: пароли не совпадают", responseErrorKey)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	hashedPassword, err := service.HashPassword(currentRequestBody.Password)
	if err != nil {
		h.logger.Errorf("ошибка при хешировании паролей - %s", err)
		responseMessage = fmt.Sprintf("%s: ошибка при регистрации", responseErrorKey)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusServiceUnavailable)
		return
	}

	alreadyExist, err := h.database.CreateUser(currentRequestBody.Login, hashedPassword)
	if err != nil {
		h.logger.Errorf("ошибка при создании пользователя - %s", err)
		responseMessage = fmt.Sprintf("%s: ошибка при создании пользователя", responseErrorKey)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		return
	}

	if alreadyExist {
		responseMessage = fmt.Sprintf("%s: пользователь с таким логином уже существует", responseErrorKey)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusBadRequest)
		return
	}

	responseData := schemas.CorrectResponse{
		Data: successfulUserCreate,
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
