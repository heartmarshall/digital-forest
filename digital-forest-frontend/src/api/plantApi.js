import axios from 'axios';

// Базовый URL вашего бэкенда. Выносим в константу, чтобы легко менять.
const API_BASE_URL = 'http://localhost:8080/v1';

/**
 * Запрашивает случайный набор растений с сервера.
 * @param {number} count - Количество растений для запроса.
 * @returns {Promise<Array>} - Промис, который разрешается массивом растений.
 */
export const fetchRandomPlants = async (count = 15) => {
  try {
    // Отправляем GET-запрос на /plants/random с параметром count
    const response = await axios.get(`${API_BASE_URL}/plants/random`, {
      params: {
        count: count
      }
    });
    // axios заворачивает данные ответа в поле `data`
    // В нашем случае, API возвращает объект { plants: [...], count: X }
    return response.data.plants; 
  } catch (error) {
    // Если произошла ошибка (сервер недоступен, 404 и т.д.),
    // выводим ее в консоль и пробрасываем дальше.
    console.error("Ошибка при загрузке растений:", error);
    throw error;
  }
};

/**
 * Отправляет данные нового растения на сервер.
 * @param {string} author - Имя автора.
 * @param {string} imageData - Данные изображения в base64.
 * @returns {Promise<Object>} - Промис, который разрешается созданным объектом растения.
 */
export const createPlant = async (author, imageData) => {
  try {
    // Отправляем POST-запрос на /plants
    const response = await axios.post(`${API_BASE_URL}/plants`, {
      author: author,
      imageData: imageData
    });
    return response.data;
  } catch (error) {
    console.error("Ошибка при создании растения:", error);
    throw error;
  }
};