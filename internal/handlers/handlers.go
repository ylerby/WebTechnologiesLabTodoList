package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"backend/internal/csv_encoder"

	"backend/internal/auth"
	"backend/internal/domain"
	"backend/internal/response"
	"backend/internal/service"
)

const (
	minPasswordSize         = 12
	maxPasswordSize         = 255
	ResponseErrorKey        = "Error"
	csvResponseFormatKey    = "csv"
	successfulValueCreate   = "запись успешно создана"
	successfulValueUpdate   = "запись успешно обновлена"
	successfulValueDelete   = "запись успешно удалена"
	successfulUserCreate    = "пользователь успешно создан"
	successfulCommentCreate = "комментарий успешно добавлен"
)

func (h *Handler) Main(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte(""))
	if err != nil {
		h.logger.Errorf("ошибка при получении ответа -%s", err)
		return
	}
}

func (h *Handler) CreateTodoList(w http.ResponseWriter, r *http.Request) {
	var responseMessage string

	reader := io.Reader(r.Body)
	requestBody, err := io.ReadAll(reader)
	if err != nil {
		h.logger.Errorf("%s: ошибка при чтении тела запроса - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при чтении тела запроса - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	currentRequestBody := domain.TodoListModel{}

	err = json.Unmarshal(requestBody, &currentRequestBody)
	if err != nil {
		h.logger.Errorf("%s: ошибка при десериализации объекта - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при десериализации объекта - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	err = h.cache.SetValue(currentRequestBody)
	if err != nil {
		h.logger.Errorf("%s: ошибка при создании записи - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при создании записи - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusServiceUnavailable)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	responseData := domain.CorrectResponse{
		Data: successfulValueCreate,
	}

	result, err := json.Marshal(&responseData)
	if err != nil {
		h.logger.Errorf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}
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
		h.logger.Errorf("%s: ошибка при чтении тела запроса - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при чтении тела запроса - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	currentRequestBody := domain.GetTodoListByTitle{}
	err = json.Unmarshal(requestBody, &currentRequestBody)
	if err != nil {
		h.logger.Errorf("%s: ошибка при десериализации объекта - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при десериализации объекта - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	todoLists, err := h.cache.GetValueByTitle(currentRequestBody.Title)
	if err != nil {
		h.logger.Errorf("%s: ошибка при получении записей - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при получении записей - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusServiceUnavailable)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	if currentRequestBody.IsCsv {
		err = csv_encoder.Encode(w, todoLists)
		if err != nil {
			h.logger.Errorf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
			responseMessage = fmt.Sprintf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
			err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
			if err != nil {
				h.logger.Errorf("ошибка при получении ответа -%s", err)
				return
			}
		}

		return
	}

	responseData := domain.CorrectResponse{
		Data: todoLists,
	}

	result, err := json.Marshal(&responseData)
	if err != nil {
		h.logger.Errorf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	err = response.CorrectResponseWriter(w, result, http.StatusOK)
	if err != nil {
		h.logger.Errorf("ошибка при получении ответа -%s", err)
		return
	}
}

func (h *Handler) GetAllTodoLists(w http.ResponseWriter, r *http.Request) {
	var responseMessage string

	todoLists, err := h.cache.GetAllValues()
	if err != nil {
		h.logger.Errorf("%s: ошибка при получении всех записей - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при получении всех записей - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusServiceUnavailable)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}
	}

	format := r.URL.Query().Get("format")
	if format == csvResponseFormatKey {
		err = csv_encoder.Encode(w, todoLists)
		if err != nil {
			h.logger.Errorf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
			responseMessage = fmt.Sprintf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
			err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
			if err != nil {
				h.logger.Errorf("ошибка при получении ответа -%s", err)
				return
			}
		}

		return
	}

	responseData := domain.CorrectResponse{
		Data: todoLists,
	}

	result, err := json.Marshal(&responseData)
	if err != nil {
		h.logger.Errorf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
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
		h.logger.Errorf("%s: ошибка при чтении тела запроса - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при чтении тела запроса - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	currentRequestBody := domain.UpdateTodoList{}

	err = json.Unmarshal(requestBody, &currentRequestBody)
	if err != nil {
		h.logger.Errorf("%s: ошибка при десериализации объекта - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при десериализации объекта - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	err = h.cache.UpdateValue(currentRequestBody)
	if err != nil {
		h.logger.Errorf("%s: ошибка при обновлении записи - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при обновлении записи - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusServiceUnavailable)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	responseData := domain.CorrectResponse{
		Data: successfulValueUpdate,
	}

	result, err := json.Marshal(&responseData)
	if err != nil {
		h.logger.Errorf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

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
		h.logger.Errorf("%s: ошибка при чтении тела запроса - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при чтении тела запроса - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	currentRequestBody := domain.TodoListModel{}

	err = json.Unmarshal(requestBody, &currentRequestBody)
	if err != nil {
		h.logger.Errorf("%s: ошибка при десериализации объекта - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при десериализации объекта - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	err = h.cache.DeleteValue(currentRequestBody)
	if err != nil {
		h.logger.Errorf("%s: ошибка при удалении записи - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при удалении записи - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusServiceUnavailable)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	responseData := domain.CorrectResponse{
		Data: successfulValueDelete,
	}

	result, err := json.Marshal(&responseData)
	if err != nil {
		h.logger.Errorf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

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
		h.logger.Errorf("%s: ошибка при чтении тела запроса - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при чтении тела запроса - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	currentRequestBody := auth.Credentials{}

	err = json.Unmarshal(requestBody, &currentRequestBody)
	if err != nil {
		h.logger.Errorf("%s: ошибка при десериализации объекта - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при десериализации объекта - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	user, isFound, err := h.database.GetUser(currentRequestBody.Login)
	if err != nil {
		h.logger.Errorf("ошибка при получении значения - %s", err)
		responseMessage = fmt.Sprintf("%s: ошибка при авторизации", ResponseErrorKey)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusServiceUnavailable)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	if !isFound {
		h.logger.Errorf("%s: пользователь с таким логином не найден", ResponseErrorKey)
		responseMessage = fmt.Sprintf("%s: пользователь с таким логином не найден", ResponseErrorKey)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusBadRequest)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	err = service.CompareHashAndPassword(currentRequestBody.Password, user.Password)
	if err != nil {
		h.logger.Errorf("%s: ошибка при авторизации", ResponseErrorKey)
		responseMessage = fmt.Sprintf("%s: ошибка при авторизации", ResponseErrorKey)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusUnauthorized)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	} else {
		token, err := auth.GenerateToken(currentRequestBody.Login)
		if err != nil {
			h.logger.Errorf("%s: ошибка при генерации токена - %s", ResponseErrorKey, err)
			responseMessage = fmt.Sprintf("%s: ошибка при генерации токена - %s", ResponseErrorKey, err)
			err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
			if err != nil {
				h.logger.Errorf("ошибка при получении ответа -%s", err)
				return
			}

			return
		}

		responseData := domain.CorrectResponse{
			Data: token,
		}

		result, err := json.Marshal(&responseData)
		if err != nil {
			h.logger.Errorf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
			responseMessage = fmt.Sprintf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
			err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
			if err != nil {
				h.logger.Errorf("ошибка при получении ответа -%s", err)
				return
			}

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
		h.logger.Errorf("%s: ошибка при чтении тела запроса - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при чтении тела запроса - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	currentRequestBody := domain.RegisterUserRequestBody{}

	err = json.Unmarshal(requestBody, &currentRequestBody)
	if err != nil {
		h.logger.Errorf("%s: ошибка при десериализации объекта - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при десериализации объекта - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	if len(currentRequestBody.Password) < minPasswordSize || len(currentRequestBody.Password) > maxPasswordSize {

		responseMessage = fmt.Sprintf("%s: некорректная длина пароля. "+
			"(мин. длина пароля - 12 символов, макс. - 255)", ResponseErrorKey)
		h.logger.Errorf(responseMessage)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusBadRequest)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	if currentRequestBody.Password != currentRequestBody.AgainPassword {
		h.logger.Errorf("%s: пароли не совпадают", ResponseErrorKey)
		responseMessage = fmt.Sprintf("%s: пароли не совпадают", ResponseErrorKey)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	hashedPassword, err := service.HashPassword(currentRequestBody.Password)
	if err != nil {
		h.logger.Errorf("ошибка при хешировании паролей - %s", err)
		responseMessage = fmt.Sprintf("%s: ошибка при регистрации", ResponseErrorKey)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusServiceUnavailable)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	alreadyExist, err := h.database.CreateUser(currentRequestBody.Login, hashedPassword)
	if err != nil {
		h.logger.Errorf("ошибка при создании пользователя - %s", err)
		responseMessage = fmt.Sprintf("%s: ошибка при создании пользователя", ResponseErrorKey)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	if alreadyExist {
		h.logger.Errorf("%s: пользователь с таким логином уже существует", ResponseErrorKey)
		responseMessage = fmt.Sprintf("%s: пользователь с таким логином уже существует", ResponseErrorKey)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusBadRequest)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	responseData := domain.CorrectResponse{
		Data: successfulUserCreate,
	}

	result, err := json.Marshal(&responseData)
	if err != nil {
		h.logger.Errorf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	err = response.CorrectResponseWriter(w, result, http.StatusOK)
	if err != nil {
		h.logger.Errorf("ошибка при получении ответа -%s", err)
		return
	}
}

func (h *Handler) SetComment(w http.ResponseWriter, r *http.Request) {
	var responseMessage string

	reader := io.Reader(r.Body)
	requestBody, err := io.ReadAll(reader)
	if err != nil {
		h.logger.Errorf("%s: ошибка при чтении тела запроса - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при чтении тела запроса - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	currentRequestBody := &domain.TodoListComment{}
	err = json.Unmarshal(requestBody, currentRequestBody)
	if err != nil {
		h.logger.Errorf("%s: ошибка при десериализации объекта - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при десериализации объекта - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	err = h.cache.SetComment(*currentRequestBody)
	if err != nil {
		h.logger.Errorf("ошибка при создании комментария - %s", err)
		responseMessage = fmt.Sprintf("%s: ошибка при добавлении комментария", ResponseErrorKey)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusServiceUnavailable)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	responseData := domain.CorrectResponse{
		Data: successfulCommentCreate,
	}

	result, err := json.Marshal(&responseData)
	if err != nil {
		h.logger.Errorf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
		responseMessage = fmt.Sprintf("%s: ошибка при сериализации объекта - %s", ResponseErrorKey, err)
		err = response.ErrorResponseWriter(w, []byte(responseMessage), http.StatusInternalServerError)
		if err != nil {
			h.logger.Errorf("ошибка при получении ответа -%s", err)
			return
		}

		return
	}

	err = response.CorrectResponseWriter(w, result, http.StatusOK)
	if err != nil {
		h.logger.Errorf("ошибка при получении ответа -%s", err)
		return
	}
}
